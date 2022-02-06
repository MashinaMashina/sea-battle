package service

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"seabattle/internal/interfaces"
)

type GameServer struct {
	pairs []interfaces.GamePair
}

func NewServer() *GameServer {
	pairs := make([]interfaces.GamePair, 0)

	return &GameServer{pairs: pairs}
}

func (gs GameServer) RegisterHandler(router *mux.Router)  {
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Traceln("Create websocket connection with user")

		upgrader := websocket.Upgrader{}
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			panic(err)
		}

		c := NewClient(ws)

		if !gs.ConnectPair(c) {
			p := NewPair()
			p.AddClient(c)

			gs.pairs = append(gs.pairs, p)
		}
	})
}

func (gs GameServer) ConnectPair(c interfaces.GameClient) bool {
	for _, p := range gs.pairs {
		if p.IsFree() {
			if err := p.AddClient(c); err != nil {
				panic(err)
				return false
			}
			return true
		}
	}

	return false
}

