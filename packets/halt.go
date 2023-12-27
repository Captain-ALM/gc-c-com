package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// Halt Sent from master server to app server; sent from app server to web client
const Halt = "h"

func NewHalt(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(Halt, nil, key)
}
