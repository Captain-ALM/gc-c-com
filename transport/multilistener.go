package transport

import (
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

func NewMultiListener(listeners []Listener, onAccept func(l Listener, t Transport) Transport, onConnect func(l Listener, t Transport), onClose func(t Transport, e error), timeout time.Duration, readLimit int64) *MultiListener {
	if readLimit < 4 {
		readLimit = 4
	}
	mL := &MultiListener{
		init:          true,
		acceptEvent:   onAccept,
		connectEvent:  onConnect,
		closeEvent:    onClose,
		listeners:     listeners,
		listenerMutex: &sync.Mutex{},
		timeout:       timeout,
		readLimit:     readLimit,
	}
	for _, cl := range listeners {
		cl.SetOnAccept(onAccept)
		cl.SetOnConnect(onConnect)
		cl.SetOnClose(onClose)
		cl.SetTimeout(timeout)
		cl.SetReadLimit(readLimit)
	}
	return mL
}

type MultiListener struct {
	init          bool
	acceptEvent   func(l Listener, t Transport) Transport
	connectEvent  func(l Listener, t Transport)
	closeEvent    func(t Transport, e error)
	listeners     []Listener
	listenerMutex *sync.Mutex
	timeout       time.Duration
	readLimit     int64
}

func (m *MultiListener) IsActive() bool {
	if m == nil {
		return false
	}
	return m.init
}

func (m *MultiListener) Close() error {
	if m == nil || !m.init {
		return nil
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	var toRet error
	for _, cl := range m.listeners {
		err := cl.Close()
		if err != nil {
			toRet = err
		}
	}
	return toRet
}

func (m *MultiListener) SetOnAccept(callback func(l Listener, t Transport) Transport) {
	if m == nil || callback == nil {
		return
	}
	m.acceptEvent = callback
	if !m.init {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	for _, cl := range m.listeners {
		cl.SetOnAccept(callback)
	}
}

func (m *MultiListener) SetOnConnect(callback func(l Listener, t Transport)) {
	if m == nil || callback == nil {
		return
	}
	m.connectEvent = callback
	if !m.init {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	for _, cl := range m.listeners {
		cl.SetOnConnect(callback)
	}
}

func (m *MultiListener) SetOnClose(callback func(t Transport, e error)) {
	if m == nil || callback == nil {
		return
	}
	m.closeEvent = callback
	if !m.init {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	for _, cl := range m.listeners {
		cl.SetOnClose(callback)
	}
}

func (m *MultiListener) CloseTransports() error {
	if m == nil || !m.init {
		return nil
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	var toRet error
	for _, cl := range m.listeners {
		err := cl.CloseTransports()
		if err != nil {
			toRet = err
		}
	}
	return toRet
}

func (m *MultiListener) SetTimeout(to time.Duration) {
	if m == nil || to < 0 {
		return
	}
	m.timeout = to
	if !m.init {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	for _, cl := range m.listeners {
		cl.SetTimeout(to)
	}
}

func (m *MultiListener) GetTimeout() time.Duration {
	if m == nil {
		return 0
	}
	return m.timeout
}

func (m *MultiListener) GetListeners() []Listener {
	if m == nil || !m.init {
		return nil
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	lsns := make([]Listener, len(m.listeners))
	copy(lsns, m.listeners)
	return lsns
}

func (m *MultiListener) AddListener(l Listener) {
	if m == nil || !m.init || l == nil {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	m.listeners = append(m.listeners, l)
	l.SetOnAccept(m.acceptEvent)
	l.SetOnConnect(m.connectEvent)
	l.SetOnClose(m.closeEvent)
	l.SetTimeout(m.timeout)
	l.SetReadLimit(m.readLimit)
}

func (m *MultiListener) RemoveListener(l Listener) {
	if m == nil || !m.init || l == nil {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	m.listeners = slices.DeleteFunc(m.listeners, func(cl Listener) bool {
		return cl == l
	})
}

func (m *MultiListener) ClearListeners() {
	if m == nil || !m.init {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	m.listeners = nil
}

func (m *MultiListener) SetReadLimit(limit int64) {
	if m == nil || limit < 4 {
		return
	}
	m.readLimit = limit
	if !m.init {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	for _, cl := range m.listeners {
		cl.SetReadLimit(limit)
	}
}

func (m *MultiListener) GetReadLimit() int64 {
	if m == nil {
		return 0
	}
	return m.readLimit
}
