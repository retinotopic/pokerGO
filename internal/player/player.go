package player

import (
	"pokerGO/pkg/deck"

	"github.com/gorilla/websocket"
)

func NewPlayer() *Player { return &Player{} }

type Player struct {
	Name     string `json:"Name"`
	Bankroll int    `json:"Stack"`
	IsActive bool   `json:"IsActive"`
	Conn     *websocket.Conn
	Place    int  `json:"Place"`
	Admin    bool `json:"IsAdmin"`
	Hand     deck.Card
	id       int
}
