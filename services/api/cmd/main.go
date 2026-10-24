package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/booking-network/bn-near/services/api/internal/chat"
	"github.com/booking-network/bn-near/services/api/internal/router"
	"github.com/booking-network/bn-near/services/api/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("=======================================================")
	log.Printf("BN NEAR — GLOBAL HYPERLOCAL MARKETPLACE API")
	log.Printf("Booking Network Co., Ltd. (Master Order BN-NEAR-MASTER-001)")
	log.Printf("=======================================================")

	// Initialize Storage
	store := storage.NewMemoryStore()
	log.Printf("[Storage] Initialized Storage Engine with seed fixtures (Users: %d, Neighborhoods: %d, Listings: %d)",
		len(store.Users), len(store.Neighborhoods), len(store.Listings))

	// Initialize WebSocket Hub
	hub := chat.NewHub(store)
	go hub.Run()
	log.Printf("[WebSocket] Real-Time WebSocket Hub running on /api/v1/ws")

	// Initialize Router
	routes := router.NewRouter(store, hub)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      routes,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[Server] API Gateway listening on http://0.0.0.0:%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] Failed to listen: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("[Server] Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[Server] Forced shutdown error: %v", err)
	}
	fmt.Println("[Server] Server stopped.")
}

