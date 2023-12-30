package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameNotFound Sent from app server to web client
const GameNotFound = "g404"

func NewGameNotFound(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameNotFound, nil, key)
}
