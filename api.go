// The API for interacting with Space Traders

package api

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
