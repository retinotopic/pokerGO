package hub

import (
	"fmt"
	"sync"

	"pokerGO/internal/player"

	"github.com/gorilla/websocket"
)

func NewLobby() *Lobby {
	l := &Lobby{Players: make(map[string]*player.Player), Occupied: make(map[int]bool), PlayerCh: make(chan player.Player)}
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
	PlayerTurn chan *player.Player
	StartGame  chan struct{}
	//Conns   chan *websocket.Conn
}

// <ol hx-swap-oob=beforeend:#piece>    <li>%v</li></ol>

func (l *Lobby) LobbyWork() {
	fmt.Println("im in")
	for {
		select {
		case x := <-l.PlayerCh: // broadcoasting one to everyone
			for _, v := range l.Players {
				xs := x
				if xs != *v {
					xs = v.PrivateSend()
				}
				v.Conn.WriteJSON(xs)
				fmt.Println(v, "piskaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			}
		case <-l.StartGame:
			l.Game()
		}
	}
}

func (l *Lobby) Connhandle(plr *player.Player, conn *websocket.Conn) {
	fmt.Println("im in2")
	data := player.InterimPlayer{}
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
		err := plr.Conn.ReadJSON(&data)
		if err != nil {
			fmt.Println(err, "conn read error")
			plr.Conn = nil
			break
		}

		if l.Admin == plr && len(l.Occupied) >= 2 && data.IsGame == true {
			l.StartGame <- struct{}{}
		}
		plr.ChangeState(l.Occupied, data)
		fmt.Println(data, "pered v ch")

		l.PlayerCh <- *plr

	}
}
func (l *Lobby) Game() {
	//timer := time.NewTicker(time.Second * 1)
	PlayerBroadcast := make(chan player.Player)
	for {
		select {
		case pb := <-PlayerBroadcast: // broadcoasting one to everyone
			for _, v := range l.Players {
				pbs := pb
				if pbs != *v {
					pbs = v.PrivateSend()
				}
				v.Conn.WriteJSON(pbs)
				fmt.Println(v, "piskaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			}
			/*case tick := <-timer.C:
				timevalue := tick.Second()
				PlayerBroadcast <- l.Players[l.CurrentTurn].SendTimeValue(timevalue)
			case x1 := <-l.PlayerCh: // evaluating hand
				for _, v := range l.Players {
					xs := x
					if xs != *v {
						xs = v.PrivateSend()
					}
					v.Conn.WriteJSON(xs)
					fmt.Println(v, "piskaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
				}*/
		}
	}
}
