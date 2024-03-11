package packets

import (
	"crypto/rsa"
	"errors"
	"github.com/Captain-ALM/gc-c-com/packet"
	"strings"
)

// AuthStatus Sent from app server to web client
const AuthStatus = "astat"

func NewAuthStatus(status EnumAuthStatus, tokenHash []byte, userEmail string, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(AuthStatus, &AuthStatusPayload{status, tokenHash, userEmail}, key)
}

type AuthStatusPayload struct {
	Status    EnumAuthStatus `json:"s"`
	TokenHash []byte         `json:"t,omitempty"`
	UserEmail string         `json:"u,omitempty"`
}

type EnumAuthStatus string

func (e *EnumAuthStatus) UnmarshalJSON(bytes []byte) error {
	if e == nil {
		return errors.New("packets.EnumAuthStatus: UnmarshalJSON on nil pointer")
	}
	*e = EnumAuthStatus(strings.Trim(string(bytes), "\""))
	return nil
}

func (e *EnumAuthStatus) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("\"null\""), nil
	}
	return []byte("\"" + *e + "\""), nil
}

const (
	EnumAuthStatusRequired  = EnumAuthStatus("required")
	EnumAuthStatusSignedOut = EnumAuthStatus("none")
	EnumAuthStatusLoggedOut
	EnumAuthStatusSignedIn = EnumAuthStatus("active")
	EnumAuthStatusLoggedIn
	EnumAuthStatusAcceptedJWT  = EnumAuthStatus("acceptedjwt")
	EnumAuthStatusRejectedJWT  = EnumAuthStatus("rejectedjwt")
	EnumAuthStatusAcceptedHash = EnumAuthStatus("acceptedhsh")
	EnumAuthStatusRejectedHash = EnumAuthStatus("rejectedhsh")
)
