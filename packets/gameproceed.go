package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameProceed Sent from web client to app server
const GameProceed = "pgame"

func NewGameProceed(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameProceed, nil, key)
}
