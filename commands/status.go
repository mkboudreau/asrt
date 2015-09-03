package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/output"
)

func cmdStatus(ctx *cli.Context) {
	c, err := config.GetConfiguration(ctx)
	if err != nil {
		cli.ShowCommandHelp(ctx, "status")
		fmt.Println("Could not get configuration. Reason:", err)
		log.Fatalln("Exiting....")
	}

	targetChannel := make(chan *config.Target, c.Workers)
	resultChannel := make(chan *output.Result)

	go processTargets(targetChannel, resultChannel)

	for _, target := range c.Targets {
		targetChannel <- target
	}
	close(targetChannel)

	formatter := c.ResultFormatter()
	writer := c.Writer()

	var exitStatus int
	if c.AggregateOutput {
		exitStatus = processAggregatedResult(resultChannel, formatter, writer, c.FailuresOnly)
	} else {
		exitStatus = processEachResult(resultChannel, formatter, writer, c.FailuresOnly)
	}

	os.Exit(exitStatus)
}
