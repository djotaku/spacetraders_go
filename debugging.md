```go
// test out Get Agent endpoint
	agent := api.GetAgent(token)
	fmt.Printf("AccountId: %s, Symbol: %s, HQ: %s, Credits: %d, Faction: %s\n", agent.AccountId, agent.Symbol, agent.Headquarters, agent.Credits, agent.StartingFaction)

	// Test out server status endpoint
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
	fmt.Printf("Announcements: \n")
	for _, announcement := range serverStatus.Announcements {
		fmt.Printf("\tTitle: %s\n", announcement["title"])
		fmt.Printf("\t\t %s\n", announcement["body"])
	}
	fmt.Printf("Links: \n")
	for _, link := range serverStatus.Links {
		fmt.Printf("\t%s: %s\n", link["name"], link["url"])
	}

	// Test out Factions Endpoints
	factionDetails := api.GetFaction(token, "COSMIC")
	fmt.Println("Faction Test:")
	fmt.Printf("\tFaction Symbol: %s\n", factionDetails.Symbol)
	fmt.Printf("\tFaction Name: %s\n", factionDetails.Name)
	fmt.Printf("\tFaction Description: %s\n", factionDetails.Description)
	fmt.Printf("\tFaction HQ Waypoint: %s\n", factionDetails.Headquarters)
	fmt.Println("\tFaction Traits:")
	for _, trait := range factionDetails.Traits {
		fmt.Printf("\t\tSymbol/Name: %s/%s\n", trait["symbol"], trait["name"])
		fmt.Printf("\t\tDescription: %s\n", trait["description"])
	}
	fmt.Printf("\tThis faction is recruiting? %v\n", factionDetails.IsRecruiting)

	allFactions := api.ListFactions(token, 10, 1)
	fmt.Println(allFactions.Factions[0].Description)
```