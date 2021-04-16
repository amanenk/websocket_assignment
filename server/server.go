package server

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/server/config"
	"github.com/fdistorted/websocket-practical/server/handlers"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/server/websocket/broadcast"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

func Start() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to laod config: %v", err)
	}

	err = logger.Load() //todo maybe add some config to loader
	if err != nil {
		log.Fatalf("Failed to laod logger: %v", err)
	}

	//storage := clients.NewStorage()
	//start broadcaster
	broadcast.InitBroadcast()

	server := &http.Server{
		Addr:         cfg.ListenUrl,
		Handler:      handlers.NewRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Get().Info("Listening...", zap.String("listen_url", cfg.ListenUrl))
	err = server.ListenAndServe()
	if err != nil {
		// logger.Get().Error("Failed to initialize HTTP server", zap.Error(err))
		fmt.Println("failed to start server")
		os.Exit(1)
	}
}
