package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// GameLeaderboard Sent from app server to web client
const GameLeaderboard = "lgame"

func NewGameLeaderboard(score uint32, streak uint32, entries []GameLeaderboardEntry, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(GameLeaderboard, &GameLeaderboardPayload{score, streak, entries}, key)
}

type GameLeaderboardPayload struct {
	MyScore  uint32                 `json:"s"`
	MyStreak uint32                 `json:"t,omitempty"`
	Entries  []GameLeaderboardEntry `json:"e"`
}

type GameLeaderboardEntry struct {
	Nickname string `json:"n"`
	Score    uint32 `json:"s"`
	Streak   uint32 `json:"t,omitempty"`
}
