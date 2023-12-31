package transport

import "time"

type Listener interface {
	IsActive() bool
	Close() error
	SetOnConnect(callback func(l Listener, t Transport) Transport)
	SetOnClose(callback func(t Transport, e error))
	CloseTransports() error
	SetTimeout(to time.Duration)
	GetTimeout() time.Duration
}
