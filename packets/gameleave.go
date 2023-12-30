package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameLeave Sent from web client to app server
const GameLeave = "lev"

func NewGameLeave(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameLeave, nil, key)
}
