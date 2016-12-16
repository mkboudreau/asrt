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

- Refresh every 5 seconds

`asrt dashboard -r 5s www.yahoo.com`

- Change format to csv

`asrt status -fmt csv www.yahoo.com`

- Change format to csv with markdown formatting

`asrt status -fmt csv-md www.yahoo.com`

- Change format to csv with no color 

`asrt status -fmt csv-no-color www.yahoo.com`

- Change format to json

`asrt status -fmt json www.yahoo.com`

- Add a site to the list and change format to compact json

`asrt status -fmt json-compact www.yahoo.com www.google.com`

- Read from a file and pretty print json

`asrt status -fmt json -f sites.list`

- Read from a file and refresh every 1 minute with json

`asrt dashboard -fmt json -r 1m -f sites.list`

- Read from a file and expose online with data refresh every 1 minute

`asrt server -fmt json -r 1m -f sites.list`

- Custom formatting using Go text template syntax.

`asrt status -f sites.list -fmt template="{{ .Label }} is {{ .Success }}" `

- Custom formatting using a file with Go text template syntax.

`asrt status -f sites.list -fmt template-file=my-output.txt`


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
- `-fmt` or `--format`: the format to be used for output. valid values are: 
	CSV, CSV-MD, CSV-NO-COLOR, TAB, TAB-MD, TAB-NO-COLOR, JSON, JSON-COMPACT, TEMPLATE="{{...}}}, TEMPLATE-FILE=<filename>. default is TAB.
- `--no-headers`: minimizes the response and contains no header nor footer.
- `-q` or `--quiet`: turns off standard output, useful for scripts.
- `--slack-url`: setting this parameter enables slack integration using incoming webhook url specified.
- `--slack-user`: overrides the user this tool will post as to slack. only works if slack-url is specified.
- `--slack-channel`: overrides the channel this tool will post to on slack. only works if slack-url is specified.
- `--slack-icon`: overrides the icon this tool will use when posting to slack. only works if slack-url is specified.
- `--http-url`: setting this parameter enables generic http integration for sending output data.
- `--http-method`: sets the http method to use. default is POST. only works if http-url is specified.
- `--http-auth`: sets the http Authorization header. only works if http-url is specified.

#### Options for Status Command
- `-ru` or `--retry-until`: retry until a time.Duration. Status command will return result once it gets 100% success from all its targets or until this time.Duration elapses. If time.Duration is reached, it will return whatever the last response was. 

#### Options for Dashboard Command
- `-r` or `--rate`: refresh rate for dashboard only in time.Duration format. defaults to 30s.

#### Options for Server Command
- `-r` or `--rate`: refresh rate for dashboard only in time.Duration format. defaults to 30s.
- `--port`: set port to listen on (only works with server command). default is 7070

### Input Format for Target Endpoints

`URL|METHOD|STATUS_CODE|LABEL|HEADER...`

- The input is order dependent.
- URL is the only field that is required
- Method, status code and label are all optional
- Method must be one of GET, POST, PUT, PATCH, HEAD, OPTIONS
- Status code must be an integer
- Label can only have spaces if it is within quotes
- Headers take the same format and rules as a label (mostly), so to differentiate them, a header must contain the {H} prefix. It is the last element and there can be as many as you need, each separated by a '|'
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

### Go Text Templating

Using the `--fmt template=...` will cause the template text on the right side of "template=" to get parsed according to the Go standard library's text templating.

	type Result struct {
		Success   bool
		Error     error
		Expected  string
		Actual    string
		Url       string
		Label     string
		Timestamp string
		Extra     map[string]interface{} 
	}


See https://golang.org/pkg/text/template/

## TODO

The following items are still outstanding:
- Add more tests!
- Add ability to pipe input
