package execution

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mkboudreau/asrt/config"
	"github.com/stretchr/testify/assert"
)

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
