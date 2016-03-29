package output

import (
	"bytes"
	"github.com/mkboudreau/asrt/log"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

/*
func NewResult(success bool, err error, expected string, actual string, url string, label string) *Result {
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
*/
var emptyModifiersForTesting *ResultFormatModifiers = new(ResultFormatModifiers)

func TestTemplateAllFields(t *testing.T) {
	log.TurnOnLogging()
	option := "template={{.Success }} {{ .Error }} {{ .Expected }}, {{ .Actual }},{{ .Url }}, {{ .Label }}"
	expected := "true <nil> hello, hello,http://hello/world, my label"
	formatter := NewTemplateResultFormatter(option, emptyModifiersForTesting)
	result := NewResult(true, nil, "hello", "hello", "http://hello/world", "my label")
	reader := formatter.Reader(result)
	assert.Equal(t, expected, readerToString(reader))
}

func TestTemplateAllFieldsPlusOneNotRecognized(t *testing.T) {
	log.TurnOnLogging()
	option := "template={{.Success }} {{ .Error }} {{ .Expected }}, {{ .Actual }},{{ .Url }}, {{ .Label }} {{ .WHAT }}"
	expected := "true <nil> hello, hello,http://hello/world, my label "
	formatter := NewTemplateResultFormatter(option, emptyModifiersForTesting)
	result := NewResult(true, nil, "hello", "hello", "http://hello/world", "my label")
	reader := formatter.Reader(result)
	assert.Equal(t, expected, readerToString(reader))
}

func TestTemplateAllFieldsWithMultipleLines(t *testing.T) {
	log.TurnOnLogging()
	option := `template=[ {{.Success}} ] {{ .Label }}, {{ .Url}}
Expected {{ .Expected }}; Got {{ .Actual }}`
	expected := "[ true ] my label, http://hello/world\nExpected hello; Got hello"
	formatter := NewTemplateResultFormatter(option, emptyModifiersForTesting)
	result := NewResult(true, nil, "hello", "hello", "http://hello/world", "my label")
	reader := formatter.Reader(result)
	assert.Equal(t, expected, readerToString(reader))
}

func TestTemplateConditional(t *testing.T) {
	log.TurnOnLogging()
	option := `template=[ {{if .Success}}ok{{else}}fail{{end}} ] {{ .Label }}, {{ .Url}}
Expected {{ .Expected }}; Got {{ .Actual }}`
	expected := "[ ok ] my label, http://hello/world\nExpected hello; Got hello"
	formatter := NewTemplateResultFormatter(option, emptyModifiersForTesting)
	result := NewResult(true, nil, "hello", "hello", "http://hello/world", "my label")
	reader := formatter.Reader(result)
	assert.Equal(t, expected, readerToString(reader))
}

func readerToString(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
}
