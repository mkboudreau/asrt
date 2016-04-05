package execution

import (
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/output"
	"github.com/mkboudreau/asrt/writer"
)

type Executor struct {
	WorkerCount            int
	ReportInAggregate      bool
	ReportOnlyFailures     bool
	ReportOnlyStateChanges bool
	OutputFormatter        output.ResultFormatter
	OutputWriter           io.Writer
	resultChannel          chan *output.Result
}

func NewExecutor(isAggregate, onlyFailures, onlyStateChange bool, formatter output.ResultFormatter, writer io.Writer, workers int) *Executor {
	return &Executor{
		WorkerCount:            workers,
		ReportInAggregate:      isAggregate,
		ReportOnlyFailures:     onlyFailures,
		ReportOnlyStateChanges: onlyStateChange,
		OutputFormatter:        formatter,
		OutputWriter:           writer,
	}
}

func (executor *Executor) Execute(incomingTargets []*config.Target) int {
	if executor.resultChannel != nil {
		if _, ok := <-executor.resultChannel; !ok {
			close(executor.resultChannel)
		}
		executor.resultChannel = nil
	}

	executor.resultChannel = make(chan *output.Result)
	targetChannel := make(chan *config.Target, executor.WorkerCount)

	go executor.processTargets(targetChannel, executor.resultChannel)

	for _, target := range incomingTargets {
		targetChannel <- target
	}
	close(targetChannel)

	if executor.ReportInAggregate {
		return executor.processAggregatedResult(executor.resultChannel)
	} else {
		return executor.processEachResult(executor.resultChannel)
	}
}

func (executor *Executor) processTargets(incomingTargets <-chan *config.Target, resultChannel chan<- *output.Result) {
	var wg sync.WaitGroup

	for t := range incomingTargets {
		wg.Add(1)
		go func(target *config.Target) {
			execResult := ExecuteWithTimoutAndHeaders(string(target.Method), target.URL, target.Timeout, target.Headers, target.ExpectedStatus)

			result := output.NewResult(execResult.Success(), execResult.Error, strconv.Itoa(execResult.Expected), strconv.Itoa(execResult.Actual), execResult.URL, target.Label)
			if target.Extra != nil {
				result.Extra = target.Extra
			}
			result.Timestamp = output.NewTimeStringForJSON(time.Now())
			resultChannel <- result
			wg.Done()
		}(t)
	}

	wg.Wait()
	close(resultChannel)
}

func (executor *Executor) processEachResult(resultChannel <-chan *output.Result) int {
	exitStatus := 0
	counter := 0
	formatter := executor.OutputFormatter
	w := executor.OutputWriter
	for r := range resultChannel {
		if r.Success && executor.ReportOnlyFailures {
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

func (executor *Executor) processAggregatedResult(resultChannel <-chan *output.Result) int {
	var results []*output.Result
	exitStatus := 0
	formatter := executor.OutputFormatter
	w := executor.OutputWriter
	for r := range resultChannel {
		results = append(results, r)
		if !r.Success {
			exitStatus = 1
		}
	}

	if !executor.ReportOnlyFailures || (exitStatus == 1 && executor.ReportOnlyFailures) {
		reader := formatter.AggregateReader(results)
		writer.WriteToWriter(w, formatter.Header())
		writer.WriteToWriter(w, reader)
		writer.WriteToWriter(w, formatter.Footer())
	}

	writer.DoneWithWriter(w)

	return exitStatus
}
