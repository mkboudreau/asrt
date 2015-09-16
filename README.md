# API Simple Reporting Tool (ASRT)

## Summary
API Simple Reporting Tool (ASRT) tool has a simple goal: Report up/down status of API endpoints.

## Description
ASRT (pronounced assert) was created to simply report on api endpoint statuses. It is not intended to be a stress tester like other tools, such as vegeta. It is just meant to be a simple reporter that a dashboard like tool could use. It can be easily integrated with any build processes. Not only is the output easily parseable (especitally with cut), but it also sets the exit status to 1 if it is not 100% successful.

## Installation

`go get github.com/mkboudreau/asrt`

## Example Usage
This tool is simple. Just dive in and start using it.

- Simple tab delimited status

`asrt status www.yahoo.com`

- Add pretty printing

`asrt status -p www.yahoo.com`

- Refresh every 5 seconds

`asrt dashboard -p -r 5s www.yahoo.com`

- Change format to csv

`asrt status -p -fmt csv www.yahoo.com`

- Change format to json

`asrt status -p -fmt json www.yahoo.com`

- Add a site to the list and pretty print json

`asrt status -p -fmt json www.yahoo.com www.google.com`

- Read from a file and pretty print json

`asrt status -p -fmt json -f sites.list`

- Read from a file and refresh every 1 minute with pretty print json

`asrt dashboard -p -fmt json -r 1m -f sites.list`

- Read from a file and expose online with data refresh every 1 minute

`asrt server -fmt json -r 1m -f sites.list`

Example File
```
www.yahoo.com
www.yahoo.com|200
www.yahoo.com|GET|200
```

## Configuration

### Main Commands
- status
- dashboard
- server

### Options
- `-d` or `--debug`: turns on the logger

#### Input-related Options
- `-f` or `--file`: input file (see formats below). at least a file or urls on the command line must be specified.

#### Processing-related Options
- `-w` or `--workers`: this is the number of workers or goroutines used to connect to the client sites.
- `-t` or `--timeout`: timeout for connections in time.Duration format. defaults to no timeout.
- `-a` or `--aggregate`: aggregates all sites into a single true/false response. includes the total sites unless -q is specified
- `--failures-only`: only submits data to the writer upon failure. Useful when using something like slack since you may only want to perform an http.POST upon a failure.

#### Output-related Options
- `-fmt` or `--format`: the format to be used for output. valid values are CSV,TAB,JSON. default is TAB.
- `-q` or `--quiet`: minimizes the response and contains no header nor footer.
- `-qq` or `--quieter`: turns off standard output, superceding -q. useful for scripts.
- `-p` or `--pretty`: makes response have some formatting using escape codes. Mutually exclusive with markdown.
- `-md` or `--markdown`: makes response have some formatting using markdown. Mutually exclusive with pretty.
- `--slack-url`: setting this parameter enables slack integration using incoming webhook url specified.
- `--slack-user`: overrides the user this tool will post as to slack. only works if slack-url is specified.
- `--slack-channel`: overrides the channel this tool will post to on slack. only works if slack-url is specified.
- `--slack-icon`: overrides the icon this tool will use when posting to slack. only works if slack-url is specified.

#### Options for Dashboard Command
- `-r` or `--rate`: refresh rate for dashboard only in time.Duration format. defaults to 30s.

#### Options for Server Command
- `-r` or `--rate`: refresh rate for dashboard only in time.Duration format. defaults to 30s.
- `--port`: set port to listen on (only works with server command). default is 7070

### Input Format for Target Endpoints

`URL|METHOD|STATUS_CODE|LABEL`

- The input is order dependent.
- URL is the only field that is required
- Method, status code and label are all optional
- Method must be one of GET, POST, PUT, PATCH, HEAD, OPTIONS
- Status code must be an integer
- Label can only have spaces if it is within quotes
- *if the url has the | character, it should also be placed within quotes*
- Examples
    + `www.yahoo.com`
    + `www.yahoo.com|200`
    + `www.yahoo.com|GET|200`
    + `www.microsoft.com|POST`
    + `www.microsoft.com|POST|201`
    + `www.yahoo.com/not_found|GET|404`
    + `www.yahoo.com/not_found|GET|404|{H}"Authorization: Bearer 123"`
    + `data.asrt.io|GET|200|"Main ASRT API Endpoint"`

### Differences between passing input via command line parameter and by input file

Command line parameter should be in a single line and each target be separated by spaces.
Example: `asrt www.yahoo.com www.microsoft.com|POST|201 data.asrt.io|GET|200|"Main ASRT API Endpoint"`

Input file should have one line per target.
Example:
`
www.yahoo.com
www.microsoft.com|POST|201
data.asrt.io|GET|200|"Main ASRT API Endpoint"
`

## TODO

The following items are still outstanding:
- Add support for headers, both as an option and from within a file
- Add more tests!
- Add ability to pipe input
