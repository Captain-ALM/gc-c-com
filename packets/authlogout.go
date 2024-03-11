package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// AuthLogout Sent from web client to app server
const AuthLogout = "alout"

func NewAuthLogout(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(AuthLogout, nil, key)
}
