package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// AuthCheck Sent from web client to app server
const AuthCheck = "acheck"

func NewAuthCheck(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(AuthCheck, nil, key)
}
