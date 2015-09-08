package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CommandMethod string

const (
	MethodGet    CommandMethod = "GET"
	MethodPost                 = "POST"
	MethodPut                  = "PUT"
	MethodDelete               = "DELETE"
	MethodHead                 = "HEAD"
	MethodPatch                = "PATCH"
)

type Target struct {
	Label          string
	Method         CommandMethod
	Timeout        time.Duration
	ExpectedStatus int
	URL            string
	Headers        map[string]string
}

func NewTarget(label string, urlString string, method CommandMethod, expectedStatus int) (*Target, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	t := &Target{
		Label:          label,
		Method:         method,
		URL:            u.String(),
		ExpectedStatus: expectedStatus,
	}

	return t, nil
}

func ParseTarget(targetString string) (*Target, error) {
	url, theRest := extractURL(targetString)
	var method CommandMethod
	var label string
	var statusCode int
	for _, part := range theRest {
		if string(method) == "" && isHttpMethod(part) {
			method = extractMethod(part)
		} else if statusCode == 0 && isStatusCode(part) {
			statusCode = extractStatusCode(part, method)
		} else if label == "" {
			label = extractLabel(part)
		}
	}
	if string(method) == "" {
		method = extractMethod("")
	}
	if statusCode == 0 {
		statusCode = extractStatusCode("", method)
	}

	t, err := NewTarget(label, url, method, statusCode)
	if err != nil {
		return nil, fmt.Errorf("could not create target with url %v: %v", url, err)
	}
	t.Method = method
	return t, nil
}

func extractURL(targetString string) (url string, theRest []string) {
	if strings.HasPrefix(targetString, "\"") && strings.Index(targetString[1:], "\"") != -1 {
		index := strings.Index(targetString[1:], "\"")
		url = targetString[1 : index+1]
		if index+2 < len(targetString) {
			theRest = strings.Split(targetString[index+2:], "|")
		}
	} else if strings.HasPrefix(targetString, "'") && strings.Index(targetString[1:], "'") != -1 {
		index := strings.Index(targetString[1:], "'")
		url = targetString[1 : index+1]
		if index+2 < len(targetString) {
			theRest = strings.Split(targetString[index+2:], "|")
		}
	} else {
		allParts := strings.Split(targetString, "|")
		url = allParts[0]
		if len(allParts) > 1 {
			theRest = allParts[1:]
		}
	}
	return
}

func extractMethod(method string) CommandMethod {
	switch method {
	case string(MethodGet):
		return MethodGet
	case string(MethodPost):
		return MethodPost
	case string(MethodPut):
		return MethodPut
	case string(MethodDelete):
		return MethodDelete
	case string(MethodHead):
		return MethodHead
	case string(MethodPatch):
		return MethodPatch
	default:
		return MethodGet
	}
}

func isHttpMethod(method string) bool {
	switch method {
	case string(MethodGet):
		return true
	case string(MethodPost):
		return true
	case string(MethodPut):
		return true
	case string(MethodDelete):
		return true
	case string(MethodHead):
		return true
	case string(MethodPatch):
		return true
	default:
		return false
	}
}

func extractStatusCode(statusCodeString string, currentMethod CommandMethod) int {
	if statusCode, err := strconv.Atoi(statusCodeString); err == nil {
		return statusCode
	}
	return DefaultHttpStatuses[string(currentMethod)]
}
func isStatusCode(statusCodeString string) bool {
	if statusCode, err := strconv.Atoi(statusCodeString); err == nil && statusCode > 0 && statusCode < 600 {
		return true
	}
	return false
}

func extractLabel(labelString string) string {
	if strings.HasPrefix(labelString, "\"") && strings.HasSuffix(labelString, "\"") {
		return labelString[1 : len(labelString)-1]
	} else if strings.HasPrefix(labelString, "'") && strings.HasSuffix(labelString, "'") {
		return labelString[1 : len(labelString)-1]
	} else {
		return labelString
	}
}
