package transport

import (
	"golang.local/gc-c-com/packet"
	"time"
)

type Transport interface {
	GetID() string
	IsActive() bool
	Send(p *packet.Packet) error
	Receive() (*packet.Packet, error)
	Close() error
	SetOnClose(callback CloseCallback)
	SetTimeout(to time.Duration)
	GetTimeout() time.Duration
}

type CloseCallback = func(t Transport, e error)
