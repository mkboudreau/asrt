package output

import (
	"fmt"
	"io"
	"strings"
)

const (
	// pretty strings have three sections:
	// 1. before color
	// 2. test
	// 3. after color

	fmtCsvRaw    string = "%v,%v,%v"
	fmtCsvPretty        = "%v%v%v,%v%v%v,%v%v%v"
	fmtCsvMd            = "*%v*,%v,%v"

	fmtCsvQuietRaw    = "%v"
	fmtCsvQuietPretty = "%v%v%v"
	fmtCsvQuietMd     = "*%v*"

	fmtCsvAggregateRaw    = "%v,%v"
	fmtCsvAggregatePretty = "%v%v%v,%v%v%v"
	fmtCsvAggregateMd     = "*%v*,%v"
)

type CsvResultFormatter struct {
	Modifiers *ResultFormatModifiers
}

func NewCsvResultFormatter(m *ResultFormatModifiers) *CsvResultFormatter {
	return &CsvResultFormatter{Modifiers: m}
}

func (rf *CsvResultFormatter) Header() io.Reader {
	if !rf.Modifiers.Quiet {
		if rf.Modifiers.Markdown {
			if rf.Modifiers.Aggregate {
				return strings.NewReader("*RESULT*,*COUNT*\n")
			} else {
				return strings.NewReader("*RESULT*,*EXPECT*,*URL*\n")
			}
		} else if rf.Modifiers.Pretty {
			if rf.Modifiers.Aggregate {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v,%vCOUNT%v\n", colorYellow, colorReset, colorYellow, colorReset))
			} else {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v,%vEXPECT%v,%vURL%v\n", colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset))
			}
		} else {
			if rf.Modifiers.Aggregate {
				return strings.NewReader("RESULT,COUNT\n")
			} else {
				return strings.NewReader("RESULT,EXPECT,URL\n")
			}
		}
	}
	return strings.NewReader("")
}

func (rf *CsvResultFormatter) Footer() io.Reader {
	return strings.NewReader("\n")
}
func (rf *CsvResultFormatter) RecordSeparator() io.Reader {
	return strings.NewReader("\n")
}

func (rf *CsvResultFormatter) AggregateReader(result []*Result) io.Reader {
	var s string

	switch {
	case rf.Modifiers.Markdown && !rf.Modifiers.Quiet:
		s = csvAggregateResultMarkdownString(result)
	case rf.Modifiers.Markdown && rf.Modifiers.Quiet:
		s = csvAggregateQuietResultMarkdownString(result)
	case !rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = csvAggregateResultString(result)
	case rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = csvAggregateResultPrettyString(result)
	case !rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = csvAggregateQuietResultString(result)
	case rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = csvAggregateQuietResultPrettyString(result)
	}

	return strings.NewReader(s)
}
func (rf *CsvResultFormatter) Reader(result *Result) io.Reader {
	var s string

	switch {
	case rf.Modifiers.Markdown && !rf.Modifiers.Quiet:
		s = csvResultMarkdownString(result)
	case rf.Modifiers.Markdown && rf.Modifiers.Quiet:
		s = csvQuietResultMarkdownString(result)
	case !rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = csvResultString(result)
	case rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = csvResultPrettyString(result)
	case !rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = csvQuietResultString(result)
	case rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = csvQuietResultPrettyString(result)
	}

	return strings.NewReader(s)
}

func csvResultString(result *Result) string {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	return fmt.Sprintf(fmtCsvRaw, status, result.Expected, result.Url)
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

func csvAggregateResultString(results []*Result) string {
	aggResult := newAggregateResult(results)
	return fmt.Sprintf(fmtCsvAggregateRaw, aggResult.Success, aggResult.Success)
}

func csvAggregateResultPrettyString(results []*Result) string {

	aggResult := newAggregateResult(results)

	statusColor := colorGreen
	statusText := statusTextOk
	if !aggResult.Success {
		statusColor = colorRed
		statusText = statusTextNotOk
	}

	bStatus := statusColor
	aStatus := colorReset
	bCount := ""
	aCount := ""

	return fmt.Sprintf(fmtCsvAggregatePretty, bStatus, statusText, aStatus, bCount, aggResult.Count, aCount)
}

func csvQuietResultString(result *Result) string {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	return fmt.Sprintf(fmtCsvQuietRaw, status)
}

func csvQuietResultPrettyString(result *Result) string {
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

	return fmt.Sprintf(fmtCsvQuietPretty, bStatus, statusText, aStatus)
}

func csvAggregateQuietResultString(results []*Result) string {
	aggResult := newAggregateQuietResult(results)
	return fmt.Sprintf(fmtCsvQuietRaw, aggResult.Success)
}

func csvAggregateQuietResultPrettyString(results []*Result) string {
	aggResult := newAggregateQuietResult(results)

	statusColor := colorGreen
	statusText := statusTextOk
	if !aggResult.Success {
		statusColor = colorRed
		statusText = statusTextNotOk
	}

	bStatus := statusColor
	aStatus := colorReset

	return fmt.Sprintf(fmtCsvQuietPretty, bStatus, statusText, aStatus)
}

func csvResultMarkdownString(result *Result) string {
	statusText := statusTextOk
	if !result.Success {
		statusText = statusTextNotOk
	}
	if result.Error != nil {
		statusText = statusTextError
	}

	return fmt.Sprintf(fmtCsvMd, statusText, result.Expected, result.Url)
}

func csvQuietResultMarkdownString(result *Result) string {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	return fmt.Sprintf(fmtCsvQuietMd, status)
}

func csvAggregateResultMarkdownString(results []*Result) string {
	aggResult := newAggregateResult(results)

	statusText := statusTextOk
	if !aggResult.Success {
		statusText = statusTextNotOk
	}

	return fmt.Sprintf(fmtCsvAggregateMd, statusText, aggResult.Count)
}

func csvAggregateQuietResultMarkdownString(results []*Result) string {
	aggResult := newAggregateQuietResult(results)
	return fmt.Sprintf(fmtCsvQuietMd, aggResult.Success)
}
