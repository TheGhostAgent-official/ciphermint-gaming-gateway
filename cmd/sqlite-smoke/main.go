package main

import (
	"ciphermint-gaming-gateway/internal/sqlstore"
	"fmt"
)

func main() {
	store, err := sqlstore.OpenDefault()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	fmt.Println("SQLite smoke test OK")
}
