package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"pokerGO/internal/hub"
	rnds "pokerGO/pkg/Strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
)

type Server struct {
	listenaddr string
	urllobby   map[string]*hub.Lobby
	str        string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer(listenaddr string) *Server {
	return &Server{listenaddr: listenaddr, urllobby: make(map[string]*hub.Lobby)}
}

func (s *Server) Run() error {
	http.HandleFunc("/startgame", s.startgame)
	http.HandleFunc("/", s.handlelobby)
	http.HandleFunc("/lobby/", s.lobby)
	http.Handle("/main", templ.Handler(pagemain()))
	return http.ListenAndServe(s.listenaddr, nil)
}

// rand.New(rand.NewSource(time.Now().Unix()))
func (s *Server) startgame(w http.ResponseWriter, r *http.Request) {
	strr := ""
	fmt.Println("im here212")
	for _, ok := s.urllobby[strr]; !ok; {
		strr = rnds.RandomString(25, rand.New(rand.NewSource(time.Now().Unix())))
		s.str = strr
		s.urllobby[strr] = hub.NewLobby()
		go s.urllobby[strr].LobbyWork()
		break
	}
	http.Redirect(w, r, "/lobby/"+s.str, http.StatusFound)
	//w.Header().Set("HX-Redirect","/"+strr) //remember
}
func (s *Server) lobby(w http.ResponseWriter, r *http.Request) {
	wsurl := r.URL.Path[7:]
	fmt.Println(wsurl, "///")
	if s.urllobby[wsurl].Clients[5] == nil {
		turner("/"+wsurl).Render(r.Context(), w)
	} else {
		w.Write([]byte("the room is crowded"))
	}
}
func connhandle(l *hub.Lobby, conn *websocket.Conn) {
	fmt.Println("im in2")
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("aboba2")
		l.Counter++
		l.Recv <- struct{}{}
	}

}

var check bool

func (s *Server) handlelobby(w http.ResponseWriter, r *http.Request) {
	if hub, ok := s.urllobby[r.URL.Path[1:]]; ok {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		for i := 0; i < len(hub.Clients); i++ {
			if hub.Clients[i] == nil {
				hub.Clients[i] = conn
				go connhandle(hub, conn)
				check = true
				break
			}
		}
		if check == false {
			conn.Close()
			w.Write([]byte("lobby dont exist"))
		}
	} else {
		w.Write([]byte("lobby dont exist"))
	}

}
