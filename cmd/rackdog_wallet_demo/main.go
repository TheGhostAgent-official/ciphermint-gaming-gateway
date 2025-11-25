package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type SignupRequest struct {
	Alias string `json:"alias"`
}

type WalletResponse struct {
	PlayerID string            `json:"player_id"`
	Alias    string            `json:"alias"`
	Balances map[string]uint64 `json:"balances"`
}

func main() {
	mux := http.NewServeMux()

	// In-memory demo state (for presentation)
	var current WalletResponse

	// POST /api/signup – simulate account creation + 100 RACKDOG bonus
	mux.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if req.Alias == "" {
			req.Alias = "GhostPlayer"
		}

		current = WalletResponse{
			PlayerID: "player_demo_001",
			Alias:    req.Alias,
			Balances: map[string]uint64{
				"RACKDOG": 100, // signup bonus
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(current)
	})

	// GET /api/wallet – return current wallet state
	mux.HandleFunc("/api/wallet", func(w http.ResponseWriter, r *http.Request) {
		if current.PlayerID == "" {
			http.Error(w, "no wallet yet", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(current)
	})

	// Serve the front-end
	fs := http.FileServer(http.Dir("web/rackdog_wallet_demo"))
	mux.Handle("/", fs)

	addr := ":8090"
	log.Printf("RackDog Gaming Wallet demo listening on %s ...", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
