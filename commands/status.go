package commands

import (
	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/output"
	"log"
	"os"
)

func cmdStatus(c *cli.Context) {
	config, err := getConfiguration(c)
	if err != nil {
		cli.ShowCommandHelp(c, "status")
		log.Fatal(err)
	}

	targetChannel := make(chan *target, config.Workers)
	resultChannel := make(chan *output.Result)

	go processTargets(targetChannel, resultChannel)

	for _, target := range config.Targets {
		targetChannel <- target
	}
	close(targetChannel)

	formatter := getResultFormatter(config)
	var exitStatus int
	if config.AggregateOutput {
		exitStatus = processAggregatedResult(resultChannel, formatter)
	} else {
		exitStatus = processEachResult(resultChannel, formatter)
	}

	os.Exit(exitStatus)
}

func processEachResult(resultChannel <-chan *output.Result, formatter output.ResultFormatter) int {
	exitStatus := 0
	for r := range resultChannel {
		reader := formatter.Reader(r)
		if !r.Success {
			exitStatus = 1
		}
		output.WriteToConsole(reader)
	}

	return exitStatus
}

func processAggregatedResult(resultChannel <-chan *output.Result, formatter output.ResultFormatter) int {
	exitStatus := 0
	results := make([]*output.Result, 0)
	for r := range resultChannel {
		results = append(results, r)
		if !r.Success {
			exitStatus = 1
		}
	}

	reader := formatter.AggregateReader(results)
	output.WriteToConsole(reader)

	return exitStatus
}
