package domain

// RendezvousRequest : A struct that represents a single rendezvous API request
type RendezvousRequest struct {
	ID     string                 `json:"id"`
	Events map[string]interface{} `json:"events"`
	Data   string                 `json:"data"`
}
