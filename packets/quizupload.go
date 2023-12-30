package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// QuizUpload Sent from web client to app server
const QuizUpload = "uquiz"

func NewQuizUpload(id uint32, name string, questions QuizQuestions, answers QuizAnswers, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizUpload, &QuizDataPayload{id, name, questions, answers}, key)
}

//This uses QuizDataPayload as the packet payload
