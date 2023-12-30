package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// JoinGame Sent from web client to app server
const JoinGame = "jgame"

func NewJoinGame(gameID uint32, nickName string, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(JoinGame, &JoinGamePayload{gameID, nickName}, key)
}

type JoinGamePayload struct {
	ID       uint32 `json:"i"`
	Nickname string `json:"n"`
}

// JoinGamePayload used by HostedGamePayload
