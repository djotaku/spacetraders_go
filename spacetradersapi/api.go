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

// RegisterNewAgent will register a new agent and automatically accept your first contract
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

type deliveredContract struct {
	contract Contract
	cargo    ShipCargo
}

func DeliverCargoToContract(authToken string, contractID string, shipSymbol string, tradeSymbol string, units int) (Contract, ShipCargo) {
	parameters := fmt.Sprintf(`{"shipSymbol": %s,"tradeSymbol": %s,"units": %d}`, shipSymbol, tradeSymbol, units)
	deliveryURL := fmt.Sprintf("my/contracts/%s/deliver", contractID)
	deliveryResults, err := SpaceTradersCommand(parameters, deliveryURL, authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Deliver Cargo to Contract endpoint. Error: ", err)
	}
	delivery := &deliveredContract{}
	json.Unmarshal([]byte(deliveryResults), &apiResult{delivery})
	return delivery.contract, delivery.cargo
}

type ContractFulFilled struct {
	agent    Agent
	contract Contract
}

func FulfillContract(authToken string, contractID string) (Agent, Contract) {
	parameters := "{}"
	fulfillURL := fmt.Sprintf("my/contracts/%s/fulfill", contractID)
	fulfillResults, err := SpaceTradersCommand(parameters, fulfillURL, authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Fulfill Contract endpoint. Error: ", err)
	}
	fulfilledContract := &ContractFulFilled{}
	json.Unmarshal([]byte(fulfillResults), &apiResult{fulfilledContract})
	return fulfilledContract.agent, fulfilledContract.contract
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

type ShipList struct {
	Ships []Ship `json:"data"`
	Meta  map[string]int
}

func ListShips(authToken string, limit int, page int) ShipList {
	parameters := fmt.Sprintf(`{"limit": %d, "page": %d}`, limit, page)
	shipResult, err := SpaceTradersCommand(parameters, "my/ships", authToken, "get")
	if err != nil {
		fmt.Println("Error accessing List Ships endpoint. Error: ", err)
	}
	var shipList ShipList
	json.Unmarshal([]byte(shipResult), &shipList)
	return shipList
}

type Transaction struct {
	waypointSymbol string
	shipsymbol     string
	price          int
	agentSymbol    string
	timestamp      string
}
type PurchaseShipResult struct {
	agent       Agent
	ship        Ship
	transaction Transaction
}

func PurchaseShip(authToken string, shipType string, waypointSymbol string) (Agent, Ship, Transaction) {
	parameters := fmt.Sprintf(`{"shipType": %s, "wapointSymbol": %s}`, shipType, waypointSymbol)
	shipResult, err := SpaceTradersCommand(parameters, "my/ships", authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Purchase Ships endpoint. Error: ", err)
	}
	transactionResult := &PurchaseShipResult{}
	json.Unmarshal([]byte(shipResult), &apiResult{shipResult})
	return transactionResult.agent, transactionResult.ship, transactionResult.transaction
}

func GetShip(authToken, string, shipSymbol string) Ship {
	parameters := "{}"
	url := fmt.Sprintf("my/ships/%s", shipSymbol)
	shipResult, err := SpaceTradersCommand(parameters, url, authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Ship endpoint. Error: ", err)
	}
	ship := &Ship{}
	json.Unmarshal([]byte(shipResult), &apiResult{ship})
	return *ship
}

func GetShipCargo(authToken, string, shipSymbol string) ShipCargo {
	parameters := "{}"
	url := fmt.Sprintf("my/ships/%s/cargo", shipSymbol)
	shipResult, err := SpaceTradersCommand(parameters, url, authToken, "get")
	if err != nil {
		fmt.Println("Error accessing Get Ship Cargo endpoint. Error: ", err)
	}
	shipCargo := &ShipCargo{}
	json.Unmarshal([]byte(shipResult), &apiResult{shipCargo})
	return *shipCargo
}

func OrbitShip(authToken, string, shipSymbol string) ShipNav {
	parameters := "{}"
	url := fmt.Sprintf("my/ships/%s/orbit", shipSymbol)
	shipResult, err := SpaceTradersCommand(parameters, url, authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Orbit Ship endpoint. Error: ", err)
	}
	shipNav := &ShipNav{}
	json.Unmarshal([]byte(shipResult), &apiResult{shipNav})
	return *shipNav
}

type Cooldown struct {
	ShipSymbol       string
	totalSeconds     int
	remainingSeconds int
	expiration       string
}

// Ship Entpoints Needed, but not in the Quickstart
// Ship Refine
// Create Chart
// Get Ship Cooldown

func DockShip(authToken, string, shipSymbol string) ShipNav {
	parameters := "{}"
	url := fmt.Sprintf("my/ships/%s/dock", shipSymbol)
	shipResult, err := SpaceTradersCommand(parameters, url, authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Dock Ship endpoint. Error: ", err)
	}
	shipNav := &ShipNav{}
	json.Unmarshal([]byte(shipResult), &apiResult{shipNav})
	return *shipNav
}

// Ship Entpoints Needed, but not in the Quickstart
// Create Survey

type Yield struct {
	symbol string
	units  int
}

type Extraction struct {
	shipSymbol string
	yield      Yield
}

type ExtractionResult struct {
	Cooldown   Cooldown
	Extraction Extraction
	Cargo      ShipCargo
}

func ExtractResources(authToken, string, shipSymbol string) (Cooldown, Extraction, ShipCargo) {
	parameters := "{}"
	url := fmt.Sprintf("my/ships/%s/extract", shipSymbol)
	shipResult, err := SpaceTradersCommand(parameters, url, authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Extract Resources endpoint. Error: ", err)
	}
	extracted := &ExtractionResult{}
	json.Unmarshal([]byte(shipResult), &apiResult{extracted})
	return extracted.Cooldown, extracted.Extraction, extracted.Cargo
}

// Ship Entpoints Needed, but not in the Quickstart
// Extract Resources with Survey
// Jettison Cargo
// Jump Ship

type NavigateShipResult struct{
	fuel ShipFuel
	nav ShipNav
}

func NavigateShip(authToken, string, shipSymbol string, waypointSymbol string) (ShipFuel, ShipNav) {
	parameters := fmt.Sprintf(`{"wapointSymbol": "%s"}`, waypointSymbol)
	url := fmt.Sprintf("my/ships/%s/navigate", shipSymbol)
	shipResult, err := SpaceTradersCommand(parameters, url, authToken, "post")
	if err != nil {
		fmt.Println("Error accessing Navigate Ship endpoint. Error: ", err)
	}
	nav := &NavigateShipResult{}
	json.Unmarshal([]byte(shipResult), &apiResult{nav})
	return nav.fuel, nav.nav
}

// Ship Entpoints Needed, but not in the Quickstart
// Patch Ship Nav
// Get Ship Nav
// Warp Ship



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
		if err != nil {
			return "Error doing request", err
		}

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
