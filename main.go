package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"path"

	"github.com/mkboudreau/asrt/commands"
	_ "github.com/mkboudreau/asrt/log"
	"github.com/mkboudreau/asrt/version"
)

func main() {

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = "Michael Boudreau"
	app.Email = "https://github.com/mkboudreau/asrt"
	app.Commands = commands.Commands
	app.CommandNotFound = cmdNotFound
	app.Usage = "API Status Reporting Tool"
	app.Version = version.Version + " (" + version.GitCommit + ")"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			EnvVar: "DEBUG",
			Name:   "debug, d",
			Usage:  "Enable debug logging mode",
		},
	}

	app.Run(os.Args)
}

func cmdNotFound(c *cli.Context, command string) {
	log.Fatalf(
		"%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name,
		command,
		c.App.Name,
		c.App.Name,
	)
}

// <cmd> status ...url
// -t | --timeout
//
// -g = GET
// -p = POST
// -u = PUT
// -d = DELETE
// -h = HEAD
// or
// -f <file> = FILE

// <cmd> dashboard ...url
// -t | --timeout
// -r | --rate
//
// -g | --get = GET
// -p | --post = POST
// -u | --put = PUT
// -d | --delete = DELETE
// -h | --head = HEAD
// -hdr | --header = 1..*
// or
// -f <file> = FILE

// FILE FORMAT
// <method> <url>
// <header 1>
// <header n>
// examples
// GET http://abc.com
// GET http://abcd.com
// Accept: application/json
// Content-Type: application/json
// POST http://abd.com
