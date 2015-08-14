package output

import (
	"io"
)

const (
	colorGreen  string = "\033[1;32m"
	colorRed           = "\033[1;31m"
	colorYellow        = "\033[1;33m"
	colorBlue          = "\033[1;34m"
	colorReset         = "\033[0m"

	statusTextOk    string = "[ok]"
	statusTextNotOk string = "[!ok]"
	statusTextError string = "[err]"
)

type Result struct {
	Success  bool   `json:"ok"`
	Error    error  `json:"error,omitempty"`
	Expected string `json:"expectation,omitempty"`
	Url      string `json:"url,omitempty"`
}

func NewResult(success bool, err error, expected string, url string) *Result {
	return &Result{
		Success:  success,
		Error:    err,
		Expected: expected,
		Url:      url,
	}
}

type ResultFormatter interface {
	Reader(result *Result) io.Reader
	AggregateReader(result []*Result) io.Reader
}

type ResultFormatModifiers struct {
	Pretty    bool
	Aggregate bool
	Quiet     bool
	Markdown  bool
}

type quietResult struct {
	Success bool `json:"ok"`
}

type quietAggregateResult struct {
	Success bool `json:"ok"`
}

type aggregateResult struct {
	Success bool `json:"ok"`
	Count   int  `json:"count,omitempty"`
}

func newQuietResult(result *Result) *quietResult {
	return &quietResult{Success: result.Success}
}

func newAggregateQuietResult(results []*Result) *quietAggregateResult {
	success := true
	for _, r := range results {
		if !r.Success {
			success = false
		}
	}
	return &quietAggregateResult{Success: success}
}

func newAggregateResult(results []*Result) *aggregateResult {
	success := true
	for _, r := range results {
		if !r.Success {
			success = false
		}
	}
	return &aggregateResult{Success: success, Count: len(results)}
}
