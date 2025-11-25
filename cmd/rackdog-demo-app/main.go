package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve static RackDog™ demo UI
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	addr := ":8090"
	log.Printf("RackDog™ Presentation App listening on http://localhost%s\n", addr)
	log.Println("Powered by CipherMint™ Gaming Gateway")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
