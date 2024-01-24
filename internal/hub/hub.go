package hub

import (
	"fmt"
	"sync"
	"time"

	"github.com/retinotopic/pokerGO/internal/player"

	"github.com/gorilla/websocket"

	"github.com/retinotopic/pokerGO/pkg/randfuncs"
)

func NewLobby() *Lobby {
	l := &Lobby{Players: make(map[string]*player.Player), Occupied: make(map[int]bool), PlayerCh: make(chan player.Player), StartGame: make(chan struct{}, 1)}
	return l
}

type Lobby struct {
	Players  map[string]*player.Player
	Admin    *player.Player
	Occupied map[int]bool
	sync.Mutex
	GameState  int
	AdminOnce  sync.Once
	PlayerCh   chan player.Player
	StartGame  chan struct{}
	PlayerTurn player.Player

	//Conns   chan *websocket.Conn
}

// <ol hx-swap-oob=beforeend:#piece>    <li>%v</li></ol>

func (l *Lobby) LobbyWork() {
	fmt.Println("im in")
	for {
		select {
		case x := <-l.PlayerCh: // broadcoasting one seat to everyone
			for _, v := range l.Players {
				v.Conn.WriteJSON(x)
			}
		case <-l.StartGame:
			l.Game()
		}
	}
}

func (l *Lobby) Connhandle(plr *player.Player, conn *websocket.Conn) {
	fmt.Println("im in2")
	l.AdminOnce.Do(func() {
		l.Admin = plr
		plr.Admin = true
	})

	defer func() {
		fmt.Println("rip connection")

	}()
	plr.Conn = conn
	for _, v := range l.Players { // load current state of the game
		vs := *v
		if vs != *plr {
			vs = v.PrivateSend()
		}
		err := plr.Conn.WriteJSON(vs)
		if err != nil {
			fmt.Println(err, "WriteJSON start")
		}
		fmt.Println("start")
	}
	for {
		err := plr.Conn.ReadJSON(player.MiddlewarePlayer{})
		if err != nil {
			fmt.Println(err, "conn read error")
			plr.Conn = nil
			break
		}
		plr.ChangeState(l.Occupied, l.StartGame)
		fmt.Println(player.MiddlewarePlayer{}, "pered v ch")
		l.PlayerCh <- *plr
	}
}
func (l *Lobby) Game() {
	plorder := []player.Player{}
	for _, v := range l.Players {
		if v.CurrentState == v.InGame {
			plorder = append(plorder, *v)
		}
	}
	timer := time.NewTicker(time.Second * 1)
	PlayerBroadcast := make(chan player.Player)
	k := randfuncs.NewSource().Intn(len(l.Occupied))
	l.PlayerTurn = plorder[k]
	for {
		select {
		case pb := <-PlayerBroadcast: // broadcoasting one to everyone
			for _, v := range l.Players {
				if pb != *v {
					pb = v.PrivateSend()
				}
				v.Conn.WriteJSON(pb)
				fmt.Println(v, "aaaaaaaaaaaaaaaaaaaaaaaaaaa")
			}
		case tick := <-timer.C:
			timevalue := tick.Second()
			PlayerBroadcast <- l.PlayerTurn.SendTimeValue(timevalue)
		case pl := <-l.PlayerCh: // evaluating hand
			if pl == l.PlayerTurn {
				// evaluate hand
			}
			for _, v := range l.Players {
				pls := pl
				if pls != *v {
					pls = v.PrivateSend()
				}
				v.Conn.WriteJSON(pls)
				fmt.Println(v, "898aaaaaaaaaaaaaaaaaaaaaaaaaa")
			}
		}
	}
}
