package writer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HttpWriter implements io.WriteCloser.
// Note: calls to Write only write to an internal buffer.
// Calling Close performs an http.Post
type HttpWriter struct {
	Url    string
	Method string
	Auth   string
	buffer *bytes.Buffer
}

// Creates a new HttpWriter with sensible defaults
func NewHttpWriter(url string) *HttpWriter {
	return &HttpWriter{
		Url:    url,
		Method: "POST",
		Auth:   "",
	}
}

// HttpWriter Write method writes only to an internal buffer.
// Caller must Close the writer in order for a the submission to occur
func (w *HttpWriter) Write(p []byte) (n int, err error) {
	if w.buffer == nil {
		w.buffer = new(bytes.Buffer)
	}

	return w.buffer.Write(p)
}

// Close method does the actual http.Post(...) to Slack
func (w *HttpWriter) Close() error {
	if buf, err := w.buildRequestBody(); err != nil {
		if err != ErrNoContent {
			return fmt.Errorf("could not close http writer: %v", err)
		}
	} else {

		req, err := http.NewRequest(w.Method, w.Url, buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		if w.Auth != "" {
			req.Header.Set("Authorization", w.Auth)
		}

		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
		if _, err := client.Do(req); err != nil {
			return err
		}
	}

	w.reset()
	return nil
}

func (w *HttpWriter) buildRequestBody() (io.Reader, error) {
	if w.buffer == nil || w.buffer.Len() == 0 {
		return nil, ErrNoContent
	}
	data := w.buffer.String()
	_, err := json.Marshal(&data)
	if err != nil {
		return nil, fmt.Errorf("could not create json request body: %v", err)
	}

	//buf := bytes.NewBuffer(jsonBytes)
	return w.buffer, nil
}

func (w *HttpWriter) reset() {
	w.buffer = nil
}

func (http *HttpWriter) HttpMethod(method string) *HttpWriter {
	http.Method = method
	return http
}

func (http *HttpWriter) HttpAuth(auth string) *HttpWriter {
	http.Auth = auth
	return http
}
