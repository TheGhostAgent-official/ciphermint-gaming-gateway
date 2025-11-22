package api

import (
    "ciphermint-gaming-gateway/internal/models"
    "ciphermint-gaming-gateway/internal/store"
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/v1/player", createPlayer).Methods("POST")
    router.HandleFunc("/v1/player/{id}", getPlayer).Methods("GET")
    router.HandleFunc("/v1/player/{id}/earn", rewardTokens).Methods("POST")
    router.HandleFunc("/v1/player/{id}/spend", spendTokens).Methods("POST")

    return router
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
    var req models.CreatePlayerRequest
    json.NewDecoder(r.Body).Decode(&req)

    player := store.CreatePlayer(req.PlayerID, req.Alias)
    json.NewEncoder(w).Encode(player)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    player := store.GetPlayer(models.PlayerID(id))
    json.NewEncoder(w).Encode(player)
}

func rewardTokens(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    var req models.EarnRequest
    json.NewDecoder(r.Body).Decode(&req)

    updated := store.Earn(models.PlayerID(id), req.Token, req.Amount)
    json.NewEncoder(w).Encode(updated)
}

func spendTokens(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    var req models.SpendRequest
    json.NewDecoder(r.Body).Decode(&req)

    updated := store.Spend(models.PlayerID(id), req.Token, req.Amount)
    json.NewEncoder(w).Encode(updated)
}