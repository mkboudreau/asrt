package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/output"
)

func cmdStatus(c *cli.Context) {
	config, err := getConfiguration(c)
	if err != nil {
		cli.ShowCommandHelp(c, "status")
		fmt.Println("Could not get configuration. Reason:", err)
		log.Fatalln("Exiting....")
	}

	targetChannel := make(chan *target, config.Workers)
	resultChannel := make(chan *output.Result)

	go processTargets(targetChannel, resultChannel)

	for _, target := range config.Targets {
		targetChannel <- target
	}
	close(targetChannel)

	formatter := config.ResultFormatter()
	writer := config.Writer()

	var exitStatus int
	if config.AggregateOutput {
		exitStatus = processAggregatedResult(resultChannel, formatter, writer, config.FailuresOnly)
	} else {
		exitStatus = processEachResult(resultChannel, formatter, writer, config.FailuresOnly)
	}

	os.Exit(exitStatus)
}
