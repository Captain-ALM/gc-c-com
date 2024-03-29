package transport

import (
	"github.com/Captain-ALM/gc-c-com/packet"
	"time"
)

type Transport interface {
	GetID() string
	IsActive() bool
	Send(p *packet.Packet) error
	Receive() (*packet.Packet, error)
	Close() error
	SetOnClose(callback func(t Transport, e error))
	SetTimeout(to time.Duration)
	GetTimeout() time.Duration
	SetReadLimit(limit int64)
	GetReadLimit() int64
}
