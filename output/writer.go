package output

import (
	"bufio"
	"fmt"
	"io"
)

var (
	codeClear       string = "\033[2J"
	codeResetCursor        = "\033[H"
)

func WriteToWriter(w io.Writer, r io.Reader) {

}

func ClearConsole() {
	fmt.Printf("%v%v", codeClear, codeResetCursor)
}

func WriteToConsole(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if line != "" {
			fmt.Print(line)
		}
		if err != nil {
			break
		}
	}
}
