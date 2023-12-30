package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameScore Sent from app server to web client
const GameScore = "sgame"

func NewGameScore(value uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameScore, &GameValuePayload{value}, key)
}

// GameValuePayload is used by this packet as a payload
