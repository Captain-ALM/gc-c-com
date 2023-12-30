package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// NewGame Sent from web client to app server
const NewGame = "ngame"

func NewNewGame(quizID uint32, maxCountdown uint32, streakEnabled bool, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(NewGame, &NewGamePayload{quizID, maxCountdown, streakEnabled}, key)
}

type NewGamePayload struct {
	QuizID        uint32 `json:"qi"`
	MaxCountdown  uint32 `json:"mc"`
	StreakEnabled bool   `json:"se"`
}
