package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type expectation struct {
	url     string
	label   string
	method  string
	status  int
	headers map[string]string
}

var testCasesForTargetParsing = []struct {
	target string
	expect *expectation
}{
	{"www.yahoo.com", &expectation{"http://www.yahoo.com", "", "GET", 200, map[string]string{}}},
	{"http://www.yahoo.com", &expectation{"http://www.yahoo.com", "", "GET", 200, map[string]string{}}},
	{"https://www.yahoo.com", &expectation{"https://www.yahoo.com", "", "GET", 200, map[string]string{}}},
	{"'www.yahoo.com?abc=123|456'", &expectation{"http://www.yahoo.com?abc=123|456", "", "GET", 200, map[string]string{}}},
	{"\"www.yahoo.com?abc=123|456\"", &expectation{"http://www.yahoo.com?abc=123|456", "", "GET", 200, map[string]string{}}},
	{"'www.yahoo.com?abc=123|456'|\"label me\"", &expectation{"http://www.yahoo.com?abc=123|456", "label me", "GET", 200, nil}},
	{"\"www.yahoo.com?abc=123|456\"|300|what", &expectation{"http://www.yahoo.com?abc=123|456", "what", "GET", 300, nil}},
	{"www.yahoo.com|GET", &expectation{"http://www.yahoo.com", "", "GET", 200, nil}},
	{"www.yahoo.com|GET|202", &expectation{"http://www.yahoo.com", "", "GET", 202, nil}},
	{"www.yahoo.com|GET|203|Hello", &expectation{"http://www.yahoo.com", "Hello", "GET", 203, nil}},
	{"www.yahoo.com|201", &expectation{"http://www.yahoo.com", "", "GET", 201, nil}},
	{"www.yahoo.com|POST", &expectation{"http://www.yahoo.com", "", "POST", 201, nil}},
	{"www.yahoo.com|Test", &expectation{"http://www.yahoo.com", "Test", "GET", 200, nil}},
	{"www.yahoo.com|\"Test Label\"", &expectation{"http://www.yahoo.com", "Test Label", "GET", 200, nil}},
	{"www.yahoo.com|{H}Test", &expectation{"http://www.yahoo.com", "", "GET", 200, map[string]string{"Test": ""}}},
	{"www.yahoo.com|{H}Test:", &expectation{"http://www.yahoo.com", "", "GET", 200, map[string]string{"Test": ""}}},
	{"www.yahoo.com|{H}Test: Value", &expectation{"http://www.yahoo.com", "", "GET", 200, map[string]string{"Test": "Value"}}},
	{"www.yahoo.com|{H}Test: Value|{H}Another: V", &expectation{"http://www.yahoo.com", "", "GET", 200, map[string]string{"Test": "Value", "Another": "V"}}},
	{"www.yahoo.com|{H}Test: Value|{H}\"Another: X Y Z\"", &expectation{"http://www.yahoo.com", "", "GET", 200, map[string]string{"Test": "Value", "Another": "X Y Z"}}},
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
		assertMapContains(t, target.Headers, tc.expect.headers)
	}
}

func assertMapContains(t *testing.T, targetMap, checkMap map[string]string) {
	if checkMap == nil {
		return
	}
	if targetMap == nil && len(checkMap) > 0 {
		assert.Fail(t, "No Headers in Parsed Target")
	}
	for k, v := range checkMap {
		assert.Equal(t, v, targetMap[k])
	}
}
