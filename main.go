package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

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
	appWindow.SetContent(widget.NewLabel("hello world!"))
	appWindow.ShowAndRun()
}
