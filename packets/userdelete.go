package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// UserDelete Sent from web client to app server
const UserDelete = "udel"

func NewUserDelete(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(UserDelete, nil, key)
}
