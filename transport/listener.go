package transport

import "time"

type ConnectCallback = func(l Listener, t Transport) Transport

type Listener interface {
	IsActive() bool
	Close() error
	SetOnConnect(callback ConnectCallback)
	SetOnClose(callback CloseCallback)
	CloseTransports() error
	SetTimeout(to time.Duration)
	GetTimeout() time.Duration
}
