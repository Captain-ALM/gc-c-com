package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameCommit Sent from web client to app server
const GameCommit = "cgame"

func NewGameCommit(index uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameCommit, &GameAnswerPayload{index}, key)
}

//This packet uses the GameAnswerPayload
