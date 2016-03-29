package output

import (
	"fmt"
	"io"
	"strings"
)

func getResultStringWithSeparator(result *Result, modifiers *ResultFormatModifiers, separator string) io.Reader {
	var s string

	switch {
	case modifiers.Markdown && !modifiers.NoHeader:
		s = separatorResultForMarkdown(result, separator)
	case modifiers.Markdown && modifiers.NoHeader:
		s = separatorResultForQuietMarkdown(result, separator)
	case !modifiers.Pretty && !modifiers.NoHeader:
		s = separatorResult(result, separator)
	case modifiers.Pretty && !modifiers.NoHeader:
		s = separatorResultForPretty(result, separator)
	case !modifiers.Pretty && modifiers.NoHeader:
		s = separatorResultForQuiet(result, separator)
	case modifiers.Pretty && modifiers.NoHeader:
		s = separatorResultForQuietPretty(result, separator)
	}

	return strings.NewReader(s)
}

func getResultStringAggregateWithSeparator(results []*Result, modifiers *ResultFormatModifiers, separator string) io.Reader {
	var s string

	switch {
	case modifiers.Markdown && !modifiers.NoHeader:
		s = separatorAggregateResultForMarkdown(results, separator)
	case modifiers.Markdown && modifiers.NoHeader:
		s = separatorAggregateResultForQuietMarkdown(results, separator)
	case !modifiers.Pretty && !modifiers.NoHeader:
		s = separatorAggregateResult(results, separator)
	case modifiers.Pretty && !modifiers.NoHeader:
		s = separatorAggregateResultForPretty(results, separator)
	case !modifiers.Pretty && modifiers.NoHeader:
		s = separatorAggregateResultForQuiet(results, separator)
	case modifiers.Pretty && modifiers.NoHeader:
		s = separatorAggregateResultForQuietPretty(results, separator)
	}

	return strings.NewReader(s)
}

// Normal Result
// +pretty
// +quiet
func separatorResultForQuietPretty(result *Result, separator string) string {
	statusColor := colorGreen
	if !result.Success || result.Error != nil {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v%v%v", statusColor, result.StatusMessage(), colorReset, separator, result.Url)
}

// Normal Result
func separatorResult(result *Result, separator string) string {
	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v", result.StatusMessage(), separator, result.Expected, separator, result.StatusCodeActual(), separator, result.Label, separator, result.Url)
}

// Normal Result
// +pretty
func separatorResultForPretty(result *Result, separator string) string {
	statusColor := colorGreen
	if !result.Success || result.Error != nil {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v", statusColor, result.StatusMessage(), colorReset, separator, result.Expected, separator, result.StatusCodeActual(), separator, result.Label, separator, result.Url)
}

// Normal Result
// +quiet
func separatorResultForQuiet(result *Result, separator string) string {
	return fmt.Sprintf("%v%v%v", result.StatusMessage(), separator, result.Url)
}

// Normal Result
// +markdown
func separatorResultForMarkdown(result *Result, separator string) string {
	return fmt.Sprintf("*%v*%v%v%v%v%v%v%v%v", result.StatusMessage(), separator, result.Expected, separator, result.StatusCodeActual(), separator, result.Label, separator, result.Url)
}

// Normal Result
// +markdown
// +quiet
func separatorResultForQuietMarkdown(result *Result, separator string) string {
	return fmt.Sprintf("*%v*%v%v", result.StatusMessage(), separator, result.Url)
}

// Aggregate Result
func separatorAggregateResult(results []*Result, separator string) string {
	aggResult := newAggregateResult(results)
	return fmt.Sprintf("%v%v%v", aggResult.StatusMessage(), separator, aggResult.Count)
}

// Aggregate Result
// +pretty
func separatorAggregateResultForPretty(results []*Result, separator string) string {
	aggResult := newAggregateResult(results)

	statusColor := colorGreen
	if !aggResult.Success {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v%v%v", statusColor, aggResult.StatusMessage(), colorReset, separator, aggResult.Count)
}

// Aggregate Result
// +quiet
func separatorAggregateResultForQuiet(results []*Result, separator string) string {
	aggResult := newAggregateQuietResult(results)
	return fmt.Sprintf("%v", aggResult.StatusMessage())
}

// Aggregate Result
// +pretty
// +quiet
func separatorAggregateResultForQuietPretty(results []*Result, separator string) string {
	aggResult := newAggregateQuietResult(results)

	statusColor := colorGreen
	if !aggResult.Success {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v", statusColor, aggResult.StatusMessage(), colorReset)
}

// Aggregate Result
// +markdown
func separatorAggregateResultForMarkdown(results []*Result, separator string) string {
	aggResult := newAggregateResult(results)

	return fmt.Sprintf("*%v*%v%v", aggResult.StatusMessage(), separator, aggResult.Count)
}

// Aggregate Result
// +markdown
// +quiet
func separatorAggregateResultForQuietMarkdown(results []*Result, separator string) string {
	aggResult := newAggregateQuietResult(results)
	return fmt.Sprintf("*%v*", aggResult.StatusMessage())
}
