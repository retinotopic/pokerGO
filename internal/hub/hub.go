package hub

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"pokerGO/internal/player"
	rnds "pokerGO/pkg/Strings"
	htmpl "pokerGO/pkg/htmx"

	"github.com/gorilla/websocket"
)

func NewLobby() *Lobby {
	l := &Lobby{Recv: make(chan struct{}), Players: make(map[string]*player.Player)}
	return l
}

type Lobby struct {
	Recv    chan struct{}
	Players map[string]*player.Player
	Admin   *player.Player
	sync.Mutex
	sync.Once
	//Conns   chan *websocket.Conn
}

var button = []byte("<form id=count ws-send>\n    <button>\n        Take seat\n    </button>\n</form>")
var form = []byte("<form id=form name=Stack name=Name ws-send>\n    <input id=Name name=Name>Enter your name</input>\n    <input id=Stack name=Stack>Enter your wished stack</input>\n    <button type=submit>Send data</button>\n</form>")

// <ol hx-swap-oob=beforeend:#piece>    <li>%v</li></ol>

func (l *Lobby) LobbyWork() {
	fmt.Println("im in")

	//var startgame sync.Once
	//startgame.Do()
	/*for {
		for _, v := range l.Players {
			v.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<div> id=count %v</div>", l.Counter)))

		}
	}*/
}
func (l *Lobby) Connhandle(player *player.Player, conn *websocket.Conn) {
	fmt.Println("im in2")
	l.Do(func() {
		l.Admin = player
		go func() {
			fmt.Println("im in3")
			l.Connhandle(player, conn)
		}()
	})

	defer func() {
		fmt.Println("rip connection")

	}()
	player.Conn = conn
	//player.Conn.SetReadDeadline(time.Now().Add(time.Second * 20))
	player.Conn.SetPongHandler(func(string) error {
		fmt.Println("ping")
		return nil
	})

	var jsonString map[string]interface{}
	for {
		err := player.Conn.ReadJSON(&jsonString)

		if err != nil {
			fmt.Println(err, "conn read error")
			player.Conn = nil
			break
		}
		fmt.Println(jsonString, "json")
		byteBuff := bytes.NewBuffer(make([]byte, 0))
		err = htmpl.Form().Render(context.Background(), byteBuff)
		if err != nil {
			fmt.Println(err, "render in bytes errir")
		}
		if val, ok := jsonString["Name"]; ok {
			fmt.Println(val)
		}
		err = player.Conn.WriteMessage(websocket.TextMessage, byteBuff.Bytes())
		if err != nil {
			fmt.Println(err, "conn writer render error")
		}

	}
}
func (l *Lobby) Game() {
	for {
		select {
		case <-l.Recv:
			rnd := len(l.Players)
			rnd2 := rnds.NewSource().Intn(rnd)
			i := 0
			for _, v := range l.Players {
				if i == rnd2 {
					v.Conn.WriteMessage(websocket.TextMessage, button)
				} else {
					v.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<div> id=count %v</div>", 7)))
				}
				i++
			}
		}
	}
}
