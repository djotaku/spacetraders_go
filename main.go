package main

import (
	"fmt"
	"os"

	api "www.github.com/djotaku/spacetraders_go/spacetradersapi"
)

func main() {

	// get the auth token
	token_bytes, err := os.ReadFile("token")
	if err != nil {
		os.Exit(1)
	}
	token := string(token_bytes)
	fmt.Println(token) // jsut for now for debugging

	// test out registration
	//response, _ := api.SpaceTradersCommand(`{"symbol":"TestOtaku", "faction": "COSMIC" }`, "register", "", "post")
	//fmt.Print(response)
	// test out Get Agent endpoint
	agent := api.GetAgent(token)
	fmt.Printf("AccountId: %s, Symbol: %s, HQ: %s, Credits: %d, Faction: ", agent.AccountId, agent.Symbol, agent.Credits, agent.StartingFaction)
}
