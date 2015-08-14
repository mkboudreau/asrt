package main

import (
	"os"
	"path"

	"github.com/codegangsta/cli"

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
	app.CommandNotFound = commands.CommandNotFound
	app.Usage = "API Status Reporting Tool"
	app.Version = version.Version + " (" + version.GitCommit + ")"

	app.Run(os.Args)
}
