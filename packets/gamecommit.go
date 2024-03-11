package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// GameCommit Sent from web client to app server
const GameCommit = "cgame"

func NewGameCommit(index uint32, questionNumber uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameCommit, &GameAnswerPayload{questionNumber, index}, key)
}

//This packet uses the GameAnswerPayload
