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
	Place    int       `json:"Place"`
	Admin    bool      `json:"IsAdmin"`
	Hand     deck.Card `json:"Hand,omitempty"`
	ValueSec int       `json:"Time,omitempty"`
}
type InterimPlayer struct {
	Name     string `json:"Name"`
	Stack    int    `json:"Stack"`
	IsActive bool   `json:"IsActive"`
	Place    int    `json:"Place"`
	Bet      int    `json:"Bet"`
	IsGame   bool   `json:"IsGame"`
}

func (p Player) PrivateSend() Player {
	p.Hand = deck.Card{}
	return p
}
func (p Player) SendTimeValue(time int) Player {
	return Player{ValueSec: time, Place: p.Place}
}
func (p *Player) ChangeState(Occupied map[int]bool, data InterimPlayer) {
	if p.IsActive == false && Occupied[data.Place] == false && data.IsActive == true {
		p.Name = data.Name
		p.Bankroll = data.Stack
		p.IsActive = data.IsActive
		p.Place = data.Place
		Occupied[data.Place] = true
	} else if p.IsActive == true && data.IsActive == false {
		p.IsActive = data.IsActive
		delete(Occupied, p.Place)
	}
}
