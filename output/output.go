package output

import "io"

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
	Success   bool                   `json:"ok"`
	Error     error                  `json:"error,omitempty"`
	Expected  string                 `json:"expectation,omitempty"`
	Actual    string                 `json:"actual,omitempty"`
	Url       string                 `json:"url,omitempty"`
	Label     string                 `json:"label,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

func NewResult(success bool, err error, expected string, actual string, url string, label string) *Result {
	return &Result{
		Success:  success,
		Error:    err,
		Expected: expected,
		Actual:   actual,
		Url:      url,
		Label:    label,
	}
}

func (r *Result) AddExtra(key string, data interface{}) {
	if r.Extra == nil {
		r.Extra = make(map[string]interface{})
	}
	r.Extra[key] = data
}

type ResultFormatter interface {
	Reader(result *Result) io.Reader
	AggregateReader(result []*Result) io.Reader
	Header() io.Reader
	Footer() io.Reader
	RecordSeparator() io.Reader
}

type ResultFormatModifiers struct {
	Pretty    bool
	Aggregate bool
	NoHeader  bool
	Markdown  bool
}

type quietResult struct {
	Success bool   `json:"ok"`
	Url     string `json:"url,omitempty"`
}

type quietAggregateResult struct {
	Success bool `json:"ok"`
}

type aggregateResult struct {
	Success bool `json:"ok"`
	Count   int  `json:"count,omitempty"`
}

func newQuietResult(result *Result) *quietResult {
	return &quietResult{Success: result.Success, Url: result.Url}
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

type StatusMessager interface {
	StatusMessage() string
}

func (result *Result) StatusCodeActual() string {
	if result.Actual == "" {
		return "n/a"
	}
	return result.Actual
}

func (result *Result) StatusMessage() string {
	if result.Success {
		return statusTextOk
	}
	if result.Error != nil {
		return statusTextError
	}
	return statusTextNotOk
}

func (result *quietResult) StatusMessage() string {
	if result.Success {
		return statusTextOk
	} else {
		return statusTextNotOk
	}
}

func (result *quietAggregateResult) StatusMessage() string {
	if result.Success {
		return statusTextOk
	} else {
		return statusTextNotOk
	}
}

func (result *aggregateResult) StatusMessage() string {
	if result.Success {
		return statusTextOk
	} else {
		return statusTextNotOk
	}
}
