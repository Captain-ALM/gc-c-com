package transport

import (
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type ListenWebsocket struct {
	active       bool
	acceptEvent  func(l Listener, t Transport) Transport
	connectEvent func(l Listener, t Transport)
	socketMap    map[string]*Websocket
	socketMutex  *sync.Mutex
	closeEvent   func(t Transport, e error)
	timeout      time.Duration
	Upgrader     websocket.Upgrader
}

func (l *ListenWebsocket) Activate() {
	if l == nil || l.IsActive() {
		return
	}
	l.socketMap = make(map[string]*Websocket)
	l.socketMutex = &sync.Mutex{}
	l.active = true
}

func (l *ListenWebsocket) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !l.IsActive() {
		return
	}
	conn, err := l.Upgrader.Upgrade(writer, request, nil)
	if err == nil {
		nID, err := uuid.NewRandom()
		if err != nil {
			return
		}
		l.socketMutex.Lock()
		defer l.socketMutex.Unlock()
		socket := &Websocket{
			ID: nID.String() + "-" + conn.NetConn().RemoteAddr().String(),
			closeEvent: func(t Transport, e error) {
				l.socketMutex.Lock()
				defer l.socketMutex.Unlock()
				delete(l.socketMap, t.GetID())
				l.closeEvent(t, e)
			},
			timeout: l.timeout,
		}
		if l.acceptEvent != nil {
			socket = l.acceptEvent(l, socket).(*Websocket)
		}
		if socket == nil {
			_ = conn.Close()
			return
		}
		l.socketMap[socket.ID] = socket
		socket.Activate(conn)
		if l.connectEvent != nil {
			l.connectEvent(l, socket)
		}
	}
}

func (l *ListenWebsocket) IsActive() bool {
	return l != nil && l.active
}

func (l *ListenWebsocket) Close() error {
	if l == nil {
		return nil
	}
	err := l.CloseTransports()
	l.active = false
	return err
}

func (l *ListenWebsocket) SetOnAccept(callback func(l Listener, t Transport) Transport) {
	if l == nil || callback == nil {
		return
	}
	l.acceptEvent = callback
}

func (l *ListenWebsocket) SetOnConnect(callback func(l Listener, t Transport)) {
	if l == nil || callback == nil {
		return
	}
	l.connectEvent = callback
}

func (l *ListenWebsocket) SetOnClose(callback func(t Transport, e error)) {
	if l == nil || callback == nil {
		return
	}
	l.closeEvent = callback
}

func (l *ListenWebsocket) CloseTransports() error {
	if !l.IsActive() {
		return errors.New("listen handler inactive")
	}
	var err error
	l.socketMutex.Lock()
	defer l.socketMutex.Unlock()
	for _, socket := range l.socketMap {
		er := socket.Close()
		if er != nil {
			err = er
		}
	}
	return err
}

func (l *ListenWebsocket) SetTimeout(to time.Duration) {
	if l == nil {
		return
	}
	l.timeout = to
	if !l.IsActive() {
		return
	}
	l.socketMutex.Lock()
	defer l.socketMutex.Unlock()
	for _, socket := range l.socketMap {
		socket.SetTimeout(to)
	}
}

func (l *ListenWebsocket) GetTimeout() time.Duration {
	if l == nil {
		return 0
	}
	return l.timeout
}
