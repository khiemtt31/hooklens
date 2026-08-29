package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/internal/health"
	"server/internal/httpapi"
	"server/internal/inbox"
	"server/internal/ping"
	"server/internal/server"
	"server/internal/storage"
)

const version = "0.1.0"

func main() {
	database, err := storage.Open()
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}
	defer database.Close()

	healthService := health.NewService(
		"hooklens",
		version,
	)

	healthHandler := httpapi.NewHealthHandler(
		healthService,
	)

	pingService := ping.NewService()

	pingHandler := httpapi.NewPingHandler(
		pingService,
	)

	inboxRepository := inbox.NewRepository(database)
	inboxService := inbox.NewService(inboxRepository)
	inboxHandler := httpapi.NewInboxHandler(inboxService)

	router := httpapi.NewRouter(
		healthHandler,
		pingHandler,
		inboxHandler,
	)

	httpServer := server.New(
		":8080",
		router,
	)

	errorChannel := make(chan error, 1)

	go func() {
		log.Println("HookLens API running on http://localhost:8080")

		errorChannel <- httpServer.Start()
	}()

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(
		signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-errorChannel:
		if err != nil {
			log.Fatalf("server failed: %v", err)
		}

	case sig := <-signalChannel:
		log.Printf("received signal: %s", sig)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)

		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("server shutdown failed: %v", err)
		}

		log.Println("server stopped")
	}
}
