package packet

import (
	"crypto/rsa"
	"encoding/json"
)

const (
	Ping = "i"
	Pong = "o"
)

func NewPing() *Packet {
	return &Packet{Command: Ping}
}

func NewPong() *Packet {
	return &Packet{Command: Pong}
}

func New(command string, payload any, key *rsa.PrivateKey) (*Packet, error) {
	var py []byte
	var err error
	if payload != nil {
		py, err = json.Marshal(payload)
	}
	if err != nil {
		return nil, err
	}
	pk := &Packet{
		Command: command,
		Payload: py,
	}
	if key != nil {
		return pk, pk.Sign(key)
	}
	return pk, nil
}

func From(packet []byte) (*Packet, error) {
	var pk Packet
	err := json.Unmarshal(packet, &pk)
	if err != nil {
		return nil, err
	}
	return &pk, nil
}

func FromIgnoreError(packet []byte) *Packet {
	pk, err := From(packet)
	debugErrIsNil(err)
	return pk
}

func FromNew(pk *Packet, err error) *Packet {
	if debugErrIsNil(err) {
		return pk
	}
	return &Packet{}
}

func GetCommand(packet []byte) (string, error) {
	pk := struct {
		Command string `json:"c,omitempty"`
	}{}
	err := json.Unmarshal(packet, &pk)
	if err != nil {
		return "", err
	}
	return pk.Command, nil
}

func GetCommandIgnoreError(packet []byte) string {
	pk, err := GetCommand(packet)
	debugErrIsNil(err)
	return pk
}
