package output

import (
	"fmt"
	"io"
	"os"
)

var (
	codeClear       string = "\033[2J"
	codeResetCursor        = "\033[H"
)

func WriteToWriter(w io.Writer, r io.Reader) {
	if _, err := io.Copy(w, r); err != nil {
		fmt.Sprintf("error writing to writer: %v", err)
	}
}

func DoneWithWriter(w io.Writer) {
	if c, ok := w.(io.Closer); ok {
		if err := c.Close(); err != nil {
			fmt.Printf("error closing to write closer: %v", err)
		}
	}
}

func ClearConsole() {
	fmt.Printf("%v%v", codeClear, codeResetCursor)
}

func WriteToConsole(r io.Reader) {
	WriteToWriter(os.Stdout, r)
}
