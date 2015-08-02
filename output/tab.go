package output

import (
	"fmt"
	"io"
	"strings"
)

const (
	fmtTabRaw string = "%v\t%v\t%v\n"

	// three sections: status text, expected value, url
	// for each section:
	// structure: color before, data, color after, tab
	fmtTabPretty = "%v%v%v\t%v%v%v\t%v%v%v\n"
)

type TabResultFormatter struct {
	Modifiers *ResultFormatModifiers
}

func NewTabResultFormatter(m *ResultFormatModifiers) *TabResultFormatter {
	return &TabResultFormatter{Modifiers: m}
}

func (rf *TabResultFormatter) AggregateReader(result []*Result) io.Reader {
	return rf.Reader(result[0])
}

func (rf *TabResultFormatter) Reader(result *Result) io.Reader {
	var s string
	if rf.Modifiers.Pretty {
		s = tabResultPrettyString(result)
	} else {
		s = tabResultString(result)
	}

	return strings.NewReader(s)
}
func tabResultString(result *Result) string {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	return fmt.Sprintf(fmtTabRaw, status, result.Expected, result.Url)
}

func tabResultPrettyString(result *Result) string {
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

	return fmt.Sprintf(fmtTabPretty, bStatus, statusText, aStatus, bExpected, result.Expected, aExpected, bUrl, result.Url, aUrl)
}
