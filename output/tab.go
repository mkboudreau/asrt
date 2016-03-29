package output

import (
	"fmt"
	"io"
	"strings"
)

type TabResultFormatter struct {
	Modifiers *ResultFormatModifiers
}

func NewTabResultFormatter(m *ResultFormatModifiers) *TabResultFormatter {
	return &TabResultFormatter{Modifiers: m}
}

func (rf *TabResultFormatter) Header() io.Reader {
	if !rf.Modifiers.NoHeader {
		if rf.Modifiers.Markdown {
			if rf.Modifiers.Aggregate {
				return strings.NewReader("*RESULT*\t*COUNT*\n")
			} else {
				return strings.NewReader("*RESULT*\t*EXPECT*\t*ACTUAL*\t*LABEL*\t\t\t*URL*\n")
			}
		} else if rf.Modifiers.Pretty {
			if rf.Modifiers.Aggregate {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v\t%vCOUNT%v\n", colorYellow, colorReset, colorYellow, colorReset))
			} else {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v\t%vEXPECT%v\t%vACTUAL%v\t%vLABEL%v\t\t\t%vURL%v\n", colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset))
			}
		} else {
			if rf.Modifiers.Aggregate {
				return strings.NewReader("RESULT\tCOUNT\n")
			} else {
				return strings.NewReader("RESULT\tEXPECT\tACTUAL\tLABEL\t\t\tURL\n")
			}
		}
	}
	return strings.NewReader("")
}

func (rf *TabResultFormatter) Footer() io.Reader {
	if rf.Modifiers.NoHeader {
		return strings.NewReader("")
	}
	return strings.NewReader("\n")
}

func (rf *TabResultFormatter) RecordSeparator() io.Reader {
	return strings.NewReader("\n")
}

func (rf *TabResultFormatter) Reader(result *Result) io.Reader {
	if result.Label == "" {
		result.Label = "\t"
	} else if len(result.Label) < 8 {
		result.Label = fmt.Sprintf("%v\t\t", result.Label)
	} else if len(result.Label) < 16 {
		result.Label = fmt.Sprintf("%v\t", result.Label)
	}
	return getResultStringWithSeparator(result, rf.Modifiers, "\t")
}

func (rf *TabResultFormatter) AggregateReader(results []*Result) io.Reader {
	return getResultStringAggregateWithSeparator(results, rf.Modifiers, "\t")
}
