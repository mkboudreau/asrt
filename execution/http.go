package execution

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

var DefaultHeaders = map[string]string{
	"Accept":          "application/json",
	"Content-Type":    "application/json",
	"Accept-Encoding": "gzip, deflate",
}
var DefaultTimeout = 0 * time.Second

type ExecutionResult struct {
	URL      string
	Method   string
	Expected int
	Actual   int
	Error    error
}

func (r *ExecutionResult) Success() bool {
	return r.Error == nil && r.Expected == r.Actual
}

func Execute(method string, url string, expectation int) *ExecutionResult {
	return ExecuteWithTimoutAndHeaders(method, url, DefaultTimeout, DefaultHeaders, expectation)
}

func ExecuteWithHeaders(method string, url string, headers map[string]string, expectation int) *ExecutionResult {
	return ExecuteWithTimoutAndHeaders(method, url, DefaultTimeout, headers, expectation)
}

func ExecuteWithTimeout(method string, url string, timeout time.Duration, expectation int) *ExecutionResult {
	return ExecuteWithTimoutAndHeaders(method, url, timeout, DefaultHeaders, expectation)
}

func ExecuteWithTimoutAndHeaders(method string, url string, timeout time.Duration, headers map[string]string, expectation int) *ExecutionResult {
	statusCode, err := execute(method, url, timeout, headers)

	return &ExecutionResult{
		URL:      url,
		Method:   method,
		Expected: expectation,
		Actual:   statusCode,
		Error:    err,
	}
}

func execute(method string, url string, timeout time.Duration, headers map[string]string) (int, error) {
	req, reqErr := buildRequest(method, url, headers)
	if reqErr != nil {
		return 0, reqErr
	}

	resp, respErr := connectWithInsecureCert(req, timeout)
	log.Println("Connecting")
	log.Printf(" - Request: %+v/n", req)
	log.Printf(" - Response: %+v\n", resp)
	log.Printf(" - Error: %+v\n", respErr)
	if respErr != nil {
		return 0, respErr
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return resp.StatusCode, nil
}

func buildRequest(method string, url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	header := http.Header{}
	for key, value := range headers {
		header.Add(key, value)
	}
	req.Header = header

	return req, nil
}

func connectWithInsecureCert(req *http.Request, timeout time.Duration) (*http.Response, error) {
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: timeout}

	return client.Do(req)
}
