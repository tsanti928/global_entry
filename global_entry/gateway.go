package global_entry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// timeInRange returns true if `input` is inbetween `begin` and `end` times.
// input should be in time.DateTime format.
// begin and end should be in time.DateOnly format.
// Times aren't converted to a local time because for comparison purposes the
// outcome will be the same, since they all would shift the same.
func timeInRange(input, begin, end string) (bool, error) {
	inputTime, err := time.Parse(time.DateTime, input)
	if err != nil {
		return false, err
	}
	beginTime, err := time.Parse(time.DateOnly, begin)
	if err != nil {
		return false, err
	}
	endTime, err := time.Parse(time.DateOnly, end)
	if err != nil {
		return false, err
	}
	inputUnix := inputTime.Unix()
	beginUnix := beginTime.Unix()
	endUnix := endTime.Unix()

	return inputUnix >= beginUnix && inputUnix <= endUnix, nil
}

// APITimeToDateTime converts global entry API time data to time.DateTime.
// API returns time in the format: 2024-05-15T06:45.
// This doesn't match a known time format in Go, but it can converted to
// time.DateTime easily.
func APITimeToDateTime(input string) string {
	return strings.Replace(input, "T", " ", 1) + ":00"
}

// PollSlotAvailabilityURL polls the slot availability URL for matching
// available appointments for input location IDs in provided time range.
// Only the first matching time for each location is used as a candidate.
func PollSlotAvailabilityURL(options PollOptions) error {
	var data SlotAvailabilityData
	for {
		fmt.Println("Starting the queries...")

		var candidates []string
		for _, id := range options.IDs {
			name := options.IDMap[id].Name
			loc, err := time.LoadLocation(options.IDMap[id].TimeZone)
			if err != nil {
				return err
			}
			url := slotAvailabilityURL(id)
			//fmt.Printf("Querying url: %q\n", url)

			decode := func() error {
				resp, err := http.Get(url)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					return err
				}
				return nil
			}
			if err := decode(); err != nil {
				return err
			}

			if len(data.SlotData) > 0 {
				dateTime := APITimeToDateTime(data.SlotData[0].StartTimestamp)
				inRange, err := timeInRange(dateTime, options.RangeBegin, options.RangeEnd)
				if err != nil {
					return err
				}

				if inRange {
					time, err := time.ParseInLocation(time.DateTime, dateTime, loc)
					if err != nil {
						return err
					}

					candidates = append(candidates, fmt.Sprintf("City: %s Name: %s ID: %d Time: %s", options.IDMap[id].City, name, id, time.String()))
				}
			}
		}
		if len(candidates) > 0 {
			options.OnSuccess(candidates)
			return nil
		}
		fmt.Printf("No matches found. Trying again in %d seconds.\n", options.SleepDuration)
		time.Sleep(options.SleepDuration * time.Second)
	}
}

// PollLocationRangesURL polls the location URL for matching
// available appointments for input location IDs in provided time range.
// Only the first matching time for each location is used as a candidate.
func PollLocationRangesURL(options PollOptions) error {
	for {
		fmt.Println("Starting the queries...")

		var candidates []string
		for _, id := range options.IDs {
			var data []LocationRangesData
			name := options.IDMap[id].Name
			loc, err := time.LoadLocation(options.IDMap[id].TimeZone)
			if err != nil {
				return err
			}
			url := locationRangesURL(id, options.RangeBegin, options.RangeEnd)
			//fmt.Printf("Querying url: %q\n", url)

			decode := func() error {
				resp, err := http.Get(url)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					return err
				}
				return nil
			}
			if err := decode(); err != nil {
				return err
			}

			if len(data) > 0 {
				for _, d := range data {
					if d.Active == 0 {
						continue
					}

					dateTime := APITimeToDateTime(d.Timestamp)
					time, err := time.ParseInLocation(time.DateTime, dateTime, loc)
					if err != nil {
						return err
					}

					candidates = append(candidates, fmt.Sprintf("City: %s Name: %s ID: %d Time: %s", options.IDMap[id].City, name, id, time.String()))
					break
				}
			}
		}
		if len(candidates) > 0 {
			options.OnSuccess(candidates)
			return nil
		}
		fmt.Printf("No matches found. Trying again in %d seconds.\n", options.SleepDuration)
		time.Sleep(options.SleepDuration * time.Second)
	}
}

// Locations querys locations EP and returns the data.
func Locations() ([]*LocationInfo, error) {
	var data []*LocationInfo

	url := locationsURL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
