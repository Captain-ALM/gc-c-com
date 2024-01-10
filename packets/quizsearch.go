package packets

import (
	"crypto/rsa"
	"errors"
	"golang.local/gc-c-com/packet"
	"strings"
)

// QuizSearch Sent from web client to app server
const QuizSearch = "squiz"

func NewQuizSearch(name string, filter EnumQuizSearchFilter, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizSearch, &QuizSearchPayload{name, filter}, key)
}

type QuizSearchPayload struct {
	Name   string               `json:"n"`
	Filter EnumQuizSearchFilter `json:"f"`
}

type EnumQuizSearchFilter string

func (e *EnumQuizSearchFilter) UnmarshalJSON(bytes []byte) error {
	if e == nil {
		return errors.New("packets.EnumQuizSearchFilter: UnmarshalJSON on nil pointer")
	}
	*e = EnumQuizSearchFilter(strings.Trim(string(bytes), "\""))
	return nil
}

func (e *EnumQuizSearchFilter) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("\"null\""), nil
	}
	return []byte("\"" + *e + "\""), nil
}

const (
	EnumQuizSearchFilterAll        = EnumQuizSearchFilter("all")
	EnumQuizSearchFilterOtherUsers = EnumQuizSearchFilter("othr")
	EnumQuizSearchFilterMine       = EnumQuizSearchFilter("mine")
	EnumQuizSearchFilterMyPublic   = EnumQuizSearchFilter("mpub")
	EnumQuizSearchFilterMyPrivate  = EnumQuizSearchFilter("mprv")
)
