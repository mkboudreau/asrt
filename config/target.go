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
	Extra          map[string]interface{}
}

// NewTarget creates a new config.Target object with the required fields.
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

// AddExtra is for writers of TargetConfigurers. This enables a way to attach data or track data from a target to its output.
func (t *Target) AddExtra(key string, data interface{}) {
	if t.Extra == nil {
		t.Extra = make(map[string]interface{})
	}
	t.Extra[key] = data
}

// ParseTarget takes a string of format <url>|<method>|<status_code>|<label>|<header> and parses it into a config.Target object. URL is the only required field.
func ParseTarget(targetString string) (*Target, error) {
	url, theRest := extractURL(targetString)
	var method CommandMethod
	var label string
	var statusCode int
	var headers = make(map[string]string)

	for _, part := range theRest {
		if string(method) == "" && isHttpMethod(part) {
			method = extractMethod(part)
		} else if statusCode == 0 && isStatusCode(part) {
			statusCode = extractStatusCode(part, method)
		} else if !isHeader(part) && label == "" {
			label = extractLabel(part)
		} else if isHeader(part) {
			headerKey, headerValue := extractHeader(part)
			if headerKey != "" {
				headers[headerKey] = headerValue
			}
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
	if len(headers) > 0 {
		t.Headers = headers
	}

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

func isHeader(header string) bool {
	return strings.HasPrefix(header, "{H}")
}

func extractHeader(headerWithPrefix string) (string, string) {
	headerString := strings.TrimPrefix(headerWithPrefix, "{H}")

	header := extractQuotedString(headerString)
	headerParts := strings.SplitN(header, ":", 2)
	if len(headerParts) == 1 {
		return header, ""
	} else if len(headerParts) != 2 {
		return "", ""
	}

	return strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1])
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
	return extractQuotedString(labelString)
}
func extractQuotedString(quoted string) string {
	if strings.HasPrefix(quoted, "\"") && strings.HasSuffix(quoted, "\"") {
		return quoted[1 : len(quoted)-1]
	} else if strings.HasPrefix(quoted, "'") && strings.HasSuffix(quoted, "'") {
		return quoted[1 : len(quoted)-1]
	} else {
		return quoted
	}
}
