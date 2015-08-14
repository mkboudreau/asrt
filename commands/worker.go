package commands

import (
	"io"
	"strconv"
	"sync"

	"github.com/mkboudreau/asrt/execution"
	"github.com/mkboudreau/asrt/output"
)

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

func processEachResult(resultChannel <-chan *output.Result, formatter output.ResultFormatter, writer io.Writer) int {
	exitStatus := 0
	output.WriteToWriter(writer, formatter.Header())
	for r := range resultChannel {
		reader := formatter.Reader(r)
		if !r.Success {
			exitStatus = 1
		}
		output.WriteToWriter(writer, reader)
	}
	output.WriteToWriter(writer, formatter.Footer())
	output.DoneWithWriter(writer)

	return exitStatus
}

func processAggregatedResult(resultChannel <-chan *output.Result, formatter output.ResultFormatter, writer io.Writer) int {
	exitStatus := 0
	results := make([]*output.Result, 0)
	for r := range resultChannel {
		results = append(results, r)
		if !r.Success {
			exitStatus = 1
		}
	}

	reader := formatter.AggregateReader(results)
	output.WriteToWriter(writer, formatter.Header())
	output.WriteToWriter(writer, reader)
	output.WriteToWriter(writer, formatter.Footer())
	output.DoneWithWriter(writer)

	return exitStatus
}
