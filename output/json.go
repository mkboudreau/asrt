package output

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
)

type JsonResultFormatter struct {
	Modifiers *ResultFormatModifiers
}

func NewJsonResultFormatter(m *ResultFormatModifiers) *JsonResultFormatter {
	return &JsonResultFormatter{Modifiers: m}
}

func (rf *JsonResultFormatter) Header() io.Reader {
	if !rf.Modifiers.Aggregate {
		if rf.Modifiers.Pretty {
			return strings.NewReader("[\n")
		} else {
			return strings.NewReader("[")
		}
	}
	return strings.NewReader("")
}

func (rf *JsonResultFormatter) RecordSeparator() io.Reader {
	return strings.NewReader(",")
}

func (rf *JsonResultFormatter) Footer() io.Reader {
	if !rf.Modifiers.Aggregate {
		if rf.Modifiers.Pretty {
			return strings.NewReader("\n]\n")
		} else {
			return strings.NewReader("]\n")
		}

	}
	return strings.NewReader("\n")
}

func (rf *JsonResultFormatter) AggregateReader(results []*Result) io.Reader {
	if rf.Modifiers.Quiet {
		return rf.getReaderForInterface(newAggregateQuietResult(results))
	} else {
		return rf.getReaderForInterface(newAggregateResult(results))
	}
}

func (rf *JsonResultFormatter) Reader(result *Result) io.Reader {
	if rf.Modifiers.Quiet {
		return rf.getReaderForInterface(newQuietResult(result))
	} else {
		return rf.getReaderForInterface(result)
	}
}

func (rf *JsonResultFormatter) getReaderForInterface(obj interface{}) io.Reader {
	var b []byte
	var err error

	if rf.Modifiers.Pretty {
		b, err = jsonResultPrettyString(obj)
	} else {
		b, err = jsonResultString(obj)
	}

	if err != nil {
		log.Printf("Could not get io.Reader for result: %v", err)
		return nil
	}
	return bytes.NewReader(b)
}

func jsonResultString(result interface{}) ([]byte, error) {
	return json.Marshal(result)
}
func jsonResultPrettyString(result interface{}) ([]byte, error) {
	return json.MarshalIndent(result, "", "\t")
}
