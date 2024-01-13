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
	readLimit    int64
	Upgrader     websocket.Upgrader
}

func (l *ListenWebsocket) Activate() {
	if l == nil || l.IsActive() {
		return
	}
	l.socketMap = make(map[string]*Websocket)
	l.socketMutex = &sync.Mutex{}
	if l.readLimit < 4 {
		l.readLimit = 8192
	}
	l.active = true
}

func (l *ListenWebsocket) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !l.IsActive() {
		return
	}
	conn, err := l.Upgrader.Upgrade(writer, request, nil)
	if debugErrIsNil(err) {
		nID, err := uuid.NewRandom()
		if err != nil {
			debugErrIsNil(err)
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
			timeout:   l.timeout,
			readLimit: l.readLimit,
		}
		if l.acceptEvent != nil {
			socket, _ = l.acceptEvent(l, socket).(*Websocket)
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
	l.active = false
	err := l.intCloseTransports()
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

func (l *ListenWebsocket) getTransports() []Transport {
	l.socketMutex.Lock()
	defer l.socketMutex.Unlock()
	var trnsp []Transport
	for _, ctp := range l.socketMap {
		trnsp = append(trnsp, ctp)
	}
	return trnsp
}

func (l *ListenWebsocket) intCloseTransports() error {
	var err error
	trnsp := l.getTransports()
	for _, socket := range trnsp {
		er := socket.Close()
		if er != nil {
			err = er
		}
	}
	return err
}

func (l *ListenWebsocket) CloseTransports() error {
	if !l.IsActive() {
		return errors.New("listen handler inactive")
	}
	return l.intCloseTransports()
}

func (l *ListenWebsocket) SetTimeout(to time.Duration) {
	if l == nil || to < 0 {
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

func (l *ListenWebsocket) SetReadLimit(limit int64) {
	if l == nil || limit < 4 {
		return
	}
	l.readLimit = limit
	if !l.IsActive() {
		return
	}
	l.socketMutex.Lock()
	defer l.socketMutex.Unlock()
	for _, socket := range l.socketMap {
		socket.SetReadLimit(limit)
	}
}

func (l *ListenWebsocket) GetReadLimit() int64 {
	if l == nil {
		return 0
	}
	return l.readLimit
}
