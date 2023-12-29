package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// AuthLogout Sent from web client to app server
const AuthLogout = "alout"

func NewAuthLogout(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(AuthLogout, nil, key)
}
