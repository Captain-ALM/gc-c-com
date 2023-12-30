package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// KickGuest Sent from web client to app server
const KickGuest = "kg"

func NewKickGuest(id uint32, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(KickGuest, &IDPayload{id}, key)
}

// IDPayload is in use for this packet as a payload
