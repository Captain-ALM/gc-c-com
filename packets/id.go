package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// ID Sent signed from master server to app server; sent unsigned the other way
const ID = "id"

func NewID(id uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(ID, &IDPayload{id}, key)
}

type IDPayload struct {
	ID uint32 `json:"i"`
}

//This payload is also used by IDGuest, KickGuest, QuizRequest and QuizDelete
