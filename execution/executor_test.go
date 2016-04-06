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

type execTestCase struct {
	checkStatusCode  int
	expectedExitCode int
	expectedSuccess  bool
}

var executorTypicalTestCases = []execTestCase{
	{200, 0, true},
	{300, 1, false},
	{400, 1, false},
}

func TestExecutorTypical(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	for _, test := range executorTypicalTestCases {
		writer := new(bytes.Buffer)
		formatter := &testResultFormatter{}
		exec := NewExecutor(false, false, false, formatter, writer, 1)
		target, _ := config.NewTarget("abc", testServer.URL, config.MethodGet, test.checkStatusCode)
		targets := []*config.Target{target}

		actualExitCode := exec.Execute(targets)
		actualResult := formatter.lastResult

		assert.Equal(t, test.expectedExitCode, actualExitCode)

		assert.NotNil(t, actualResult)
		assert.Equal(t, test.expectedSuccess, actualResult.Success)
	}
}

var executorOnlyFailureTestCases = []execTestCase{
	{200, 0, true},
	{300, 1, false},
	{400, 1, false},
}

func TestExecutorOnlyFailure(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	for _, test := range executorOnlyFailureTestCases {
		writer := new(bytes.Buffer)
		formatter := &testResultFormatter{}
		exec := NewExecutor(false, true, false, formatter, writer, 1)
		target, _ := config.NewTarget("abc", testServer.URL, config.MethodGet, test.checkStatusCode)
		targets := []*config.Target{target}

		actualExitCode := exec.Execute(targets)
		actualResult := formatter.lastResult
		actualOutput := writer.String()

		assert.Equal(t, test.expectedExitCode, actualExitCode)
		if test.expectedSuccess {
			assert.Empty(t, actualOutput)
			assert.Nil(t, actualResult)
		} else {
			assert.NotEmpty(t, actualOutput)
			assert.NotNil(t, actualResult)
		}

		if actualResult != nil {
			assert.Equal(t, test.expectedSuccess, actualResult.Success)
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
