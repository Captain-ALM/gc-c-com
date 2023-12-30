package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameEnd Sent from web client to app server
const GameEnd = "end"

func NewGameEnd(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameEnd, nil, key)
}