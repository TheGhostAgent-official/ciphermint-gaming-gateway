package main

import (
	"log"
	"net/http"

	"ciphermint-gaming-gateway/internal/api"
	"ciphermint-gaming-gateway/internal/sqlstore"
)

func main() {
	// Open the SQLite-backed store
	store, err := sqlstore.OpenDefault()
	if err != nil {
		log.Fatalf("open sqlite store: %v", err)
	}
	defer store.Close()

	// Build the HTTP router (health + game/player/balance API)
	handler := api.NewRouter(store)

	log.Println("CipherMint Gaming Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}