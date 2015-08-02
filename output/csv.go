package output

import (
	"fmt"
	"io"
	"strings"
)

const (
	fmtCsvRaw string = "%v,%v,%v\n"

	// three sections: status text, expected value, url
	// for each section:
	// structure: color before, data, color after, comma
	fmtCsvPretty = "%v%v%v,%v%v%v,%v%v%v\n"
)

type CsvResultFormatter struct {
	Modifiers *ResultFormatModifiers
}

func NewCsvResultFormatter(m *ResultFormatModifiers) *CsvResultFormatter {
	return &CsvResultFormatter{Modifiers: m}
}
func (rf *CsvResultFormatter) AggregateReader(result []*Result) io.Reader {
	return rf.Reader(result[0])
}
func (rf *CsvResultFormatter) Reader(result *Result) io.Reader {
	var s string
	if rf.Modifiers.Pretty {
		s = csvResultPrettyString(result)
	} else {
		s = csvResultString(result)
	}

	return strings.NewReader(s)
}

func csvResultPrettyString(result *Result) string {
	statusColor := colorGreen
	statusText := statusTextOk
	if !result.Success {
		statusColor = colorRed
		statusText = statusTextNotOk
	}
	if result.Error != nil {
		statusColor = colorRed
		statusText = statusTextError
	}

	bStatus := statusColor
	aStatus := colorReset
	bExpected := ""
	aExpected := ""
	bUrl := ""
	aUrl := ""

	return fmt.Sprintf(fmtCsvPretty, bStatus, statusText, aStatus, bExpected, result.Expected, aExpected, bUrl, result.Url, aUrl)
}
func csvResultString(result *Result) string {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	return fmt.Sprintf(fmtCsvRaw, status, result.Expected, result.Url)
}
