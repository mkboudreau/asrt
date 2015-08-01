package commands

import (
	"github.com/codegangsta/cli"
	"log"
)

func cmdDashboard(c *cli.Context) {
	_, err := getConfiguration(c)
	if err != nil {
		cli.ShowCommandHelp(c, "dashboard")
		log.Fatal(err)
	}
}
