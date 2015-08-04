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

- Change format to csv
`asrt status -p -fmt csv www.yahoo.com`

- Change format to json
`asrt status -p -fmt json www.yahoo.com`

- Add a site to the list
`asrt status -p -fmt json www.yahoo.com www.google.com`

- Read from a file
`asrt status -p -fmt json -f sites.list`

Example File
```
GET www.yahoo.com 
GET www.google.com 200
```

## Configuration

### Main Commands
- status
- dashboard (not implemented yet)
- server (not implemented yet)

### Options
- `-d` or `--debug`: turns on the logger
- `-t` or `--timeout`: timeout for connections in time.Duration format. defaults to no timeout.
- `-r` or `--rate`: refresh rate for dashboard only in time.Duration format. defaults to 30s.
- `-f` or `--file`: input file (see formate below). at least a file or urls on the command line must be specified.
- `-fmt` or `--format`: the format to be used for output. valid values are CSV,TAB,JSON. default is TAB.
- `-a` or `--aggregate`: aggregates all sites into a single true/false response. includes the total sites unless -q is specified
- `-q` or `--quiet`: makes response only show true/false.
- `-w` or `--workers`: this is the number of workers or goroutines used to connect to the client sites.
- `-m` or `--method`: the method to be used for all urls given on the command line. valid values are GET,POST,PUT,DELETE,HEAD,PATCH. default is GET. (not implemented yet)

### Input File Format
- Each line must contain a method and a url. 
- Each line may optionally include an expected status code.
- These values must be in order and separated with either a tab or a space.
- Examples
    + `GET www.yahoo.com`
    + `GET www.yahoo.com 200`
    + `GET www.yahoo.com/hello 404`
    + `POST www.microsoft.com/hello 404`
    + `POST www.microsoft.com 201`
    + `POST www.microsoft.com`

## TODO

The following items are still outstanding:
- Add support for headers, both as an option and from within a file
- Add more tests!
- Add ability to pipe input 
- Add dashboard that loops and rewrites/repaints the screen
- Add http server to it can report its findings outside the terminal

