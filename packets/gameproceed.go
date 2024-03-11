package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// GameProceed Sent from web client to app server
const GameProceed = "pgame"

func NewGameProceed(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameProceed, nil, key)
}
