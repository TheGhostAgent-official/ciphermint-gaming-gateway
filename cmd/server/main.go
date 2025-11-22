package server

import (
    "ciphermint-gaming-gateway/internal/api"
    "fmt"
    "log"
    "net/http"
)

func Start() {
    router := api.NewRouter()

    fmt.Println("CypherMint Gaming Gateway running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}