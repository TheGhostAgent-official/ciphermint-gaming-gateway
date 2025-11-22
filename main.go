package main

import (
	"log"
	"net/http"

	"ciphermint-gaming-gateway/internal/api"
	"ciphermint-gaming-gateway/internal/store"
)

func main() {
	s := store.NewMemoryStore()
	h := api.NewHandler(s)
	router := api.NewRouter(h)

	log.Println("CipherMint Gaming Gateway listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}