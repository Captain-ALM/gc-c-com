package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// GameStatus Sent from app server to web client
const GameStatus = "gstat"

func NewGameStatus(message string, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameStatus, &GameMessagePayload{message}, key)
}

// GameMessagePayload is used for this packet as a payload
