package transport

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ListenHandler struct {
	active       bool
	connectEvent func(l Listener, t Transport) Transport
	handlerMap   map[string]*Handler
	handlerMutex *sync.Mutex
	closeEvent   func(t Transport, e error)
	timeout      time.Duration
}

func (l *ListenHandler) Activate() {
	if l == nil || l.IsActive() {
		return
	}
	l.handlerMap = make(map[string]*Handler)
	l.handlerMutex = &sync.Mutex{}
	l.active = true
}

func (l *ListenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !l.IsActive() {
		return
	}
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	if request.URL.Query().Has("s") {
		chid := request.URL.Query().Get("s")
		l.handlerMutex.Lock()
		ch := l.handlerMap[chid]
		l.handlerMutex.Unlock()
		if ch == nil {
			writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			eMsg := "Session Not Found"
			writer.Header().Set("Content-Length", strconv.Itoa(len(eMsg)))
			writer.WriteHeader(http.StatusNotFound)
			_, _ = writer.Write([]byte(eMsg))
			return
		}
		ch.ServeHTTP(writer, request)
	} else if request.Method == http.MethodGet {
		nID, err := uuid.NewRandom()
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if err != nil {
			eMsg := "Internal Error: " + err.Error()
			writer.Header().Set("Content-Length", strconv.Itoa(len(eMsg)))
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(eMsg))
		}
		l.handlerMutex.Lock()
		defer l.handlerMutex.Unlock()
		hndl := &Handler{
			ID: nID.String() + "-" + request.RemoteAddr,
			closeEvent: func(t Transport, e error) {
				l.handlerMutex.Lock()
				defer l.handlerMutex.Unlock()
				delete(l.handlerMap, t.GetID())
				l.closeEvent(t, e)
			},
			timeout: l.timeout,
		}
		if l.connectEvent != nil {
			hndl = l.connectEvent(l, hndl).(*Handler)
		}
		l.handlerMap[hndl.ID] = hndl
		hndl.Activate()
		writer.Header().Set("Content-Length", strconv.Itoa(len(hndl.ID)))
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(hndl.ID))
	} else if request.Method == http.MethodOptions {
		writer.Header().Set("Allow", http.MethodOptions+", "+http.MethodGet)
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (l *ListenHandler) IsActive() bool {
	return l != nil && l.active
}

func (l *ListenHandler) Close() error {
	if l == nil {
		return nil
	}
	err := l.CloseTransports()
	l.active = false
	return err
}

func (l *ListenHandler) SetOnConnect(callback func(l Listener, t Transport) Transport) {
	if l == nil || callback == nil {
		return
	}
	l.connectEvent = callback
}

func (l *ListenHandler) SetOnClose(callback func(t Transport, e error)) {
	if l == nil || callback == nil {
		return
	}
	l.closeEvent = callback
}

func (l *ListenHandler) CloseTransports() error {
	if !l.IsActive() {
		return errors.New("listen handler inactive")
	}
	var err error
	l.handlerMutex.Lock()
	defer l.handlerMutex.Unlock()
	for _, handler := range l.handlerMap {
		er := handler.Close()
		if er != nil {
			err = er
		}
	}
	return err
}

func (l *ListenHandler) SetTimeout(to time.Duration) {
	if l == nil {
		return
	}
	l.timeout = to
	if !l.IsActive() {
		return
	}
	l.handlerMutex.Lock()
	defer l.handlerMutex.Unlock()
	for _, handler := range l.handlerMap {
		handler.SetTimeout(to)
	}
}

func (l *ListenHandler) GetTimeout() time.Duration {
	if l == nil {
		return 0
	}
	return l.timeout
}
