package commands

import (
	"github.com/mkboudreau/asrt/execution"
	"github.com/mkboudreau/asrt/output"
	"strconv"
	"sync"
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
