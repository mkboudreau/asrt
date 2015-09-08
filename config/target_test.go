package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectation struct {
	url    string
	label  string
	method string
	status int
}

var testCasesForTargetParsing = []struct {
	target string
	expect *expectation
}{
	{"www.yahoo.com", &expectation{"http://www.yahoo.com", "", "GET", 200}},
	{"http://www.yahoo.com", &expectation{"http://www.yahoo.com", "", "GET", 200}},
	{"https://www.yahoo.com", &expectation{"https://www.yahoo.com", "", "GET", 200}},
	{"'www.yahoo.com?abc=123|456'", &expectation{"http://www.yahoo.com?abc=123|456", "", "GET", 200}},
	{"\"www.yahoo.com?abc=123|456\"", &expectation{"http://www.yahoo.com?abc=123|456", "", "GET", 200}},
	{"'www.yahoo.com?abc=123|456'|\"label me\"", &expectation{"http://www.yahoo.com?abc=123|456", "label me", "GET", 200}},
	{"\"www.yahoo.com?abc=123|456\"|300|what", &expectation{"http://www.yahoo.com?abc=123|456", "what", "GET", 300}},
	{"www.yahoo.com|GET", &expectation{"http://www.yahoo.com", "", "GET", 200}},
	{"www.yahoo.com|GET|202", &expectation{"http://www.yahoo.com", "", "GET", 202}},
	{"www.yahoo.com|GET|203|Hello", &expectation{"http://www.yahoo.com", "Hello", "GET", 203}},
	{"www.yahoo.com|201", &expectation{"http://www.yahoo.com", "", "GET", 201}},
	{"www.yahoo.com|POST", &expectation{"http://www.yahoo.com", "", "POST", 201}},
	{"www.yahoo.com|Test", &expectation{"http://www.yahoo.com", "Test", "GET", 200}},
	{"www.yahoo.com|\"Test Label\"", &expectation{"http://www.yahoo.com", "Test Label", "GET", 200}},
}

func TestTargetParsing(t *testing.T) {
	for _, tc := range testCasesForTargetParsing {
		target, err := ParseTarget(tc.target)
		assert.Nil(t, err)
		assert.NotNil(t, target)
		assert.Equal(t, tc.expect.url, target.URL)
		assert.Equal(t, tc.expect.label, target.Label)
		assert.Equal(t, tc.expect.status, target.ExpectedStatus)
		assert.Equal(t, tc.expect.method, string(target.Method))
	}
}
