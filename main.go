package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"blockchess/internal/game"
	"blockchess/internal/websocket"

	"github.com/gorilla/mux"
)

func main() {
	var addr = flag.String("addr", ":8080", "http service address")
	flag.Parse()

	// Create game manager
	gameManager := game.NewManager()

	// Create WebSocket hub
	hub := websocket.NewHub(gameManager)
	go hub.Run()

	// Create router
	r := mux.NewRouter()

	// WebSocket endpoint
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r)
	})

	// Serve static files and handle client-side routing
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the path
		path := r.URL.Path

		// Check if it's a static asset (has file extension)
		if strings.Contains(path, ".") {
			// Serve static file
			http.FileServer(http.Dir("./dist/")).ServeHTTP(w, r)
			return
		}

		// For all other paths (including root and /game), serve index.html
		indexPath := filepath.Join("./dist", "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, indexPath)
	})

	log.Printf("Server starting on %s", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
