# Global Entry Appointment Notifier

## Build

`go build ge.go`

## Help

```
global_entry % ./ge -h
Usage of ./ge:
  -begin string
        The start of the desired time range to search for appointments in time.DateOnly format. (default "2023-12-01")
  -c string
        'The command to run. locations' will print out known locations. 'notify' will poll for open appointments.
  -end string
        The end of the desired time range to search for appointments in time.DateOnly format. (default "2024-03-01")
  -ids locations
        Comma separated list of location IDs to search for appointment openings. IDs can be found by the locations command. (default "5180, 5002, 16547")
  -sleep_duration uint
        The time to sleep between polls in seconds. (default 30)
  -to string
        Comma separated list of destination emails.
```
 ## List Locations

```
./ge -c locations
```

 ## Notify

```
./ge -c notify [-to <comma separated list of emails>]

```

## Configure

### OnSuccess Callback

The actions to take with the matching candidates are configurable. This can be
changed by re-defining the `OnSuccess` field of `PollOptions` in `ge.go`.

## Default OnSuccess Callback

The default callback will simply print the matching candidates to stdout.

### Send Mail Flag

By specifying the `-to` flag, an email will be sent using gmail's SMTP
server. These two environment variables must be defined:

* GE_USER: Email username.
* GE_APP_PW: Gmail app password.
