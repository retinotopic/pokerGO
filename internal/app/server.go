package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/retinotopic/pokerGO/internal/auth"
	"github.com/retinotopic/pokerGO/internal/hub"
	"github.com/retinotopic/pokerGO/internal/player"
	htmpl "github.com/retinotopic/pokerGO/pkg/htmx"
	"github.com/retinotopic/pokerGO/pkg/randfuncs"

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
	http.HandleFunc("/lobby/", s.lobbyMW)
	http.HandleFunc("/", s.handleWS)
	http.Handle("/main", templ.Handler(htmpl.Pagemain()))
	//return http.ListenAndServeTLS(s.listenaddr, nil, nil,nil)
	return http.ListenAndServe(s.listenaddr, nil)
}

// rand.New(rand.NewSource(time.Now().Unix()))
func (s *Server) startgame(w http.ResponseWriter, r *http.Request) {
	strr := ""
	for _, ok := s.urllobby[strr]; !ok; {
		strr = randfuncs.RandomString(25, randfuncs.NewSource())
		s.str = strr
		s.urllobby[strr] = hub.NewLobby()

		go s.urllobby[strr].LobbyWork()
		break
	}
	http.Redirect(w, r, "/lobby/"+s.str, http.StatusFound)
	//w.Header().Set("HX-Redirect","/"+strr) //remember
}
func (s *Server) lobbyMW(w http.ResponseWriter, r *http.Request) {
	wsurl := r.URL.Path[7:]
	fmt.Println(wsurl, "///")
	if hub, ok := s.urllobby[wsurl]; ok {
		cookie, err := auth.ReadHashCookie(r, auth.Secretkey, r.Cookies())
		if err != nil {
			cookie = auth.WriteHashCookie(w, auth.Secretkey)
		}
		if _, ok := hub.Players[cookie.Value]; ok {
			if hub.Players[cookie.Value].Conn == nil {
				htmpl.Turner("/"+wsurl).Render(r.Context(), w)
			} else {
				htmpl.Refresh(strconv.Itoa(1)).Render(r.Context(), w)
				return
			}
		} else {
			hub.Lock()
			hub.Players[cookie.Value] = player.NewPlayer()
			hub.Unlock()
			htmpl.Turner("/"+wsurl).Render(r.Context(), w)
		}
	} else {
		http.Redirect(w, r, "/main", http.StatusNotFound)
	}

}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path[1:])
	if hub, ok := s.urllobby[r.URL.Path[1:]]; ok {
		cookie, err := auth.ReadHashCookie(r, auth.Secretkey, r.Cookies())
		if err != nil {
			log.Fatalln(err)
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalln(err)
		}
		go hub.Connhandle(hub.Players[cookie.Value], conn)

	} else {
		http.Redirect(w, r, "/main", http.StatusNotFound)
	}

}
