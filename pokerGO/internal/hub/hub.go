package hub

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func NewLobby() *Lobby {
	l := &Lobby{Recv: make(chan struct{}), Clients: [6]*websocket.Conn{}}
	return l
}

type Lobby struct {
	Recv    chan struct{}
	Clients [6]*websocket.Conn
	Counter int
}

func (l *Lobby) LobbyWork() {
	fmt.Println("im in")
	for {
		select {
		case <-l.Recv:
			fmt.Println("aboba")
			for i := 0; i < len(l.Clients); i++ {
				if l.Clients[i] != nil {
					l.Clients[i].WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("<div id=count>%v</div>", l.Counter)))
				}
			}
		}
	}
}
