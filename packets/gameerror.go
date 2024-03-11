package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// GameError Sent from app server to web client
const GameError = "egame"

func NewGameError(message string, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameError, &GameMessagePayload{message}, key)
}

type GameMessagePayload struct {
	Message string `json:"m"`
}

// GameMessagePayload is also used by GameStatus
