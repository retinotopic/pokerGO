package player

import (
	"github.com/gorilla/websocket"
)

func NewPlayer() *Player { return &Player{} }

type Player struct {
	Name     string
	Bankroll int
	Conn     *websocket.Conn
}
