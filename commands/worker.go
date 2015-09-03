package commands

import (
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/execution"
	"github.com/mkboudreau/asrt/output"
	"github.com/mkboudreau/asrt/writer"
)

func processTargets(incomingTargets <-chan *config.Target, resultChannel chan<- *output.Result) {
	var wg sync.WaitGroup

	for t := range incomingTargets {
		wg.Add(1)
		go func(target *config.Target) {
			execResult := execution.ExecuteWithTimoutAndHeaders(string(target.Method), target.URL, target.Timeout, target.Headers, target.ExpectedStatus)

			result := output.NewResult(execResult.Success(), execResult.Error, strconv.Itoa(execResult.Expected), strconv.Itoa(execResult.Actual), execResult.URL)
			result.Timestamp = output.NewTimeStringForJSON(time.Now())
			resultChannel <- result
			wg.Done()
		}(t)
	}

	wg.Wait()
	close(resultChannel)
}

func processEachResult(resultChannel <-chan *output.Result, formatter output.ResultFormatter, w io.Writer, onlyFailures bool) int {
	exitStatus := 0
	counter := 0
	for r := range resultChannel {
		if r.Success && onlyFailures {
			continue
		}
		if counter == 0 {
			writer.WriteToWriter(w, formatter.Header())
		}
		if counter != 0 {
			writer.WriteToWriter(w, formatter.RecordSeparator())
		}
		reader := formatter.Reader(r)
		if !r.Success {
			exitStatus = 1
		}
		writer.WriteToWriter(w, reader)
		counter++
	}
	if counter > 0 {
		writer.WriteToWriter(w, formatter.Footer())
	}
	writer.DoneWithWriter(w)

	return exitStatus
}

func processAggregatedResult(resultChannel <-chan *output.Result, formatter output.ResultFormatter, w io.Writer, onlyFailures bool) int {
	exitStatus := 0
	results := make([]*output.Result, 0)
	for r := range resultChannel {
		results = append(results, r)
		if !r.Success {
			exitStatus = 1
		}
	}

	if !onlyFailures || (exitStatus == 1 && onlyFailures) {
		reader := formatter.AggregateReader(results)
		writer.WriteToWriter(w, formatter.Header())
		writer.WriteToWriter(w, reader)
		writer.WriteToWriter(w, formatter.Footer())
	}

	writer.DoneWithWriter(w)

	return exitStatus
}
