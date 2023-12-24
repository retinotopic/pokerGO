package hub

import (
	"fmt"
	"sync"

	"pokerGO/internal/player"

	"github.com/gorilla/websocket"
)

func NewLobby() *Lobby {
	l := &Lobby{Players: make(map[string]*player.Player), Occupied: make(map[int]bool), PlayerCh: make(chan *player.Player)}
	return l
}

type Lobby struct {
	Players  map[string]*player.Player
	Admin    *player.Player
	Occupied map[int]bool
	sync.Mutex
	sync.Once
	PlayerCh   chan *player.Player
	PlayerTurn chan *player.Player
	//Conns   chan *websocket.Conn
}

// <ol hx-swap-oob=beforeend:#piece>    <li>%v</li></ol>

func (l *Lobby) LobbyWork() {
	fmt.Println("im in")
	for {
		select {
		case x := <-l.PlayerCh:
			for _, v := range l.Players {
				v.Conn.WriteJSON(x)
				fmt.Println(v, "piskaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			}
		}
	}
}

var data = map[string]interface{}{
	"Name":     "",
	"Stack":    int(0),
	"IsActive": false,
	"Place":    int(0),
	"IsAdmin":  false,
	"IsFold":   false,
	"IsBet":    false,
}

func (l *Lobby) Connhandle(player *player.Player, conn *websocket.Conn) {
	fmt.Println("im in2")
	l.Do(func() {
		l.Admin = player

	})
	defer func() {
		fmt.Println("rip connection")

	}()
	player.Conn = conn
	for _, v := range l.Players {
		data["Name"] = v.Name
		data["Stack"] = v.Bankroll
		data["IsActive"] = v.IsActive
		data["Place"] = v.Place
		err := player.Conn.WriteJSON(data)
		if err != nil {
			fmt.Println(err, "WriteJSON start")
		}
		fmt.Println(data, "start")
	}

	for {
		err := player.Conn.ReadJSON(&data)
		if err != nil {
			fmt.Println(err, "conn read error")
			player.Conn = nil
			break
		}
		if l.Admin == player && len(l.Occupied) >= 2 && data["IsGame"].(bool) == true {
			data["IsAdmin"] = true
		} else if player.IsActive == false && l.Occupied[int(data["Place"].(float64))] == false && data["IsActive"].(bool) == true {
			player.Name = data["Name"].(string)
			player.Bankroll = int(data["Stack"].(float64))
			player.IsActive = data["IsActive"].(bool)
			player.Place = int(data["Place"].(float64))
			l.Occupied[int(data["Place"].(float64))] = true
		} else if player.IsActive == true && data["IsActive"].(bool) == false {
			player.IsActive = data["IsActive"].(bool)
			delete(l.Occupied, player.Place)
		}
		fmt.Println(data, "pered v ch")
		l.PlayerCh <- player
	}
}
