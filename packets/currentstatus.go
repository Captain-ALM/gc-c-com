package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// CurrentStatus Sent from app server to master server
const CurrentStatus = "cstat"

func NewCurrentStatus(id uint32, current uint32, max uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(CurrentStatus, &CurrentStatusPayload{id, current, max}, key)
}

type CurrentStatusPayload struct {
	ID      uint32 `json:"i"`
	Current uint32 `json:"c"`
	Max     uint32 `json:"m"`
}
