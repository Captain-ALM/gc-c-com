package transport

import (
	"golang.org/x/exp/slices"
	"sync"
	"time"
)

func NewMultiListener(listeners []Listener, onAccept func(l Listener, t Transport) Transport, onConnect func(l Listener, t Transport), onClose func(t Transport, e error), timeout time.Duration) *MultiListener {
	mL := &MultiListener{
		init:          true,
		acceptEvent:   onAccept,
		connectEvent:  onConnect,
		closeEvent:    onClose,
		listeners:     listeners,
		listenerMutex: &sync.Mutex{},
		timeout:       timeout,
	}
	for _, cl := range listeners {
		cl.SetOnAccept(onAccept)
		cl.SetOnConnect(onConnect)
		cl.SetOnClose(onClose)
		cl.SetTimeout(timeout)
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
	if m == nil || !m.init || callback == nil {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	m.acceptEvent = callback
	for _, cl := range m.listeners {
		cl.SetOnAccept(callback)
	}
}

func (m *MultiListener) SetOnConnect(callback func(l Listener, t Transport)) {
	if m == nil || !m.init || callback == nil {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	m.connectEvent = callback
	for _, cl := range m.listeners {
		cl.SetOnConnect(callback)
	}
}

func (m *MultiListener) SetOnClose(callback func(t Transport, e error)) {
	if m == nil || !m.init || callback == nil {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	m.closeEvent = callback
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
	if m == nil || !m.init {
		return
	}
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()
	m.timeout = to
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
