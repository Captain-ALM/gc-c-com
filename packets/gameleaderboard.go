package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameLeaderboard Sent from app server to web client
const GameLeaderboard = "lgame"

func NewGameLeaderboard(entries []GameLeaderboardEntry, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameLeaderboard, &GameLeaderboardPayload{entries}, key)
}

type GameLeaderboardPayload struct {
	Entries []GameLeaderboardEntry `json:"e"`
}

type GameLeaderboardEntry struct {
	ID       uint32 `json:"i"`
	Nickname string `json:"n"`
	Score    uint32 `json:"s"`
	Streak   uint32 `json:"t,omitempty"`
}
