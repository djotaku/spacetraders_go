// The API for interacting with Space Traders
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type apiResult struct {
	Data any
}

type Leaderboard struct {
	MostCredits         []map[string]any
	MostSubmittedCharts []map[string]any
}

type ServerStatus struct {
	Status        string
	Version       string
	ResetDate     string
	Description   string
	Stats         map[string]int64
	Leaderboards  Leaderboard
	ServerResets  map[string]string
	Announcements []map[string]string
	Links         []map[string]string
}

type Agent struct {
	AccountId       string
	Symbol          string
	Headquarters    string
	Credits         int64
	StartingFaction string
}

type Contract struct {
	Id               string
	FactionSymbol    string
	Type             string
	Terms            []map[string]string
	Accepted         bool
	Fulfilled        bool
	DeadlineToAccept string
}

type FactionList struct {
	Factions []Faction `json:"data"`
	Meta     map[string]int
}

type Faction struct {
	Symbol       string
	Name         string
	Description  string
	Headquarters string
	Traits       []map[string]string
	IsRecruiting bool
}

// continue with ship

// SpaceTradersCommand sends a command to the Space Traders API
func SpaceTradersCommand(parameters string, endpoint string, authToken string, httpVerb string) (string, error) {
	apiURLBase := "https://api.spacetraders.io/v2/"
	fullURL := apiURLBase + endpoint
	jsonBody := []byte(parameters)
	bodyReader := bytes.NewReader(jsonBody)
	bearer := "Bearer " + authToken

	if httpVerb == "get" {
		request, err := http.NewRequest(http.MethodGet, fullURL, nil)
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", bearer)
		if err != nil {
			return "Error accessing URL", err
		}
		result, err := http.DefaultClient.Do(request)
		if err != nil {
			return "Error accessing URL", err
		}

		resultBody, err := io.ReadAll(result.Body)
		if result.StatusCode > 299 {
			fmt.Printf("Response failed with status code: %d and \nbody: %s\n", result.StatusCode, resultBody)
		}
		if err != nil {
			return "Error reading response", err
		}
		return string(resultBody), err

	} else if httpVerb == "post" {
		request, err := http.NewRequest(http.MethodPost, fullURL, bodyReader)
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", bearer)
		if err != nil {
			return "Error accessing URL", err
		}
		result, err := http.DefaultClient.Do(request)
		resultBody, err := io.ReadAll(result.Body)
		if result.StatusCode > 299 {
			fmt.Printf("Response failed with status code: %d and \nbody: %s\n", result.StatusCode, resultBody)
		}
		if err != nil {
			return "Error reading response", err
		}

		return string(resultBody), err
	}
	return "Something went wrong", fmt.Errorf("error")
}

func GetAgent(authToken string) Agent {
	agentResult, err := SpaceTradersCommand("", "my/agent", authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Agent endpoint. Error: ", err)
	}
	agent := &Agent{}
	json.Unmarshal([]byte(agentResult), &apiResult{agent})
	return *agent
}

func GetStatus(authToken string) ServerStatus {
	statusResult, err := SpaceTradersCommand("", "", authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Status endpoint. Error: ", err)
	}
	//fmt.Println(statusResult)
	var status ServerStatus
	json.Unmarshal([]byte(statusResult), &status)
	return status
}

func ListContracts(authToken string, limit int, page int) []Contract {
	parameters := fmt.Sprintf(`{"limit": %d, "page": %d}`, limit, page)
	contractResult, err := SpaceTradersCommand(parameters, "my/contracts", authToken, "get")
	if err != nil {
		fmt.Println("Error accessing List Contracts endpoint. Error: ", err)
	}
	fmt.Printf("Contract debugging from API: %s", contractResult)
	return nil
}

func GetContract(authToken string, contractID string) Contract {
	parameters := fmt.Sprintf(`{"contractID": %s}`, contractID)
	contractURL := fmt.Sprintf("my/contracts/%s", contractID)
	factionResults, err := SpaceTradersCommand(parameters, contractURL, authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Contract endpoint. Error: ", err)
	}
	contract := &Contract{}
	json.Unmarshal([]byte(factionResults), &apiResult{contract})
	return *contract
}

func GetFaction(authToken string, factionSymbol string) Faction {
	parameters := fmt.Sprintf(`{"factionSymbol": %s}`, factionSymbol)
	factionURL := fmt.Sprintf("factions/%s", factionSymbol)
	factionResults, err := SpaceTradersCommand(parameters, factionURL, authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Faction endpoint. Error: ", err)
	}
	faction := &Faction{}
	json.Unmarshal([]byte(factionResults), &apiResult{faction})
	return *faction
}

func ListFactions(authToken string, limit int, page int) FactionList {
	parameters := fmt.Sprintf(`{"limit": %d, "page": %d}`, limit, page)
	factionResult, err := SpaceTradersCommand(parameters, "factions", authToken, "get")
	if err != nil {
		fmt.Println("Error accessing List Contracts endpoint. Error: ", err)
	}
	//fmt.Printf("List Factions debugging from API: %s", factionResult)
	var factionList FactionList
	json.Unmarshal([]byte(factionResult), &factionList)
	return factionList
}
