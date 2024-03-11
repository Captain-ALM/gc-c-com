package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// QuizList Sent from app server to web client
const QuizList = "lquiz"

func NewQuizList(entries []QuizListEntry, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizList, &QuizListPayload{entries}, key)
}

type QuizListPayload struct {
	Entries []QuizListEntry `json:"e"`
}

type QuizListEntry struct {
	ID     uint32 `json:"i"`
	Name   string `json:"n"`
	Mine   bool   `json:"m"`
	Public bool   `json:"p"`
}
