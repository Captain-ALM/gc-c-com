package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// QuizVisibility Sent from web client to app server
const QuizVisibility = "vquiz"

func NewQuizVisibility(id uint32, isPublic bool, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizVisibility, &QuizVisibilityPayload{id, isPublic}, key)
}

type QuizVisibilityPayload struct {
	ID     uint32 `json:"i"`
	Public bool   `json:"p"`
}
