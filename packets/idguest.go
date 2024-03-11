package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// IDGuest Sent from app server to web client or web client to app server
const IDGuest = "ig"

func NewIDGuest(id uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(IDGuest, &IDPayload{id}, key)
}

// IDPayload is used as a payload for this packet
