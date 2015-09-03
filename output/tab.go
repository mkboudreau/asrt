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
	if !rf.Modifiers.Quiet {
		if rf.Modifiers.Markdown {
			if rf.Modifiers.Aggregate {
				return strings.NewReader("*RESULT*\t*COUNT*\n")
			} else {
				return strings.NewReader("*RESULT*\t*EXPECT*\t*ACTUAL*\t*URL*\n")
			}
		} else if rf.Modifiers.Pretty {
			if rf.Modifiers.Aggregate {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v\t%vCOUNT%v\n", colorYellow, colorReset, colorYellow, colorReset))
			} else {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v\t%vEXPECT%v\t%vACTUAL%v\t%vURL%v\n", colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset))
			}
		} else {
			if rf.Modifiers.Aggregate {
				return strings.NewReader("RESULT\tCOUNT\n")
			} else {
				return strings.NewReader("RESULT\tEXPECT\tACTUAL\tURL\n")
			}
		}
	}
	return strings.NewReader("")
}

func (rf *TabResultFormatter) Footer() io.Reader {
	if rf.Modifiers.Quiet {
		return strings.NewReader("")
	}
	return strings.NewReader("\n")
}

func (rf *TabResultFormatter) RecordSeparator() io.Reader {
	return strings.NewReader("\n")
}

func (rf *TabResultFormatter) Reader(result *Result) io.Reader {
	return getResultStringWithSeparator(result, rf.Modifiers, "\t")
}

func (rf *TabResultFormatter) AggregateReader(results []*Result) io.Reader {
	return getResultStringAggregateWithSeparator(results, rf.Modifiers, "\t")
}
