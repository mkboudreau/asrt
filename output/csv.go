package output

import (
	"fmt"
	"io"
	"strings"
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
				return strings.NewReader("*RESULT*,*EXPECT*,*ACTUAL*,*LABEL*,*URL*\n")
			}
		} else if rf.Modifiers.Pretty {
			if rf.Modifiers.Aggregate {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v,%vCOUNT%v\n", colorYellow, colorReset, colorYellow, colorReset))
			} else {
				return strings.NewReader(fmt.Sprintf("%vRESULT%v,%vEXPECT%v,%vACTUAL%v,%vLABEL%v,%vURL%v\n", colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset, colorYellow, colorReset))
			}
		} else {
			if rf.Modifiers.Aggregate {
				return strings.NewReader("RESULT,COUNT\n")
			} else {
				return strings.NewReader("RESULT,EXPECT,ACTUAL,LABEL,URL\n")
			}
		}
	}
	return strings.NewReader("")
}

func (rf *CsvResultFormatter) Footer() io.Reader {
	if rf.Modifiers.Quiet {
		return strings.NewReader("")
	}
	return strings.NewReader("\n")
}
func (rf *CsvResultFormatter) RecordSeparator() io.Reader {
	return strings.NewReader("\n")
}

func (rf *CsvResultFormatter) Reader(result *Result) io.Reader {
	if result.Label == "" {
		result.Label = "n/a"
	}
	return getResultStringWithSeparator(result, rf.Modifiers, ",")
}

func (rf *CsvResultFormatter) AggregateReader(results []*Result) io.Reader {
	return getResultStringAggregateWithSeparator(results, rf.Modifiers, ",")
}
