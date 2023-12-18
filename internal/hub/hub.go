package hub

import (
	"fmt"
	"sync"

	"pokerGO/internal/player"

	"github.com/gorilla/websocket"
)

func NewLobby() *Lobby {
	l := &Lobby{Players: make(map[string]*player.Player)}
	return l
}

type Lobby struct {
	Players  map[string]*player.Player
	Admin    *player.Player
	Occupied map[int]bool
	sync.Mutex
	sync.Once
	//Conns   chan *websocket.Conn
}

var button = []byte("<form id=count ws-send>\n    <button>\n        Take seat\n    </button>\n</form>")
var form = []byte("<form id=form name=Stack name=Name ws-send>\n    <input id=Name name=Name>Enter your name</input>\n    <input id=Stack name=Stack>Enter your wished stack</input>\n    <button type=submit>Send data</button>\n</form>")

// <ol hx-swap-oob=beforeend:#piece>    <li>%v</li></ol>
type tempPlayer struct {
	Name     string `json:"Name"`
	Bankroll int    `json:"Stack"`
	IsActive bool   `json:"IsActive"`
	Place    int    `json:"Place"`
}

func (l *Lobby) LobbyWork() {
	fmt.Println("im in")
	/*for {
		select {
		case <-l.Recv:
			fmt.Println("im in3")
			l.LobbyWork2()
		default:
			fmt.Println("im in4")
			time.Sleep(time.Second)
		}
		for _, v := range l.Players {
			v.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<div> id=count %v</div>", l.Counter)))

		}
	}*/
}

var data = map[string]interface{}{
	"Name":     "",
	"Stack":    0,
	"IsActive": false,
	"Place":    0,
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
	for {
		err := player.Conn.ReadJSON(&data)
		if err != nil {
			fmt.Println(err, "conn read error")
			player.Conn = nil
			break
		}
		if player.IsActive == false && l.Occupied[data["Place"].(int)] == false {
			player.Name = data["Name"].(string)
			player.Bankroll = data["Stack"].(int)
			player.IsActive = data["IsActive"].(bool)
			player.Place = data["Place"].(int)
			l.Occupied[data["Place"].(int)] = true
		}

		fmt.Println(data)
		err = player.Conn.WriteMessage(websocket.TextMessage, form)
		if err != nil {
			fmt.Println(err, "conn writer render error")
		}
	}
}
func (l *Lobby) Game() {
	for {
		select {}
	}
}
