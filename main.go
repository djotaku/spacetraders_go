package main

import (
	"fmt"
	"os"
	"strings"

	api "www.github.com/djotaku/spacetraders_go/spacetradersapi"
)

func main() {

	// get the auth token
	token_bytes, err := os.ReadFile("token")
	if err != nil {
		os.Exit(1)
	}
	token := string(token_bytes)
	token = strings.TrimRight(token, "\r\n")
	// fmt.Println(token) // jsut for now for debugging

	// test out registration
	//response, _ := api.SpaceTradersCommand(`{"symbol":"TestOtaku", "faction": "COSMIC" }`, "register", "", "post")
	//fmt.Print(response)
	// test out Get Agent endpoint
	agent := api.GetAgent(token)
	fmt.Printf("AccountId: %s, Symbol: %s, HQ: %s, Credits: %d, Faction: %s\n", agent.AccountId, agent.Symbol, agent.Headquarters, agent.Credits, agent.StartingFaction)
	serverStatus := api.GetStatus(token)
	fmt.Printf("##############################\nServer Status: \n")
	fmt.Printf("Status: %s\n", serverStatus.Status)
	fmt.Printf("Version: %s\n", serverStatus.Version)
	fmt.Printf("Reset Date: %s\n", serverStatus.ResetDate)
	fmt.Printf("Stats: \n")
	fmt.Printf("\t Agents: %d\n", serverStatus.Stats["agents"])
	fmt.Printf("\t Ships: %d\n", serverStatus.Stats["ships"])
	fmt.Printf("\t Systems: %d\n", serverStatus.Stats["systems"])
	fmt.Printf("\t Waypoints: %d\n", serverStatus.Stats["waypoints"])
	fmt.Printf("Leaderboards:\n")
	fmt.Printf("\t Most Credits:\n")
	for _, leaderboardParticipant := range serverStatus.Leaderboards.MostCredits {
		fmt.Printf("\t\tAgent Symbol: %s with %f credits.\n", leaderboardParticipant["agentSymbol"], leaderboardParticipant["credits"])
	}
	fmt.Printf("\n\t Most Charts:\n")
	for _, leaderboardParticipant := range serverStatus.Leaderboards.MostSubmittedCharts {
		fmt.Printf("\t\tAgent Symbol: %s with %f charts.\n", leaderboardParticipant["agentSymbol"], leaderboardParticipant["chartCount"])
	}
	fmt.Printf("Server Resets:\n")
	fmt.Printf("\t Next Date for Reset: %s\n", serverStatus.ServerResets["next"])
	fmt.Printf("\t Frequency of Resets: %s\n", serverStatus.ServerResets["frequency"])
}
