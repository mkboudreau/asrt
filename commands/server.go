package commands

import (
	"github.com/codegangsta/cli"
	//"log"
)

func cmdServer(c *cli.Context) {
	/*
		config, err := getConfiguration(c)
		if err != nil {
			cli.ShowCommandHelp(c, "status")
			log.Fatal(err)
		}

		outputFunction := getOutputFunction(config)

		targetChannel := make(chan *target, config.Workers)
		resultChannel := make(chan *output.Result)

		go processTargets(targetChannel, resultChannel)

		for _, target := range config.Targets {
			targetChannel <- target
		}
		close(targetChannel)

		if config.AggregateOutput {
			processAggregatedResult(resultChannel, config.Quiet, outputFunction)
		} else {
			processEachResult(resultChannel, config.Quiet, outputFunction)
		}
	*/
}
