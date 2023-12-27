package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// QueryStatus Sent from master server to app server
const QueryStatus = "qstat"

func NewQueryStatus(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QueryStatus, nil, key)
}
