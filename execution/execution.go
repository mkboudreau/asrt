package execution

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

var DefaultHeaders = map[string]string{
	"Accept":          "application/json",
	"Content-Type":    "application/json",
	"Accept-Encoding": "gzip, deflate",
}
var DefaultTimeout = 0 * time.Second

func Execute(method string, url string, expectation int) (bool, error) {
	return ExecuteWithTimoutAndHeaders(method, url, DefaultTimeout, DefaultHeaders, expectation)
}

func ExecuteWithHeaders(method string, url string, headers map[string]string, expectation int) (bool, error) {
	return ExecuteWithTimoutAndHeaders(method, url, DefaultTimeout, headers, expectation)
}

func ExecuteWithTimeout(method string, url string, timeout time.Duration, expectation int) (bool, error) {
	return ExecuteWithTimoutAndHeaders(method, url, timeout, DefaultHeaders, expectation)
}

func ExecuteWithTimoutAndHeaders(method string, url string, timeout time.Duration, headers map[string]string, expectation int) (bool, error) {
	statusCode, err := execute(method, url, timeout, headers)
	if err != nil {
		return false, err
	}

	if statusCode != expectation {
		log.Println("Status code for %v expected to be %d, but actual is %d", url, expectation, statusCode)
	}

	return (statusCode == expectation), nil
}

func execute(method string, url string, timeout time.Duration, headers map[string]string) (int, error) {
	req, reqErr := buildRequest(method, url, headers)
	if reqErr != nil {
		return 0, reqErr
	}

	resp, respErr := connectWithInsecureCert(req, timeout)
	if respErr != nil {
		return 0, respErr
	}
	defer resp.Body.Close()

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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: timeout}

	return client.Do(req)
}
