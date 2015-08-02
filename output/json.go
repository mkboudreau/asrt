package output

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type JsonResultFormatter struct {
	Modifiers *ResultFormatModifiers
}

func NewJsonResultFormatter(m *ResultFormatModifiers) *JsonResultFormatter {
	return &JsonResultFormatter{Modifiers: m}
}

func (rf *JsonResultFormatter) AggregateReader(result []*Result) io.Reader {
	return rf.Reader(result[0])
}

func (rf *JsonResultFormatter) Reader(result *Result) io.Reader {
	var b []byte
	var err error
	if rf.Modifiers.Pretty {
		b, err = jsonResultPrettyString(result)
	} else {
		b, err = jsonResultString(result)
	}
	if err != nil {
		log.Printf("Could not get io.Reader for result: %v", err)
		return nil
	}
	return bytes.NewReader(b)
}

func jsonResultPrettyString(result *Result) ([]byte, error) {
	return json.MarshalIndent(result, "", "\t")
}
func jsonResultString(result *Result) ([]byte, error) {
	return json.Marshal(result)
}
