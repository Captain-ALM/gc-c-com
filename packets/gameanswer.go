package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// GameAnswer Sent from app server to web client
const GameAnswer = "agame"

func NewGameAnswer(index uint32, questionNumber uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameAnswer, &GameAnswerPayload{questionNumber, index}, key)
}

type GameAnswerPayload struct {
	QuestionNumber uint32 `yaml:"q"`
	Index          uint32 `json:"x"`
}

//The GameAnswerPayload is also in use by GameCommit
