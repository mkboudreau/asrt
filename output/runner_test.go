package output

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFormat int

const (
	jsonFormat testFormat = iota
	csvFormat
	tabFormat
)

type resultFormatTestCase struct {
	expect    string
	results   []*Result
	format    testFormat
	modifiers *ResultFormatModifiers
}

func runResultFormatTestCase(t *testing.T, testcase *resultFormatTestCase) {
	expectation := testcase.expect
	results := testcase.results

	var rf ResultFormatter
	switch testcase.format {
	case tabFormat:
		rf = NewTabResultFormatter(testcase.modifiers)
	case csvFormat:
		rf = NewCsvResultFormatter(testcase.modifiers)
	}
	//rf := NewCsvResultFormatter(testcase.modifiers)

	if testcase.modifiers.Aggregate {
		reader := rf.AggregateReader(results)
		assert.NotNil(t, reader)
		assert.Equal(t, expectation, readersToString(reader))
	} else {
		readers := make([]io.Reader, len(results))
		for i, result := range results {
			readers[i] = rf.Reader(result)
		}
		assert.NotNil(t, readers)
		assert.Equal(t, expectation, readersToString(readers...))
	}
}

func readersToString(readers ...io.Reader) string {
	buf := new(bytes.Buffer)

	for _, r := range readers {
		buf.ReadFrom(r)
	}
	return buf.String()
}
