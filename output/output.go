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
	Success  bool   `json:"success,omitempty"`
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
}
