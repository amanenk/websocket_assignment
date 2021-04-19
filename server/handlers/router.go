package handlers

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/server/handlers/ws"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "hello there %d", time.Now().Unix())
	}).Methods(http.MethodGet).Schemes("http")
	r.HandleFunc("/ws", ws.Get).Methods(http.MethodGet).Schemes("http")

	return r
}
