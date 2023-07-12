# Global Entry Appointment Notifier

## Build

`go build ge.go`

## Help

```
global_entry % ./ge -h
Usage of ./ge:
  -command string
        'locations' will print out known locations. 'notify' will poll for open appointments.
  -send_mail
        Set true to send a mail on matching appointment.
  -to string
        Comma separated list of destination emails.
```
 ## List Locations

```
./ge -command locations
```

 ## Notify

```
./ge -command notify [-send_mail] [-dests <comma separated list of emails>]

```

## Configure

### Source Code Variables
1. `rangeBegin`: The start of the desired time range to search for appointments in time.DateOnly format.
1. `rangeEnd`: The end of the desired time range to search for appointments in time.DateOnly format.
1. `sleepDuration`: The time to sleep between polls in seconds.
1. `wantIDs`: The list of location IDs to search for appointment openings. IDs can be found by the `locations` command.

```
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
```

### OnSuccess Callback

The actions to take with the matching candidates are configurable. This can be
changed by re-defining the `OnSuccess` field of `PollOptions` in `ge.go`.

## Default OnSuccess Callback

The default callback will simply print the matching candidates to stdout.

### Send Mail Flag

By specifying the flag `-send_mail`, an email will be sent using gmail's SMTP
server. These two environment variables must be defined:

* GE_USER: Email username.
* GE_APP_PW: Gmail app password.

Accompanying the `-send_mail` flag should be the comma separated list of destination
email address with the `-to` flag.
