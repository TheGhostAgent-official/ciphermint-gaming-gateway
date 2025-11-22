package store

import (
    "ciphermint-gaming-gateway/internal/models"
)

var players = map[models.PlayerID]*models.Player{}

func CreatePlayer(id models.PlayerID, alias string) *models.Player {
    p := &models.Player{
        ID:       id,
        Aliases:  []string{alias},
        Balances: map[models.TokenSymbol]int64{},
    }
    players[id] = p
    return p
}

func GetPlayer(id models.PlayerID) *models.Player {
    return players[id]
}

func Earn(id models.PlayerID, token models.TokenSymbol, amount int64) *models.Player {
    p := players[id]
    if p.Balances == nil {
        p.Balances = map[models.TokenSymbol]int64{}
    }
    p.Balances[token] += amount
    return p
}

func Spend(id models.PlayerID, token models.TokenSymbol, amount int64) *models.Player {
    p := players[id]
    if p.Balances[token] >= amount {
        p.Balances[token] -= amount
    }
    return p
}