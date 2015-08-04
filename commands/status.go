package commands

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/execution"
	"github.com/mkboudreau/asrt/output"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
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

func consoleWriter(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			fmt.Print(line)
		}
		if err != nil {
			break
		}
	}
}

func processEachResult(resultChannel <-chan *output.Result, formatter output.ResultFormatter) int {
	exitStatus := 0
	for r := range resultChannel {
		reader := formatter.Reader(r)
		if !r.Success {
			exitStatus = 1
		}
		consoleWriter(reader)
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
	consoleWriter(reader)

	return exitStatus
}

func processTargets(incomingTargets <-chan *target, resultChannel chan<- *output.Result) {
	var wg sync.WaitGroup

	for t := range incomingTargets {
		wg.Add(1)
		go func(target *target) {
			ok, err := execution.ExecuteWithTimoutAndHeaders(string(target.Method), target.URL, target.Timeout, target.Headers, target.ExpectedStatus)
			result := output.NewResult(ok, err, strconv.Itoa(target.ExpectedStatus), target.URL)
			resultChannel <- result
			wg.Done()
		}(t)
	}

	wg.Wait()
	close(resultChannel)
}

func getResultFormatter(config *configuration) output.ResultFormatter {
	modifiers := &output.ResultFormatModifiers{
		Pretty:    config.Pretty,
		Aggregate: config.AggregateOutput,
		Quiet:     config.Quiet,
	}

	switch {
	case config.Output == formatJSON:
		return output.NewJsonResultFormatter(modifiers)
	case config.Output == formatCSV:
		return output.NewCsvResultFormatter(modifiers)
	case config.Output == formatTAB:
		return output.NewTabResultFormatter(modifiers)
	}

	return nil
}
