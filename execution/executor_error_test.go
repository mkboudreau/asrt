package execution

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mkboudreau/asrt/config"
	"github.com/stretchr/testify/assert"
)

func TestExecutorOnlyStateUsingError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

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

	//// CLOSE SERVER
	testServer.Close()
	////

	// third time - failure - change
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

}

func TestExecutorOnlyStateChangeUsingErrorsStartWithError(t *testing.T) {
	shouldClose := true
	var serverPointer *httptest.Server

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldClose {
			serverPointer.CloseClientConnections()
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer testServer.Close()
	serverPointer = testServer

	writer := new(bytes.Buffer)
	formatter := &testResultFormatter{}
	exec := NewExecutor(false, false, true, formatter, writer, 1)
	target, _ := config.NewTarget("abc", testServer.URL, config.MethodGet, 200)
	targets := []*config.Target{target}

	// first time - failure
	actualExitCode := exec.Execute(targets)
	actualResult := formatter.lastResult
	actualOutput := writer.String()

	assert.Equal(t, 1, actualExitCode)
	assert.NotEmpty(t, actualOutput)
	assert.NotNil(t, actualResult)
	assert.False(t, actualResult.Success)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// second time - failure - no change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	////// CHANGE THE STATE
	shouldClose = false

	// third time - success - change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.NotEmpty(t, actualOutput)
	assert.NotNil(t, actualResult)
	assert.True(t, actualResult.Success)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	// fourth time - success - no change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

	// -- reset --
	formatter.lastResult = nil
	writer.Reset()

	///// TELL SERVER TO CLOSE CLIENT CONNS
	shouldClose = true

	// fifth time - failure - change
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

	// sixth time - failure - no change
	actualExitCode = exec.Execute(targets)
	actualResult = formatter.lastResult
	actualOutput = writer.String()

	assert.Equal(t, 0, actualExitCode)
	assert.Empty(t, actualOutput)
	assert.Nil(t, actualResult)

}
