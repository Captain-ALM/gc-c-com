package packet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"errors"
)

type Packet struct {
	Command   string          `json:"c,omitempty"`
	Payload   json.RawMessage `json:"p,omitempty"`
	Signature []byte          `json:"s,omitempty"`
}

func (p *Packet) GetCommand() string {
	if p == nil {
		return ""
	}
	return p.Command
}

func (p *Packet) Sign(key *rsa.PrivateKey) error {
	if p == nil {
		return errors.New("packet is nil")
	}
	if key == nil {
		return errors.New("key is nil")
	}
	hasher := sha512.New()
	_, err := hasher.Write([]byte(p.Command))
	if err != nil {
		return err
	}
	_, err = hasher.Write(p.Payload)
	if err != nil {
		return err
	}
	p.Signature, err = rsa.SignPSS(rand.Reader, key, crypto.SHA512, hasher.Sum(nil), nil)
	return err
}

func (p *Packet) Verify(key *rsa.PublicKey) error {
	if p == nil {
		return errors.New("packet is nil")
	}
	if key == nil {
		return errors.New("key is nil")
	}
	hasher := sha512.New()
	_, err := hasher.Write([]byte(p.Command))
	if err != nil {
		return err
	}
	_, err = hasher.Write(p.Payload)
	if err != nil {
		return err
	}
	return rsa.VerifyPSS(key, crypto.SHA512, hasher.Sum(nil), p.Signature, nil)
}

func (p *Packet) Valid(key *rsa.PublicKey) bool {
	return p != nil && (p.Command != "" && (len(p.Signature) == 0 || key == nil || p.Verify(key) == nil))
}

// GetPayload parameter v must be a pointer to the payload object type
func (p *Packet) GetPayload(v any) error {
	if p == nil {
		return errors.New("packet is nil")
	}
	if len(p.Payload) == 0 {
		return nil
	}
	return json.Unmarshal(p.Payload, v)
}

func (p *Packet) ToBytes() ([]byte, error) {
	if p == nil {
		return nil, errors.New("packet is nil")
	}
	return json.Marshal(p)
}

func (p *Packet) ToBytesIgnoreError() []byte {
	bts, _ := p.ToBytes()
	return bts
}
