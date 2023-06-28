package main

import (
	"fmt"

	api "www.github.com/djotaku/spacetraders_go/spacetradersapi"
)

func main() {
	// test out registration
	response, _ := api.SpaceTradersCommand(`{"symbol":"TestOtaku", "faction": "COSMIC" }`, "register", "", "post")
	fmt.Print(response)
}
