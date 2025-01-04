// internal/api/api.go
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"opggvisualizer/internal/client"
	"opggvisualizer/internal/config"
)

var server *http.Server // Global reference to the server

func newServer() *http.Server {
	if server == nil {
		cfg := config.GetConfig()
		server = &http.Server{Addr: ":" + cfg.APIServer.Port}
	}
	return server
}

func GetServer() *http.Server {
	if server == nil {
		server = newServer()
	}
	return server
}

func Start(ctx context.Context) {
	http.HandleFunc("/refresh", handleRefresh)
	http.HandleFunc("/health", handleHealth)

	GetServer() // Initialize the server if necessary
	server.Handler = http.DefaultServeMux
	// Start the server in a goroutine
	go func() {
		log.Printf("Starting API server at: %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API server failed to start: %v", err)
		}
	}()

	// Wait for the context to be cancelled
	<-ctx.Done()

	// Create a new context with a timeout to allow the server to shut down gracefully
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed:%+v", err)
	}
	log.Println("server exited properly")
}

func Stop(ctx context.Context) error {
	log.Println("Stopping API server...")
	return server.Shutdown(ctx)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	go func() {
		if err := client.FetchAndStoreChampionData(); err != nil {
			log.Printf("Error fetching and storing champion data: %v", err)
		} else {
			log.Println("Data fetching and insertion completed successfully")
		}
		if err := client.FetchAndStoreGameData(); err != nil {
			log.Printf("Error fetching and storing game data: %v", err)
		} else {
			log.Println("Data fetching and insertion completed successfully")
		}
	}()

	response := map[string]string{
		"status":  "Data refresh initiated",
		"message": "Your data refresh request is being processed.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
