package packets

import (
	"crypto/rsa"
	"github.com/Captain-ALM/gc-c-com/packet"
)

// QueryStatus Sent from master server to app server
const QueryStatus = "qstat"

func NewQueryStatus(key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QueryStatus, nil, key)
}
