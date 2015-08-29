package output

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func NewTimeReader(t time.Time) io.Reader {
	timeString := fmt.Sprintf("%v\n", t.Format(time.RFC1123Z))
	return strings.NewReader(timeString)
}

func NewTimeString(t time.Time) string {
	return fmt.Sprintf("%v", t.Format(time.RFC1123Z))
}

func NewTimeStringForJSON(t time.Time) string {
	return fmt.Sprintf("%v", t.Format(time.RFC3339))
}

func NewPrettyTimeReader(t time.Time) io.Reader {
	timeString := fmt.Sprintf("%v%v%v\n", colorYellow, t.Format(time.RFC1123), colorReset)
	return strings.NewReader(timeString)
}
