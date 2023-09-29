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

func GetAgent(authToken string) Agent {
	agentResult, err := SpaceTradersCommand("", "my/agent", authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Agent endpoint. Error: ", err)
	}
	agent := &Agent{}
	json.Unmarshal([]byte(agentResult), &apiResult{agent})
	return *agent
}
func RegisterNewAgent(faction string, agentName string) string {
	parameters := fmt.Sprintf(`{"symbol": "%s", "faction": "%s" }`, agentName, faction)
	registrationResult, err := SpaceTradersCommand(parameters, "register", "", "post")
	if err != nil {
		fmt.Println("Error registering. Error: ", err)
	}
	//fmt.Println(registrationResult)
	agent := &NewAgentResponse{}
	json.Unmarshal([]byte(registrationResult), &apiResult{agent})
	contract := AcceptContract(agent.Token, agent.Contract.Id)
	fmt.Printf("Was the contract accepted? %t\n", contract.ContractDetails.Accepted)
	fmt.Printf("Your deadline to complete the contract is: %s\n", contract.ContractDetails.Terms.Deadline)
	return agent.Token
}

type ContractDeliveryObject struct {
	TradeSymbol       string
	DestinationSymbol string
	UnitsRequired     int
	UnitsFulfilled    int
}

type ContractTerms struct {
	Deadline string
	Payment  map[string]int
	Deliver  []ContractDeliveryObject
}

type Contract struct {
	Id               string
	FactionSymbol    string
	Type             string
	Terms            ContractTerms
	Accepted         bool
	Fulfilled        bool
	DeadlineToAccept string
}

type ContractList struct {
	Contracts []Contract `json:"data"`
	Meta      map[string]int
}

type AcceptedContract struct {
	YourAgent       Agent
	ContractDetails Contract `json:"contract"`
}

func AcceptContract(authToken string, contractID string) AcceptedContract {
	contractURL := fmt.Sprintf("my/contracts/%s/accept", contractID)

	contractResults, err := SpaceTradersCommand("", contractURL, authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Accept Contract endpoint. Error: ", err)
	}
	acceptedContract := &AcceptedContract{}
	json.Unmarshal([]byte(contractResults), &apiResult{acceptedContract})
	return *acceptedContract
}

func ListContracts(authToken string, limit int, page int) ContractList {
	parameters := fmt.Sprintf(`{"limit": %d, "page": %d}`, limit, page)
	contractResult, err := SpaceTradersCommand(parameters, "my/contracts", authToken, "get")
	if err != nil {
		fmt.Println("Error accessing List Contracts endpoint. Error: ", err)
	}
	var contractList ContractList
	json.Unmarshal([]byte(contractResult), &contractList)
	return contractList
}

func GetContract(authToken string, contractID string) Contract {
	parameters := fmt.Sprintf(`{"contractID": "%s"}`, contractID)
	contractURL := fmt.Sprintf("my/contracts/%s", contractID)
	contractResults, err := SpaceTradersCommand(parameters, contractURL, authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Contract endpoint. Error: ", err)
	}
	contract := &Contract{}
	json.Unmarshal([]byte(contractResults), &apiResult{contract})
	return *contract
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

func GetFaction(authToken string, factionSymbol string) Faction {
	parameters := fmt.Sprintf(`{"factionSymbol": "%s"}`, factionSymbol)
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

type DestinationDepartureWaypoints struct {
	Symbol       string
	Type         string
	SystemSymbol string
	X            int
	Y            int
}

type ShipRoute struct {
	Destination DestinationDepartureWaypoints
	Departure   DestinationDepartureWaypoints
}

type ShipRegistration struct {
	Name          string
	FactionSymbol string
	Role          string
}

type ShipNav struct {
	SystemSymbol   string
	WaypointSymbol string
	Route          ShipRoute
	Status         string
	FlightMode     string
}

type ShipCrew struct {
	Current  int
	Required int
	Capacity int
	Rotation int
	Morale   int
	Wages    int
}

type ShipRequirements struct {
	Power int
	Crew  int
	Slots int
}

type ShipFrame struct {
	Symbol         string
	Name           string
	Description    string
	Condition      string
	ModuleSlots    int
	MountingPoints int
	FuelCapacity   int
	Requirements   ShipRequirements
}

type ShipReactor struct {
	Symbol       string
	Name         string
	Description  string
	Condition    int
	PowerOutput  int
	Requirements ShipRequirements
}

type ShipEngine struct {
	Symbol       string
	Name         string
	Description  string
	Condition    int
	Speed        int
	Requirements ShipRequirements
}

type ShipModule struct {
	Symbol       string
	Capacity     int
	Range        int
	Name         string
	Description  string
	Requirements ShipRequirements
}

type ShipMount struct {
	Symbol       string
	Name         string
	Description  string
	Strength     int
	Deposits     []string
	Requirements ShipRequirements
}

type ShipInventory struct {
	Symbol      string
	Name        string
	Description string
	Units       int
}

type ShipCargo struct {
	Capacity  int
	Units     int
	Inventory []ShipInventory
}

type ShipFuelConsumed struct {
	Amount    int
	Timestamp string
}

type ShipFuel struct {
	Current  int
	Capacity int
	Consumed ShipFuelConsumed
}

type Ship struct {
	Symbol       string
	Registration ShipRegistration
	Nav          ShipNav
	Crew         ShipCrew
	Frame        ShipFrame
	Reactor      ShipReactor
	Engine       ShipEngine
	Modules      []ShipModule
	Mounts       []ShipMount
	Cargo        ShipCargo
	Fuel         ShipFuel
}

type NewAgentResponse struct {
	Agent    Agent
	Contract Contract
	Faction  Faction
	Ship     Ship
	Token    string
}

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
			fmt.Printf("Response failed with status code: %d and body: %s\n", result.StatusCode, resultBody)
		}
		if err != nil {
			return "Error reading response", err
		}

		return string(resultBody), err
	}
	return "Something went wrong", fmt.Errorf("error")
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
