package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// TokenLogin Sent from web client to app server
const TokenLogin = "tlogin"

func NewTokenLogin(jwtToken string, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(AuthStatus, &TokenLoginPayload{jwtToken}, key)
}

type TokenLoginPayload struct {
	Token string `json:"t"`
}
