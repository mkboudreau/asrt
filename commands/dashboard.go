package commands

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/execution"
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

	executor := execution.NewExecutor(c.AggregateOutput, c.FailuresOnly, true, c.ResultFormatter(), c.Writer(), c.Workers)
	printDashboard(c, executor)
	loopDashboard(c, executor)
}

func loopDashboard(c *config.Configuration, executor *execution.Executor) {
	done := make(chan struct{})
	fn := func() {
		close(done)
	}

	OsSignalShutdown(fn, 5)

	ticker := time.NewTicker(c.Rate)

	for {
		select {
		case <-ticker.C:
			printDashboard(c, executor)
		case <-done:
			return
		}
	}
}

func printDashboard(c *config.Configuration, executor *execution.Executor) {
	clearAndWriteHeader(c)
	executor.Execute(c.Targets)
}

func clearAndWriteHeader(c *config.Configuration) {
	writer.ClearConsole()

	var timeReader io.Reader
	if c.Pretty {
		timeReader = output.NewPrettyTimeReader(time.Now())
	} else {
		timeReader = output.NewTimeReader(time.Now())
	}
	if !c.Quiet {
		writer.WriteToConsole(timeReader)
	}
}
