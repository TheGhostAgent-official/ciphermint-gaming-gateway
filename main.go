package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"ciphermint-gaming-gateway/internal/api"
)

func main() {
	r := mux.NewRouter()
	api.RegisterRoutes(r)

	log.Println("CipherMint Gaming Gateway listening on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
