package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Player represents a player in the gaming system
type Player struct {
	PlayerID string `json:"player_id"`
	Alias    string `json:"alias"`
}

// Reward represents earned in-game value paid in CipherMint tokens
type Reward struct {
	Token  string  `json:"token"`
	Amount float64 `json:"amount"`
	Source string  `json:"source"`
}

func main() {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("CipherMint Gaming Gateway OK"))
	})

	// Register a new player
	router.HandleFunc("/v1/player", func(w http.ResponseWriter, r *http.Request) {
		var p Player
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "player_registered",
			"player": p,
		})
	}).Methods("POST")

	// Reward a player with tokenized earnings
	router.HandleFunc("/v1/player/{id}/earn", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		playerID := vars["id"]

		var reward Reward
		if err := json.NewDecoder(r.Body).Decode(&reward); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "reward_issued",
			"player_id": playerID,
			"reward":    reward,
		})
	}).Methods("POST")

	// Logging + server start (UPDATED TO :8081)
	log.Println("CipherMint Gaming Gateway listening on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}
}