package packets

import (
	"crypto/rsa"
	"golang.local/gc-c-com/packet"
)

// HostedGame Sent from app server to web client
const HostedGame = "hgame"

func NewHostedGame(gameID uint32, guestID uint32, guests []JoinGamePayload, removedGuests []JoinGamePayload, key *rsa.PrivateKey) (*packet.Packet, error) {
	return packet.New(HostedGame, &HostedGamePayload{gameID, guestID, guests, removedGuests}, key)
}

type HostedGamePayload struct {
	ID            uint32            `json:"i"`
	GuestID       uint32            `json:"gi,omitempty"`
	Guests        []JoinGamePayload `json:"gs"`
	GuestsRemoved []JoinGamePayload `json:"gr"`
}

// HostedGamePayload uses JoinGamePayload
