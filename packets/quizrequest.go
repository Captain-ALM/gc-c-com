package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// QuizRequest Sent from web client to app server
const QuizRequest = "rquiz"

func NewQuizRequest(id int, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizRequest, &IDPayload{id}, key)
}

// IDPayload is in use for this packet as a payload
