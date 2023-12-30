package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameCountdown Sent from app server to web client
const GameCountdown = "dgame"

func NewGameCountdown(value uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameCountdown, &GameValuePayload{value}, key)
}

type GameValuePayload struct {
	Value uint32 `json:"v"`
}

// GameValuePayload is also used by GameScore
