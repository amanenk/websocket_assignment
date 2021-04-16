package handlers

import (
	"github.com/fdistorted/websocket-practical/server/handlers/websocket"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ws", websocket.Get).Methods(http.MethodGet).Schemes("http")

	return r
}
