package main

import (
	"fmt"
	"github.com/codegangsta/cli"
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
	fmt.Printf(
		"%s: '%s' is not an %s command. See '%s --help'.\n",
		c.App.Name,
		command,
		c.App.Name,
		c.App.Name,
	)
}
