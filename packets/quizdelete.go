package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// QuizDelete Sent from web client to app server
const QuizDelete = "dquiz"

func NewQuizDelete(id uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizDelete, &IDPayload{id}, key)
}

// IDPayload is in use for this packet as a payload
