package execution

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/output"
	"github.com/stretchr/testify/assert"
)

func TestExecutorTypical(t *testing.T) {
	var statusCodesToTest = getIntsInRange(100, 599)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	for _, statusCodeTest := range statusCodesToTest {
		writer := new(bytes.Buffer)
		formatter := &testResultFormatter{}
		exec := NewExecutor(false, false, false, formatter, writer, 1)
		target, _ := config.NewTarget("abc", testServer.URL, config.MethodGet, statusCodeTest)
		targets := []*config.Target{target}

		actualExitCode := exec.Execute(targets)
		actualResult := formatter.lastResult
		actualOutput := writer.String()

		assert.NotEmpty(t, actualOutput)
		assert.NotNil(t, actualResult)

		if statusCodeTest == 200 {
			assert.Equal(t, 0, actualExitCode)
			assert.True(t, actualResult.Success)
		} else {
			assert.NotEmpty(t, actualOutput)
			assert.NotNil(t, actualResult)
			assert.Equal(t, 1, actualExitCode)
			assert.False(t, actualResult.Success)
		}
	}
}

func getIntsInRange(begin, end int) []int {
	var ints []int
	for i := begin; i < end; i++ {
		ints = append(ints, i)
	}
	return ints
}

func TestExecutorOnlyFailure(t *testing.T) {
	var statusCodesToTest = getIntsInRange(100, 599)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	for _, statusCodeTest := range statusCodesToTest {
		writer := new(bytes.Buffer)
		formatter := &testResultFormatter{}
		exec := NewExecutor(false, true, false, formatter, writer, 1)
		target, _ := config.NewTarget("abc", testServer.URL, config.MethodGet, statusCodeTest)
		targets := []*config.Target{target}

		actualExitCode := exec.Execute(targets)
		actualResult := formatter.lastResult
		actualOutput := writer.String()

		if statusCodeTest == 200 {
			assert.Empty(t, actualOutput)
			assert.Nil(t, actualResult)
			assert.Equal(t, 0, actualExitCode)
		} else {
			assert.NotEmpty(t, actualOutput)
			assert.NotNil(t, actualResult)
			assert.Equal(t, 1, actualExitCode)
			assert.False(t, actualResult.Success)
		}
	}
}

func TestExecutorOnlyStatusChange(t *testing.T) {
	statusToReturn := http.StatusOK
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusToReturn)
	}))
	defer testServer.Close()

	writer := new(bytes.Buffer)
	formatter := &testResultFormatter{}
	exec := NewExecutor(false, false, true, formatter, writer, 1)
	target, _ := config.NewTarget("abc", testServer.URL, config.MethodGet, 200)
	targets := []*config.Target{target}

	// first time - success
	actualExitCode := exec.Execute(targets)
	actualResult := formatter.lastResult
	actualOutput := writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.NotEmpty(t, actualOutput)
	assert.NotNil(t, actualResult)
	assert.True(t, actualResult.Success)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// second time - success - no change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// third time - failure - change
	statusToReturn = http.StatusBadRequest
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 1, actualExitCode)
	assert.NotEmpty(t, actualOutput)
	assert.NotNil(t, actualResult)
	assert.False(t, actualResult.Success)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// fourth time - failure - no change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// fifth time - success - change
	statusToReturn = http.StatusOK
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.NotEmpty(t, actualOutput)
	assert.NotNil(t, actualResult)
	assert.True(t, actualResult.Success)
}

func TestExecutorOnlyStatusChangeAndOnlyFailure(t *testing.T) {
	statusToReturn := http.StatusOK
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusToReturn)
	}))
	defer testServer.Close()

	writer := new(bytes.Buffer)
	formatter := &testResultFormatter{}
	exec := NewExecutor(false, true, true, formatter, writer, 1)
	target, _ := config.NewTarget("abc", testServer.URL, config.MethodGet, 200)
	targets := []*config.Target{target}

	// first time - success
	actualExitCode := exec.Execute(targets)
	actualResult := formatter.lastResult
	actualOutput := writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// second time - success - no change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// third time - failure - change
	statusToReturn = http.StatusBadRequest
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 1, actualExitCode)
	assert.NotEmpty(t, actualOutput)
	assert.NotNil(t, actualResult)
	assert.False(t, actualResult.Success)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// fourth time - failure - no change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// fifth time - success - change
	statusToReturn = http.StatusOK
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)
}

type testResultFormatter struct {
	lastResult          *output.Result
	lastAggregateResult []*output.Result
}

func (rf *testResultFormatter) Reader(result *output.Result) io.Reader {
	rf.lastResult = result
	return strings.NewReader("true")
}
func (rf *testResultFormatter) AggregateReader(result []*output.Result) io.Reader {
	rf.lastAggregateResult = result
	return strings.NewReader("true")
}
func (rf *testResultFormatter) Header() io.Reader {
	return strings.NewReader("")
}
func (rf *testResultFormatter) Footer() io.Reader {
	return strings.NewReader("")
}
func (rf *testResultFormatter) RecordSeparator() io.Reader {
	return strings.NewReader("")
}
