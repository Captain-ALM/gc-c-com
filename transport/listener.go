package transport

import "time"

type Listener interface {
	IsActive() bool
	Close() error
	SetOnAccept(callback func(l Listener, t Transport) Transport)
	SetOnConnect(callback func(l Listener, t Transport))
	SetOnClose(callback func(t Transport, e error))
	CloseTransports() error
	SetTimeout(to time.Duration)
	GetTimeout() time.Duration
}
