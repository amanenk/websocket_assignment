package handlers

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/server/handlers/ws"
	"github.com/fdistorted/websocket-practical/server/websocket/storage"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func NewRouter(storage *storage.Storage) *mux.Router {
	r := mux.NewRouter()

	wh := ws.NewWebsocketHandler(storage)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "hello there %d", time.Now().Unix())
	}).Methods(http.MethodGet).Schemes("http")

	r.HandleFunc("/ws", wh.Get).Methods(http.MethodGet).Schemes("http")

	return r
}
