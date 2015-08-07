package output

import (
	"fmt"
	"io"
	"strings"
	"time"
)

var (
	fmtTimeRaw    string = "%v\n"
	fmtTimePretty        = "%v%v%v\n"
)

func NewTimeReader(t time.Time) io.Reader {
	timeString := fmt.Sprintf(fmtTimeRaw, t.Format(time.RFC1123Z))
	return strings.NewReader(timeString)
}

func NewPrettyTimeReader(t time.Time) io.Reader {
	timeString := fmt.Sprintf(fmtTimePretty, colorYellow, t.Format(time.RFC1123), colorReset)
	return strings.NewReader(timeString)
}
