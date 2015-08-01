package commands

import (
	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/execution"
	"github.com/mkboudreau/asrt/output"
	"log"
	"strconv"
	"sync"
)

func cmdStatus(c *cli.Context) {
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
}

func processEachResult(resultChannel <-chan *output.Result, quiet bool, outputFunction func(result *output.Result)) {
	for r := range resultChannel {
		if quiet {
			result := &output.Result{
				Success:  r.Success,
				Error:    r.Error,
				Expected: "",
				Url:      "",
			}
			outputFunction(result)
		} else {
			outputFunction(r)
		}
	}
}

func processAggregatedResult(resultChannel <-chan *output.Result, quiet bool, outputFunction func(result *output.Result)) {
	result := &output.Result{
		Success:  true,
		Error:    nil,
		Expected: "",
		Url:      "",
	}
	count := 0
	for r := range resultChannel {
		if !r.Success {
			result.Success = false
		}
		if r.Error != nil {
			result.Success = false
			result.Error = r.Error
		}
		count = count + 1
	}

	if !quiet {
		result.Expected = strconv.Itoa(count)
	}

	outputFunction(result)
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

func getOutputFunction(config *configuration) func(result *output.Result) {
	var fn func(result *output.Result)

	switch {
	case config.Output == formatJSON:
		if config.Pretty {
			fn = output.WriteToJsonPretty
		} else {
			fn = output.WriteToJson
		}
	case config.Quiet:
		if config.Pretty {
			fn = output.WriteToQuietPretty
		} else {
			fn = output.WriteToQuiet
		}
	case config.Output == formatCSV:
		if config.Pretty {
			fn = output.WriteToCsvPretty
		} else {
			fn = output.WriteToCsv
		}
	case config.Output == formatTAB:
		if config.Pretty {
			fn = output.WriteToTabPretty
		} else {
			fn = output.WriteToTab
		}
	}

	return fn
}
