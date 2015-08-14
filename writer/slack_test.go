package writer

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlackWriter(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b, err := ioutil.ReadAll(r.Body); err != nil {
			assert.Fail(t, "Error: %v", err)
		} else {
			s := string(b)
			t.Logf("Submitted Request Body: %v", s)
			assert.False(t, strings.Contains(s, "payload"), "request body should not have the key payload")
			assert.True(t, strings.Contains(s, "text"), "request body should have the text key")
			assert.True(t, strings.HasPrefix(s, "{"), "no opening brace")
			assert.True(t, strings.HasSuffix(s, "}"), "no closing brace")
			assert.True(t, strings.Contains(s, "Hello World"), "does not contain Hello World")
		}
		w.WriteHeader(http.StatusOK)
	}))
	slackWriter := NewSlackWriter(testServer.URL)
	stringReader := strings.NewReader("Hello World")
	count, err := io.Copy(slackWriter, stringReader)
	assert.NotEqual(t, 0, count, "expecting to have written some bytes")
	assert.Nil(t, err, "expecting no error on write")
	assert.Nil(t, slackWriter.Close(), "expecting no error on close")
}

func TestSlackWriterOverrides(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b, err := ioutil.ReadAll(r.Body); err != nil {
			assert.Fail(t, "Error: %v", err)
		} else {
			s := string(b)
			t.Logf("Submitted Request Body: %v", s)
			assert.True(t, strings.Contains(s, "channel"), "does not contain key: channel")
			assert.True(t, strings.Contains(s, "username"), "does not contain key: username")
			assert.True(t, strings.Contains(s, "text"), "does not contain key: text")
			assert.True(t, strings.Contains(s, "icon_url"), "does not contain key: icon_url")

			assert.True(t, strings.Contains(s, "Another Test"), "does not contain Another Test")
			assert.True(t, strings.Contains(s, "@slackbot"), "does not contain @slackbot")
			assert.True(t, strings.Contains(s, "http://someicon"), "does not contain http://someicon")
			assert.True(t, strings.Contains(s, "testuser"), "does not contain testuser")
		}
		w.WriteHeader(http.StatusOK)
	}))
	slackWriter := NewSlackWriter(testServer.URL)
	slackWriter.SlackChannel("@slackbot")
	slackWriter.SlackUser("testuser")
	slackWriter.SlackIconUrl("http://someicon")
	stringReader := strings.NewReader("Another Test")
	count, err := io.Copy(slackWriter, stringReader)
	assert.NotEqual(t, 0, count, "expecting to have written some bytes")
	assert.Nil(t, err, "expecting no error on write")
	assert.Nil(t, slackWriter.Close(), "expecting no error on close")
}
