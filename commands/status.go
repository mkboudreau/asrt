package commands

import (
	"fmt"
	"log"
	"os"
	"time"

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

	durationString := ctx.String("retry-until")
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		cli.ShowCommandHelp(ctx, "status")
		fmt.Println("Could not get retry-until. Reason:", err)
		log.Fatalln("Exiting....")
	}

	exitStatus := retryUntil(duration, func() int {
		targetChannel := make(chan *config.Target, c.Workers)
		resultChannel := make(chan *output.Result)

		go processTargets(targetChannel, resultChannel)

		for _, target := range c.Targets {
			targetChannel <- target
		}
		close(targetChannel)

		formatter := c.ResultFormatter()
		writer := c.Writer()

		if c.AggregateOutput {
			return processAggregatedResult(resultChannel, formatter, writer, c.FailuresOnly)
		} else {
			return processEachResult(resultChannel, formatter, writer, c.FailuresOnly)
		}
	})

	os.Exit(exitStatus)
}

func retryUntil(duration time.Duration, fn func() int) int {
	retry := 2 * time.Second

	timer := time.NewTimer(duration)
	ticker := time.NewTicker(retry)

	var lastResult = 1

loop:
	for {
		lastResult = fn()
		if lastResult == 0 {
			break loop
		}
		select {
		case <-timer.C:
			break loop
		case <-ticker.C:
		}
	}

	ticker.Stop()
	timer.Stop()

	return lastResult
}
