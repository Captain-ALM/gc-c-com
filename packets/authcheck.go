package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// AuthCheck Sent from web client to app server
const AuthCheck = "acheck"

func NewAuthCheck(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(AuthCheck, nil, key)
}
