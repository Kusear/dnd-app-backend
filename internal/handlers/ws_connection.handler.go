package handlers

import (
	"net/http"

	"dnd-backend-go/internal/client"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSConnectionHandler struct {
}

func (h *WSConnectionHandler) Handle(w http.ResponseWriter, r *http.Request, hub *client.Hub, commandRouter *client.CommandRouter) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() {
		client.InitClient(ws, hub, commandRouter)
	}()
}
