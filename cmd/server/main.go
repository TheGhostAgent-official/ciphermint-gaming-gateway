package server

import (
    "ciphermint-gaming-gateway/internal/api"
    "ciphermint-gaming-gateway/internal/store"
    "log"
    "net/http"
)

// Start bootstraps the in-memory store, builds the router,
// and starts the HTTP server on :8080.
func Start() {
    s := store.NewMemoryStore()
    r := api.NewRouter(s)

    log.Println("CipherMint Gaming Gateway listening on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}