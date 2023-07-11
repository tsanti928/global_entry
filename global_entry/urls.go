package global_entry

import "fmt"

const (
	// This URL will get all locations.
	locationsURL = `https://ttp.cbp.dhs.gov/schedulerapi/locations/?temporary=false&inviteOnly=false&operational=true&serviceName=Global%20Entry`
	slotsURL     = `https://ttp.cbp.dhs.gov/schedulerapi/slots/asLocations?minimum=1&limit=5&serviceName=Global%20Entry`
)

// slotAvailabilityURL returns a URL for the slot-availability endpoint.
// This finds the next available appointment.
func slotAvailabilityURL(id int) string {
	return fmt.Sprintf("https://ttp.cbp.dhs.gov/schedulerapi/slot-availability?locationId=%d", id)
}

// slotByIDURL returns a URL for the slots endpoint with a location specified.
// This finds available appointments starting from soonest.
func slotByIDURL(id int) string {
	return fmt.Sprintf("https://ttp.cbp.dhs.gov/schedulerapi/slots?orderBy=soonest&limit=1&locationId=%d&minimum=1", id)
}

// locationRangesURL returns a URL to query for open slots for a location
// within a range.
// This finds appointments in the time range, but they may not be
// active/available, so the 'active' field should be checked.
// `begin` and `end` must take the format "YYYY-MM-DD".
// Input validation is not performed.
func locationRangesURL(id int, begin, end string) string {
	// Note that %3A is the URL encoding of ASCII ':'.
	// This timeSuffix is for start of day boundary.
	// I don't know whether the API expects a UTC time, or will adjust based
	// on the local time of the location. This means there may be a ~7 hour
	// margin of error for California.
	timeSuffix := "T00%3A00%3A00"
	return fmt.Sprintf("https://ttp.cbp.dhs.gov/schedulerapi/locations/%d/slots?startTimestamp=%s%s&endTimestamp=%s%s", id, begin, timeSuffix, end, timeSuffix)
}
