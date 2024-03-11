package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// HashLogin Sent from web client to app server
const HashLogin = "hlogin"

func NewHashLogin(tokenHash []byte, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(AuthStatus, &HashLoginPayload{tokenHash}, key)
}

type HashLoginPayload struct {
	Hash []byte `json:"h"`
}
