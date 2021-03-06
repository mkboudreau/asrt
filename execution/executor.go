package execution

import (
	"fmt"
	"io"
	"log"
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
	lastestState           map[stateKey]stateValue
	autoclose              bool
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
		lastestState:           make(map[stateKey]stateValue),
		autoclose:              true,
	}
}

func (executor *Executor) SetAutoClose() {
	executor.autoclose = true
}

func (executor *Executor) SetNoAutoClose() {
	executor.autoclose = false
}

func (executor Executor) Execute(incomingTargets []*config.Target) int {
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
			executor.updateState(r)
			continue
		}
		if executor.ReportOnlyStateChanges {
			if !executor.didStateChange(r) {
				//no change... continue
				log.Printf("no change with result")
				executor.updateState(r)
				continue
			}
			log.Printf("found change with result, sending to writer")
			executor.updateState(r)
		}

		if counter == 0 {
			writer.WriteToWriter(w, formatter.Header())
		} else if counter != 0 {
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

	if executor.autoclose {
		writer.DoneWithWriter(w)
	}

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

	if executor.autoclose {
		writer.DoneWithWriter(w)
	}

	return exitStatus
}

func (executor *Executor) didStateChange(result *output.Result) bool {
	k, v := extractStateKeyValueFromResult(result)

	lv, ok := executor.lastestState[*k]
	if !ok {
		log.Printf("Target %v First State Capture: %v", k, v)
		return true
	}

	if *v != lv {
		log.Printf("Target %v Changed from [ %v ] to [ %v ]", k, lv, *v)
	} else {
		log.Printf("Target %v No Change [ %v ] ", k, *v)
	}

	return *v != lv
}

func (executor *Executor) updateState(result *output.Result) {
	k, v := extractStateKeyValueFromResult(result)
	executor.lastestState[*k] = *v
}

type stateKey struct {
	Url, Label string
}

func (s stateKey) String() string {
	return fmt.Sprintf("Label: %v | URL: %v", s.Label, s.Url)
}

type stateValue struct {
	Success  bool
	Expected string
	Actual   string
}

func (s stateValue) String() string {
	return fmt.Sprintf("Success: %v | Expected: %v | Actual: %v", s.Success, s.Expected, s.Actual)
}

func extractStateKeyValueFromResult(result *output.Result) (*stateKey, *stateValue) {
	k := &stateKey{Url: result.Url, Label: result.Label}
	v := &stateValue{
		Success:  result.Success,
		Expected: result.Expected,
		Actual:   result.Actual,
	}
	return k, v
}
