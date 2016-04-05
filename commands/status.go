package commands

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/execution"
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

	executor := execution.NewExecutor(c.AggregateOutput, c.FailuresOnly, true, c.ResultFormatter(), c.Writer(), c.Workers)
	exitStatus := retryUntil(duration, func() int {
		return executor.Execute(c.Targets)
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
