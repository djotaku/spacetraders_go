package main

import (
	"fmt"

	"www.github.com/djotaku/spacetraders/api/api"
)

func main() {
	// test out registration
	response := api.SpaceTradersCommand(`{"symbol":"TestOtaku", "faction": "COSMIC" }`, "register", "", "get")
	fmt.Print(response)
}
