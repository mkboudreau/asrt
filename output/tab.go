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

	fmtTabRaw    string = "%v\t%v\t%v\n"
	fmtTabPretty        = "%v%v%v\t%v%v%v\t%v%v%v\n"

	fmtTabQuietRaw    = "%v\n"
	fmtTabQuietPretty = "%v%v%v\n"

	fmtTabAggregateRaw    = "%v\t%v\n"
	fmtTabAggregatePretty = "%v%v%v\t%v%v%v\n"
)

type TabResultFormatter struct {
	Modifiers *ResultFormatModifiers
}

func NewTabResultFormatter(m *ResultFormatModifiers) *TabResultFormatter {
	return &TabResultFormatter{Modifiers: m}
}

func (rf *TabResultFormatter) AggregateReader(result []*Result) io.Reader {
	var s string

	switch {
	case !rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = tabAggregateResultString(result)
	case rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = tabAggregateResultPrettyString(result)
	case !rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = tabAggregateQuietResultString(result)
	case rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = tabAggregateQuietResultPrettyString(result)
	}

	return strings.NewReader(s)
}

func (rf *TabResultFormatter) Reader(result *Result) io.Reader {
	var s string

	switch {
	case !rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = tabResultString(result)
	case rf.Modifiers.Pretty && !rf.Modifiers.Quiet:
		s = tabResultPrettyString(result)
	case !rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = tabQuietResultString(result)
	case rf.Modifiers.Pretty && rf.Modifiers.Quiet:
		s = tabQuietResultPrettyString(result)
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

func tabQuietResultString(result *Result) string {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	return fmt.Sprintf(fmtTabQuietRaw, status)
}

func tabQuietResultPrettyString(result *Result) string {
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

	return fmt.Sprintf(fmtTabQuietPretty, bStatus, statusText, aStatus)
}

func tabAggregateResultString(results []*Result) string {
	aggResult := newAggregateResult(results)
	return fmt.Sprintf(fmtTabAggregateRaw, aggResult.Success, aggResult.Success)
}

func tabAggregateResultPrettyString(results []*Result) string {

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

	return fmt.Sprintf(fmtTabAggregatePretty, bStatus, statusText, aStatus, bCount, aggResult.Count, aCount)
}
func tabAggregateQuietResultString(results []*Result) string {
	aggResult := newAggregateQuietResult(results)
	return fmt.Sprintf(fmtTabQuietRaw, aggResult.Success)
}

func tabAggregateQuietResultPrettyString(results []*Result) string {
	aggResult := newAggregateQuietResult(results)

	statusColor := colorGreen
	statusText := statusTextOk
	if !aggResult.Success {
		statusColor = colorRed
		statusText = statusTextNotOk
	}

	bStatus := statusColor
	aStatus := colorReset

	return fmt.Sprintf(fmtTabQuietPretty, bStatus, statusText, aStatus)
}
