package packets

import (
	"crypto/rsa"
	"errors"
	"golang.local/gc-c-com/packet"
)

// AuthStatus Sent from app server to web client (Except EnumAuthStatusRequest sent from web client to app server)
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
	*e = EnumAuthStatus(bytes)
	return nil
}

func (e *EnumAuthStatus) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	return []byte(*e), nil
}

const (
	// EnumAuthStatusRequest Web client requests auth status
	EnumAuthStatusRequest   = EnumAuthStatus("request")
	EnumAuthStatusRequired  = EnumAuthStatus("required")
	EnumAuthStatusSignedOut = EnumAuthStatus("none")
	EnumAuthStatusLoggedOut
	EnumAuthStatusSignedIn = EnumAuthStatus("active")
	EnumAuthStatusLoggedIn
	EnumAuthStatusAccepted = EnumAuthStatus("accepted")
	EnumAuthStatusRejected = EnumAuthStatus("rejected")
)
