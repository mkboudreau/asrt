package output

import (
	"fmt"
	"io"
	"strings"
)

func getResultStringWithSeparator(result *Result, modifiers *ResultFormatModifiers, separator string) io.Reader {
	var s string

	switch {
	case modifiers.Markdown && !modifiers.Quiet:
		s = separatorResultForMarkdown(result, separator)
	case modifiers.Markdown && modifiers.Quiet:
		s = separatorResultForQuietMarkdown(result, separator)
	case !modifiers.Pretty && !modifiers.Quiet:
		s = separatorResult(result, separator)
	case modifiers.Pretty && !modifiers.Quiet:
		s = separatorResultForPretty(result, separator)
	case !modifiers.Pretty && modifiers.Quiet:
		s = separatorResultForQuiet(result, separator)
	case modifiers.Pretty && modifiers.Quiet:
		s = separatorResultForQuietPretty(result, separator)
	}

	return strings.NewReader(s)
}

func getResultStringAggregateWithSeparator(results []*Result, modifiers *ResultFormatModifiers, separator string) io.Reader {
	var s string

	switch {
	case modifiers.Markdown && !modifiers.Quiet:
		s = separatorAggregateResultForMarkdown(results, separator)
	case modifiers.Markdown && modifiers.Quiet:
		s = separatorAggregateResultForQuietMarkdown(results, separator)
	case !modifiers.Pretty && !modifiers.Quiet:
		s = separatorAggregateResult(results, separator)
	case modifiers.Pretty && !modifiers.Quiet:
		s = separatorAggregateResultForPretty(results, separator)
	case !modifiers.Pretty && modifiers.Quiet:
		s = separatorAggregateResultForQuiet(results, separator)
	case modifiers.Pretty && modifiers.Quiet:
		s = separatorAggregateResultForQuietPretty(results, separator)
	}

	return strings.NewReader(s)
}

func separatorResultForQuietPretty(result *Result, separator string) string {
	statusColor := colorGreen
	if !result.Success || result.Error != nil {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v%v%v", statusColor, result.StatusMessage(), colorReset, separator, result.Url)
}

func separatorResult(result *Result, separator string) string {
	return fmt.Sprintf("%v%v%v%v%v", result.StatusMessage(), separator, result.Expected, separator, result.Url)
}

func separatorResultForPretty(result *Result, separator string) string {
	statusColor := colorGreen
	if !result.Success || result.Error != nil {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v%v%v%v%v", statusColor, result.StatusMessage(), colorReset, separator, result.Expected, separator, result.Url)
}

func separatorAggregateResult(results []*Result, separator string) string {
	aggResult := newAggregateResult(results)
	return fmt.Sprintf("%v%v%v", aggResult.StatusMessage(), separator, aggResult.Count)
}

func separatorAggregateResultForPretty(results []*Result, separator string) string {
	aggResult := newAggregateResult(results)

	statusColor := colorGreen
	if !aggResult.Success {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v%v%v", statusColor, aggResult.StatusMessage(), colorReset, separator, aggResult.Count)
}

func separatorResultForQuiet(result *Result, separator string) string {
	return fmt.Sprintf("%v%v%v", result.StatusMessage(), separator, result.Url)
}

func separatorAggregateResultForQuiet(results []*Result, separator string) string {
	aggResult := newAggregateQuietResult(results)
	return fmt.Sprintf("%v", aggResult.StatusMessage())
}

func separatorAggregateResultForQuietPretty(results []*Result, separator string) string {
	aggResult := newAggregateQuietResult(results)

	statusColor := colorGreen
	if !aggResult.Success {
		statusColor = colorRed
	}

	return fmt.Sprintf("%v%v%v", statusColor, aggResult.StatusMessage(), colorReset)
}

func separatorResultForMarkdown(result *Result, separator string) string {
	return fmt.Sprintf("*%v*%v%v%v%v", result.StatusMessage(), separator, result.Expected, separator, result.Url)
}

func separatorResultForQuietMarkdown(result *Result, separator string) string {
	return fmt.Sprintf("*%v*%v%v", result.StatusMessage(), separator, result.Url)
}

func separatorAggregateResultForMarkdown(results []*Result, separator string) string {
	aggResult := newAggregateResult(results)

	return fmt.Sprintf("*%v*%v%v", aggResult.StatusMessage(), separator, aggResult.Count)
}

func separatorAggregateResultForQuietMarkdown(results []*Result, separator string) string {
	aggResult := newAggregateQuietResult(results)
	return fmt.Sprintf("*%v*", aggResult.StatusMessage())
}
