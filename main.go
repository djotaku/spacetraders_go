package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	api "www.github.com/djotaku/spacetraders_go/spacetradersapi"
)

func updateAgentArea(thisContentBox *fyne.Container) {
	factionListSelection := widget.NewSelect([]string{"COSMIC", "VOID", "GALACTIC", "QUANTUM", "DOMINION", "ASTRO",
		"CORSAIRS", "OBSIDIAN", "AEGIS", "UNITED", "SOLITARY", "COBALT", "OMEGA", "ECHO", "LORDS", "CULT", "ANCIENTS", "SHADOW", "ETHEREAL"}, func(value string) {
		log.Println("Select set to", value)
	})
	callsignInput := widget.NewEntry()
	callsignInput.SetPlaceHolder("Enter your Call Sign")
	newAgentButton := widget.NewButton("Get a New Agent", func() { api.RegisterNewAgent(factionListSelection.Selected, callsignInput.Text) })
	thisContentBox.RemoveAll()
	thisContentBox.Add(factionListSelection)
	thisContentBox.Add(callsignInput)
	thisContentBox.Add(newAgentButton)
	thisContentBox.Refresh()
}

func main() {

	// get the auth token
	//token_bytes, err := os.ReadFile("token")
	//if err != nil {
	//		os.Exit(1)
	//}
	//token := string(token_bytes)
	//token = strings.TrimRight(token, "\r\n")
	// fmt.Println(token) // jsut for now for debugging

	// test out registration
	//api.RegisterNewAgent("COSMIC", "TESTOTAKU18")
	//fmt.Println(token)
	//GUI
	spaceTradersApp := app.New()
	appWindow := spaceTradersApp.NewWindow("Space Traders")
	tempContentText := widget.NewLabel("Click a button to change")
	contentBox := container.New(layout.NewVBoxLayout(), tempContentText)
	agentName := widget.NewLabel("Agent: ")
	agentNameValue := widget.NewLabel("") // needs data binding
	money := widget.NewLabel("Money: ")
	moneyValue := widget.NewLabel("") // needs data binding
	numberOfShips := widget.NewLabel("Number of Ships: ")
	numberOfShipsValue := widget.NewLabel("") // needs data binding
	activeContracts := widget.NewLabel("Active Contracts: ")
	activeContractsValue := widget.NewLabel("") // needs data binding
	timeToNextMove := widget.NewLabel("Time to Next Move: ")
	timeToNextMoveValue := widget.NewLabel("") // needs data binding
	topRowBox := container.New(layout.NewHBoxLayout(), agentName, agentNameValue, money, moneyValue, numberOfShips, numberOfShipsValue, activeContracts, activeContractsValue, timeToNextMove, timeToNextMoveValue)
	agentButton := widget.NewButton("Agents", func() { updateAgentArea(contentBox) })
	shipButton := widget.NewButton("Ships", func() { log.Println("Ships") })
	contractButton := widget.NewButton("Contracts", func() { log.Println("Contracts") })
	waypointsButton := widget.NewButton("Waypoints", func() { log.Println("Waypoints") })
	currentLocationButton := widget.NewButton("Current Location", func() { log.Println("CurrentLocation") })
	automationsButton := widget.NewButton("Automations", func() { log.Println("Automations") })
	menuBox := container.New(layout.NewVBoxLayout(), agentButton, shipButton, contractButton, waypointsButton, currentLocationButton, automationsButton)
	treeAndContentBox := container.New(layout.NewHBoxLayout(), menuBox, contentBox)
	mainContent := container.New(layout.NewVBoxLayout(), topRowBox, treeAndContentBox)
	appWindow.SetContent(mainContent)
	appWindow.ShowAndRun()
}
