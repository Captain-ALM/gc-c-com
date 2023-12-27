package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// CurrentStatus Sent from app server to master server
const CurrentStatus = "cstat"

func NewCurrentStatus(id int, current int, max int, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(CurrentStatus, &CurrentStatusPayload{id, current, max}, key)
}

type CurrentStatusPayload struct {
	ID      int `json:"i"`
	Current int `json:"c"`
	Max     int `json:"m"`
}
