package Event

// Event represents a basic event structure
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

