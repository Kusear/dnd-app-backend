package handlers

import (
	"log/slog"
	"net/http"

	"dnd-backend-go/internal/client"

	"github.com/go-chi/chi/v5"
)

func RegisterHandlers(router *chi.Mux, hub *client.Hub, commandRouter *client.CommandRouter) {
	wsHandler := &WSConnectionHandler{}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler.Handle(w, r, hub, commandRouter)
	})
	slog.Info("WS connection handler registered")
}
