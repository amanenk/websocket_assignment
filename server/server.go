package server

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/server/broadcast"
	"github.com/fdistorted/websocket-practical/server/config"
	"github.com/fdistorted/websocket-practical/server/handlers"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func Start() {
	cfg := config.Load()

	err := logger.Load() //todo maybe add some config to loader

	if err != nil {
		log.Fatalf("Failed to laod logger: %s", err.Error())
	}

	server := &http.Server{
		Addr:    cfg.ListenUrl,
		Handler: handlers.NewRouter(),
	}

	logger.Get().Info("Listening...", zap.String("listen_url", cfg.ListenUrl))

	//start broadcaster
	broadcast.InitBroadcast()

	err = server.ListenAndServe()
	if err != nil {
		// logger.Get().Error("Failed to initialize HTTP server", zap.Error(err))
		fmt.Println("failed to start server")
		os.Exit(1)
	}
}
