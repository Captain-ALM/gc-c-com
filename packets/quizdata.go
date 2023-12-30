package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// QuizData Sent from app server to web client
const QuizData = "quiz"

func NewQuizData(id int, name string, questions QuizQuestions, answers QuizAnswers, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizData, &QuizDataPayload{id, name, questions, answers}, key)
}

type QuizDataPayload struct {
	ID        int           `json:"i"`
	Name      string        `json:"n"`
	Questions QuizQuestions `json:"q"`
	Answers   QuizAnswers   `json:"a"`
}

// QuizDataPayload is also used by QuizUpload as a packet payload

type QuizQuestions struct {
	Questions []QuizQuestion `json:"qs"`
}

type QuizQuestion struct {
	Type     string `json:"t"`
	Question string `json:"q"`
}

//QuizQuestion is also used by GameQuestion

type QuizAnswers struct {
	Answers []QuizAnswerSet `json:"as"`
}

type QuizAnswerSet struct {
	CorrectAnswer int          `json:"ca"`
	Answers       []QuizAnswer `json:"as"`
}

type QuizAnswer struct {
	Answer string `json:"a"`
	Color  int    `json:"c"`
}

//QuizAnswer is also used by GameQuestion
