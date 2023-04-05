package test

type FrigatePayload struct {
	Type   string       `json:"type"`
	Before FrigateEvent `json:"before"`
	After  FrigateEvent `json:"after"`
}

type FrigateEvent struct {
	EnteredZones []string `json:"entered_zones"`
	Camera       string   `json:"camera"`
	Label        string   `json:"label"`
}
