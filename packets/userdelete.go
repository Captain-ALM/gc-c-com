package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// UserDelete Sent from web client to app server
const UserDelete = "udel"

func NewUserDelete(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(UserDelete, nil, key)
}
