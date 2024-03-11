package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// GameLeave Sent from web client to app server
const GameLeave = "lev"

func NewGameLeave(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameLeave, nil, key)
}
