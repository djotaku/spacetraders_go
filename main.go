package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	api "www.github.com/djotaku/spacetraders_go/spacetradersapi"
)

func updateAgentArea(thisContentBox *fyne.Container) *fyne.Container {
	factionListSelection := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
		log.Println("Select set to", value)
	})
	callsignInput := widget.NewEntry()
	callsignInput.SetPlaceHolder("Enter your Call Sign")
	newAgentButton := widget.NewButton("Get a New Agent", func() { api.RegisterNewAgent(factionListSelection.Selected, callsignInput.Text) })
	thisContentBox = container.New(layout.NewHBoxLayout(), factionListSelection, newAgentButton)
	thisContentBox.Refresh()
	return thisContentBox
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
	contentBox := container.New(layout.NewHBoxLayout(), tempContentText)
	tempTopRowText := canvas.NewText("Top Row", color.White)
	topRowBox := container.New(layout.NewHBoxLayout(), tempTopRowText)
	agentButton := widget.NewButton("Agents", func() { contentBox = updateAgentArea(contentBox); contentBox.Refresh() })
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
