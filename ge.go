package main

import (
	"flag"
	"fmt"
	ge "global_entry/global_entry"
	"global_entry/gmail"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type applicationConfig struct {
	rangeBegin    string
	rangeEnd      string
	sleepDuration uint
	ids           []int
}

func idsStringToIntSlice(ids string) ([]int, error) {
	stringSlice := strings.Split(ids, ",")
	intSlice := make([]int, len(stringSlice))

	for i, s := range stringSlice {
		id, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return nil, fmt.Errorf("could not decode input ids: %v", err)
		}
		intSlice[i] = id
	}

	return intSlice, nil
}

func validateConfig(config *applicationConfig) error {
	if _, err := time.Parse(time.DateOnly, config.rangeBegin); err != nil {
		return fmt.Errorf("cannot convert rangeBegin %q to time.DateOnly", config.rangeBegin)
	}
	if _, err := time.Parse(time.DateOnly, config.rangeEnd); err != nil {
		return fmt.Errorf("cannot convert rangeEnd %q to time.DateOnly", config.rangeEnd)
	}

	return nil
}

func main() {
	cmd := flag.String("c", "", "'The command to run. locations' will print out known locations. 'notify' will poll for open appointments.")
	rangeBegin := flag.String("begin", "2023-12-01", "The start of the desired time range to search for appointments in time.DateOnly format.")
	rangeEnd := flag.String("end", "2024-03-01", "The end of the desired time range to search for appointments in time.DateOnly format.")
	sleepDuration := flag.Uint("sleep_duration", 30, "The time to sleep between polls in seconds.")
	ids := flag.String("ids", "5180, 5002, 16547", "Comma separated list of location IDs to search for appointment openings. IDs can be found by the `locations` command.")
	to := flag.String("to", "", "Comma separated list of destination emails.")
	flag.Parse()

	idsSlice, err := idsStringToIntSlice(*ids)
	if err != nil {
		log.Fatalln(err)
	}
	config := &applicationConfig{
		rangeBegin:    *rangeBegin,
		rangeEnd:      *rangeEnd,
		sleepDuration: *sleepDuration,
		ids:           idsSlice,
	}
	if err := validateConfig(config); err != nil {
		log.Fatalln(err)
	}

	switch *cmd {
	case "locations":
		locations, err := ge.Locations()
		if err != nil {
			log.Fatalln(err)
		}
		for _, l := range locations {
			fmt.Printf("ID: %d Name: %q State: %q City: %q\n", l.ID, l.Name, l.State, l.City)
		}
	case "notify":
		locations, err := ge.Locations()
		if err != nil {
			log.Fatalln(err)
		}
		idMap := make(ge.IDMap)
		for _, l := range locations {
			idMap[l.ID] = l
		}

		trimmedTo := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, *to)

		err = ge.PollLocationRangesURL(ge.PollOptions{
			SleepDuration: time.Duration(config.sleepDuration),
			IDMap:         idMap,
			IDs:           config.ids,
			RangeBegin:    config.rangeBegin,
			RangeEnd:      config.rangeEnd,
			OnSuccess: func(candidates []string) {
				subject := "Global Entry Availability"
				body := []string{"Found candidate appointments. First available slot at each location is provided."}
				body = append(body, candidates...)
				if trimmedTo != "" {
					splitTo := strings.Split(trimmedTo, ",")
					if err := gmail.SendEmail(subject, body, splitTo); err != nil {
						log.Fatal(err)
					}
				} else {
					for _, b := range body {
						fmt.Println(b)
					}
				}
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
	default:
		flag.Usage()
		log.Fatalf("unknown command: %q", *cmd)
	}
}
