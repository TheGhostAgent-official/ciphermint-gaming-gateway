package main

import (
	"log"
	"net/http"

	"ciphermint-gaming-gateway/internal/api"
	"ciphermint-gaming-gateway/internal/store"
)

func main() {
	// Initialize in-memory database
	db := store.NewMemoryStore()

	// Create HTTP handler set
	handler := api.NewHandler(db)

	// Build router
	router := api.NewRouter(handler)

	log.Println("CipherMint Gaming Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}