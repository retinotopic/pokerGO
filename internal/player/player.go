package player

import (
	"github.com/retinotopic/pokerGO/pkg/deck"

	"github.com/gorilla/websocket"
)

func NewPlayer() *Player {
	return &Player{InActive: &Active{p: &Player{}}, Active: &Active{p: &Player{}}, InGame: &inGame{p: &Player{}}, CurrentState: &Inactive{p: &Player{}}}
}

type Active struct {
	p *Player
}

func (s *Active) ToActive(Occupied map[int]bool) {}
func (s *Active) ToInactive(Occupied map[int]bool) {
	delete(Occupied, s.p.Place)
	s.p.SetState(s.p.InActive)
}
func (s *Active) ToGame(Occupied map[int]bool, ch chan<- struct{}) {
	if s.p.Admin == true && len(Occupied) >= 2 && s.p.MWPlayer.IsGame == true {
		s.p.SetState(s.p.InGame)
		ch <- struct{}{}
	}
}

type Inactive struct {
	p *Player
}

func (s *Inactive) ToActive(Occupied map[int]bool) {
	if Occupied[s.p.MWPlayer.Place] == false {
		s.p.Name = s.p.MWPlayer.Name
		s.p.Bankroll = s.p.MWPlayer.Stack
		s.p.Place = s.p.MWPlayer.Place
		Occupied[s.p.MWPlayer.Place] = true
		s.p.SetState(s.p.Active)
	}
}
func (s *Inactive) ToInactive(Occupied map[int]bool)                 {}
func (s *Inactive) ToGame(Occupied map[int]bool, ch chan<- struct{}) {}

type inGame struct {
	p *Player
}

func (s *inGame) ToActive(Occupied map[int]bool)   {}
func (s *inGame) ToInactive(Occupied map[int]bool) {}
func (s *inGame) ToGame(Occupied map[int]bool, ch chan<- struct{}) {
	// bet validation will be here
}

type Stater interface {
	ToActive(map[int]bool)
	ToInactive(map[int]bool)
	ToGame(map[int]bool, chan<- struct{})
}

type Player struct {
	Name         string `json:"Name"`
	Bankroll     int    `json:"Stack"`
	Conn         *websocket.Conn
	Place        int       `json:"Place"`
	Admin        bool      `json:"IsAdmin"`
	Hand         deck.Card `json:"Hand,omitempty"`
	ValueSec     int       `json:"Time,omitempty"`
	InActive     Stater
	Active       Stater
	InGame       Stater
	CurrentState Stater
	MWPlayer     MiddlewarePlayer
}
type MiddlewarePlayer struct {
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
func (p *Player) SetState(str Stater) { p.CurrentState = str }

func (p *Player) ChangeState(Occupied map[int]bool, sch chan<- struct{}) {
	if p.MWPlayer.IsGame == true {
		p.CurrentState.ToGame(Occupied, sch)
	} else if p.MWPlayer.IsActive == true {
		p.CurrentState.ToActive(Occupied)
	} else {
		p.CurrentState.ToInactive(Occupied)
	}
}
