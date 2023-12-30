package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameAnswer Sent from app server to web client
const GameAnswer = "agame"

func NewGameAnswer(index int, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameAnswer, &GameAnswerPayload{index}, key)
}

type GameAnswerPayload struct {
	Index int `json:"x"`
}
