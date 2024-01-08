package player

import (
	"pokerGO/pkg/deck"

	"github.com/gorilla/websocket"
)

func NewPlayer() *Player {
	return &Player{inActive: &Active{p: &Player{}}, Active: &Active{p: &Player{}}, inGame: &inGame{p: &Player{}}, CurrentState: &inActive{p: &Player{}}}
}

type Active struct {
	p *Player
}

func (s *Active) ToActive(Occupied map[int]bool) {}
func (s *Active) ToInactive(Occupied map[int]bool) {
	s.p.IsActive = s.p.MWPlayer.IsActive
	delete(Occupied, s.p.Place)
	s.p.SetState(s.p.inActive)
}
func (s *Active) ToGame(Occupied map[int]bool) {
	if s.p.Admin == true && len(Occupied) >= 2 && s.p.MWPlayer.IsGame == true {
		s.p.SetState(s.p.inGame)
	}
}

type inActive struct {
	p *Player
}

func (s *inActive) ToActive(Occupied map[int]bool) {
	if Occupied[s.p.MWPlayer.Place] == false {
		s.p.Name = s.p.MWPlayer.Name
		s.p.Bankroll = s.p.MWPlayer.Stack
		s.p.IsActive = s.p.MWPlayer.IsActive
		s.p.Place = s.p.MWPlayer.Place
		Occupied[s.p.MWPlayer.Place] = true
		s.p.SetState(s.p.Active)
	}
}
func (s *inActive) ToInactive(Occupied map[int]bool) {}
func (s *inActive) ToGame(Occupied map[int]bool)     {}

type inGame struct {
	p *Player
}

func (s *inGame) ToActive(Occupied map[int]bool)   {}
func (s *inGame) ToInactive(Occupied map[int]bool) {}
func (s *inGame) ToGame(Occupied map[int]bool) {
	// bet validation will be here
}

type Stater interface {
	ToActive(map[int]bool)
	ToInactive(map[int]bool)
	ToGame(map[int]bool)
}

type Player struct {
	Name         string `json:"Name"`
	Bankroll     int    `json:"Stack"`
	IsActive     bool   `json:"IsActive"`
	IsGame       bool   `json:"IsGame"`
	Conn         *websocket.Conn
	Place        int       `json:"Place"`
	Admin        bool      `json:"IsAdmin"`
	Hand         deck.Card `json:"Hand,omitempty"`
	ValueSec     int       `json:"Time,omitempty"`
	inActive     Stater
	Active       Stater
	inGame       Stater
	CurrentState Stater
	Pl           *Player
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

func (p *Player) ChangeState(Occupied map[int]bool) {
	if p.MWPlayer.IsGame == true {
		p.CurrentState.ToGame(Occupied)
	} else if p.MWPlayer.IsActive == true {
		p.CurrentState.ToActive(Occupied)
	} else {
		p.CurrentState.ToInactive(Occupied)
	}
}
