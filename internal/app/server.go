package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"pokerGO/internal/hub"
	"pokerGO/internal/player"
	rnds "pokerGO/pkg/Strings"
	htmpl "pokerGO/pkg/htmx"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
)

type Server struct {
	listenaddr string
	urllobby   map[string]*hub.Lobby
	str        string
}

var key = []byte("cbaTd3Dx9_dfknwPsc5T0rQMx34SvJJf5xvxf7nab")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer(listenaddr string) *Server {
	return &Server{listenaddr: listenaddr, urllobby: make(map[string]*hub.Lobby)}
}

func (s *Server) Run() error {
	http.HandleFunc("/startgame", s.startgame)
	http.HandleFunc("/", s.handleWS)
	http.HandleFunc("/lobby/", s.lobby)
	http.Handle("/main", templ.Handler(htmpl.Pagemain()))
	//return http.ListenAndServeTLS(s.listenaddr, nil, nil,nil)
	return http.ListenAndServe(s.listenaddr, nil)
}

// rand.New(rand.NewSource(time.Now().Unix()))
func (s *Server) startgame(w http.ResponseWriter, r *http.Request) {
	strr := ""
	for _, ok := s.urllobby[strr]; !ok; {
		strr = rnds.RandomString(25, rnds.NewSource())
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
	if hub, ok := s.urllobby[wsurl]; ok {
		cookie, err := ReadHashCookie(r, key, r.Cookies())
		if err != nil {
			cookie = WriteHashCookie(w, key)
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
	if hub, ok := s.urllobby[r.URL.Path[1:]]; ok {
		cookie, err := ReadHashCookie(r, key, r.Cookies())
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

func WriteHashCookie(w http.ResponseWriter, key []byte) *http.Cookie {
	mac := hmac.New(sha256.New, key)
	r0 := rand.New(rand.NewSource(time.Now().Unix()))
	time.Sleep(time.Millisecond * 25)
	r1 := rand.New(rand.NewSource(time.Now().Unix()))
	r1.Seed(time.Now().UnixNano())
	cookie := &http.Cookie{Name: rnds.RandomString(15, r0), Value: rnds.RandomString(20, r1), Secure: true, Path: "/"}
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	cookie.Value = cookie.Value + signature
	http.SetCookie(w, cookie)
	return cookie
}
func ReadHashCookie(r *http.Request, key []byte, cookies []*http.Cookie) (*http.Cookie, error) {
	if len(cookies) == 0 {
		return nil, errors.New("zero cookies")
	}
	c := cookies[0]
	name := c.Name
	valueHash := c.Value
	signature := valueHash[20:]
	value := valueHash[:20]

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, errors.New("ValidationErr")
	}
	return c, nil
}
