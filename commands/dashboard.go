package commands

import (
	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/output"
	"io"
	"log"
	"time"
)

func cmdDashboard(c *cli.Context) {
	config, err := getConfiguration(c)
	if err != nil {
		cli.ShowCommandHelp(c, "dashboard")
		log.Fatal(err)
	}

	printDashboard(config)
	loopDashboard(config)
}

func loopDashboard(config *configuration) {
	done := make(chan struct{})
	fn := func() {
		close(done)
	}

	osSignalShutdown(fn, 5)

	ticker := time.NewTicker(config.Rate)

	for {
		select {
		case <-ticker.C:
			printDashboard(config)
		case <-done:
			return
		}
	}
}

func printDashboard(config *configuration) {
	output.ClearConsole()

	var timeReader io.Reader
	if config.Pretty {
		timeReader = output.NewPrettyTimeReader(time.Now())
	} else {
		timeReader = output.NewTimeReader(time.Now())
	}
	output.WriteToConsole(timeReader)

	targetChannel := make(chan *target, config.Workers)
	resultChannel := make(chan *output.Result)

	go processTargets(targetChannel, resultChannel)

	for _, target := range config.Targets {
		targetChannel <- target
	}
	close(targetChannel)

	formatter := getResultFormatter(config)
	if config.AggregateOutput {
		processAggregatedResult(resultChannel, formatter)
	} else {
		processEachResult(resultChannel, formatter)
	}
}
