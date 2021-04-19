package server

import (
	"fmt"
	"github.com/fdistorted/websocket-practical/filelimits"
	"github.com/fdistorted/websocket-practical/server/config"
	"github.com/fdistorted/websocket-practical/server/handlers"
	logger "github.com/fdistorted/websocket-practical/server/loggger"
	"github.com/fdistorted/websocket-practical/websocket/broadcast"
	storage2 "github.com/fdistorted/websocket-practical/websocket/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

func Start() {
	// prepare system to run it should be configured in deployment script
	filelimits.MaxOpenFiles()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to laod config: %v", err)
	}

	//log.Printf("config loaded: %+v\n", cfg)
	//log.Printf("port variable: %+v\n", os.Getenv("PORT"))

	err = logger.Load() //todo maybe add some config to loader
	if err != nil {
		log.Fatalf("Failed to laod logger: %v", err)
	}

	storage := storage2.NewStorage()
	//start broadcaster
	broadcast.InitBroadcast(storage)

	addr := fmt.Sprintf(":%d", cfg.Port)

	server := &http.Server{
		Addr:         addr,
		Handler:      handlers.NewRouter(storage),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Get().Info("Listening...", zap.String("listen_url", addr))
	err = server.ListenAndServe()
	if err != nil {
		// logger.Get().Error("Failed to initialize HTTP server", zap.Error(err))
		fmt.Println("failed to start server")
		os.Exit(1)
	}
}
