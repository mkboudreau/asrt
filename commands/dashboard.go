package commands

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/output"
	"github.com/mkboudreau/asrt/writer"
)

func cmdDashboard(ctx *cli.Context) {
	c, err := config.GetConfiguration(ctx)
	if err != nil {
		cli.ShowCommandHelp(ctx, "dashboard")
		fmt.Println("Could not get configuration. Reason:", err)
		log.Fatalln("Exiting....")
	}

	printDashboard(c)
	loopDashboard(c)
}

func loopDashboard(c *config.Configuration) {
	done := make(chan struct{})
	fn := func() {
		close(done)
	}

	OsSignalShutdown(fn, 5)

	ticker := time.NewTicker(c.Rate)

	for {
		select {
		case <-ticker.C:
			printDashboard(c)
		case <-done:
			return
		}
	}
}

func printDashboard(c *config.Configuration) {
	writer.ClearConsole()

	var timeReader io.Reader
	if c.Pretty {
		timeReader = output.NewPrettyTimeReader(time.Now())
	} else {
		timeReader = output.NewTimeReader(time.Now())
	}
	writer.WriteToConsole(timeReader)

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
		processAggregatedResult(resultChannel, formatter, writer, c.FailuresOnly)
	} else {
		processEachResult(resultChannel, formatter, writer, c.FailuresOnly)
	}
}
