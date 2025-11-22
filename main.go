package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"

    "ciphermint-gaming-gateway/internal/player"
)

func main() {
    r := mux.NewRouter()

    // Health check
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("CipherMint Gaming Gateway is running"))
    }).Methods("GET")

    // Player endpoints
    r.HandleFunc("/v1/player/create", player.CreatePlayerHandler).Methods("POST")
    r.HandleFunc("/v1/player/{id}", player.GetPlayerHandler).Methods("GET")
    r.HandleFunc("/v1/player/{id}/earn", player.EarnTokenHandler).Methods("POST")

    log.Println("CipherMint Gaming Gateway listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}