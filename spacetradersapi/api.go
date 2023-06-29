// The API for interacting with Space Traders
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type apiResponse struct {
	Data Agent `json:"data"`
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
	var response apiResponse
	json.Unmarshal([]byte(agentResult), &response)
	return response.Data
}
