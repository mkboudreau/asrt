package log

import (
	"io"
	"log"
	"os"
)

func init() {
	TurnOffLogging()
	for _, f := range os.Args {
		if f == "-d" || f == "--debug" || f == "-debug" {
			TurnOnLogging()
		}
	}
}

func TurnOffLogging() io.Writer {
	w := &noopWriter{}
	log.SetOutput(w)
	return w
}
func TurnOnLogging() io.Writer {
	w := os.Stdout
	log.SetOutput(w)
	return w
}

type noopWriter struct{}

func (w *noopWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
