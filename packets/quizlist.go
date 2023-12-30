package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// QuizList Sent from web client to app server
const QuizList = "lquiz"

func NewQuizList(entries []QuizListEntry, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizList, &QuizListPayload{entries}, key)
}

type QuizListPayload struct {
	Entries []QuizListEntry `json:"e"`
}

type QuizListEntry struct {
	ID     int    `json:"i"`
	Name   string `json:"n"`
	Mine   bool   `json:"m"`
	Public bool   `json:"p"`
}
