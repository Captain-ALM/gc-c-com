package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// GameQuestion Sent from app server to web client
const GameQuestion = "qgame"

func NewGameQuestion(question QuizQuestion, answers QuizAnswerSet, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameQuestion, &GameQuestionPayload{question, answers.Answers}, key)
}

type GameQuestionPayload struct {
	Question QuizQuestion `json:"q"`
	Answers  []QuizAnswer `json:"a"`
}

//The payload uses QuizQuestion and QuizAnswer structures used by QuizData
