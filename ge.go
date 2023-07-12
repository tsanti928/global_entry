package main

import (
	"flag"
	"fmt"
	ge "global_entry/global_entry"
	"global_entry/gmail"
	"log"
	"strings"
	"unicode"
)

// Update ranges here as desired.
const (
	rangeBegin    = "2023-12-01"
	rangeEnd      = "2024-03-01"
	sleepDuration = 30
)

// Update ID data structures as desired.
// Running with command=locations will provide the ids.
var (
	wantIDs = []int{5180, 5002, 16547}
)

func help() {
	fmt.Println(`Usage: ge -command <command> [-send_mail] [-dests <comma separated list of emails>]
Available commands:
locations: Display the list of known locations.
notify: Poll for appointment openings at configured location and time range.
`)
}

func main() {
	cmd := flag.String("command", "", "'locations' will print out known locations. 'notify' will poll for open appointments.")
	sendMail := flag.Bool("send_mail", false, "Set true to send a mail on matching appointment.")
	dests := flag.String("to", "", "Comma separated list of destination emails.")
	flag.Parse()

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

		trimmedDests := strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, *dests)
		to := strings.Split(trimmedDests, ",")

		err = ge.PollLocationRangesURL(ge.PollOptions{
			SleepDuration: sleepDuration,
			IDMap:         idMap,
			IDs:           wantIDs,
			RangeBegin:    rangeBegin,
			RangeEnd:      rangeEnd,
			OnSuccess: func(candidates []string) {
				subject := "Global Entry Availability"
				body := []string{"Found candidate appointments. First available slot at each location is provided."}
				body = append(body, candidates...)
				if *sendMail {
					if err := gmail.SendEmail(subject, body, to); err != nil {
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
		help()
		log.Fatalf("unknown command: %q", *cmd)
	}
}
