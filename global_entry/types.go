package global_entry

import "time"

// IDMap maps a location id to name.
type IDMap map[int]*LocationInfo

// SlotData is the nested structure returned for slot-availability requests.
type SlotData struct {
	LocationID     int    `json:"locationId"`
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
	Active         bool   `json:"active"`
	Duration       int    `json:"duration"`
	RemoteInd      bool   `json:"remoteInd"`
}

// SlotAvailabilityData is the data type sent over the wire for the
// slot-availability endpoint.
type SlotAvailabilityData struct {
	SlotData          []SlotData `json:"availableSlots"`
	LastPublishedData string     `json:"lastPublishedDate"`
}

// LocationRangesData is the data type sent over the wire for the
// location ranges endpoint.
type LocationRangesData struct {
	Active    int    `json:"active"`
	Total     int    `json:"total"`
	Pending   int    `json:"pending"`
	Conflicts int    `json:"conflicts"`
	Duration  int    `json:"duration"`
	Timestamp string `json:"Timestamp"`
	Remote    bool   `json:"remote"`
}

// LocationInfo is the data type sent over the wire for locations endpoint.
type LocationInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	TimeZone   string `json:"tzData"`
}

// PollOptions are options to configure polling functions.
type PollOptions struct {
	SleepDuration time.Duration
	IDMap         IDMap
	IDs           []int

	// Ranges should be of the form: "2023-12-01", also known as time.DateOnly.
	RangeBegin string
	RangeEnd   string

	OnSuccess func(candidates []string)
}
