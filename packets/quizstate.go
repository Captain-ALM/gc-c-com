package packets

import (
	"crypto/rsa"
	"errors"
	"golang.local/gc-c-com/packet"
)

// QuizState Sent from app server to web client
const QuizState = "qzstat"

func NewQuizState(id uint32, state EnumQuizState, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(QuizState, &QuizStatePayload{id, state}, key)
}

type QuizStatePayload struct {
	ID    uint32        `json:"i"`
	State EnumQuizState `json:"s"`
}

type EnumQuizState string

func (e *EnumQuizState) UnmarshalJSON(bytes []byte) error {
	if e == nil {
		return errors.New("packets.EnumAuthStatus: UnmarshalJSON on nil pointer")
	}
	*e = EnumQuizState(bytes)
	return nil
}

func (e *EnumQuizState) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}
	return []byte(*e), nil
}

const (
	EnumQuizStateNotFound     = EnumQuizState("404")
	EnumQuizStateUploadFailed = EnumQuizState("403")
	EnumQuizStateDeleted      = EnumQuizState("202")
	EnumQuizStateCreated      = EnumQuizState("204")
	EnumQuizStatePublic       = EnumQuizState("pub")
	EnumQuizStatePrivate      = EnumQuizState("prv")
)
