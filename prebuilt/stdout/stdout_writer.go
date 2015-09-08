package stdout

import (
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
)

func init() {
	config.RegisterWriterConfigurer(new(stdoutWriterConfigurer))
}

type stdoutWriterConfigurer struct {
}

func (wc *stdoutWriterConfigurer) String() string {
	return "Stdout Writer Configurer"
}

func (wc *stdoutWriterConfigurer) GetCommandFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "quieter, qq",
			Usage: "Turns off standard output",
		},
	}
}

func (wc *stdoutWriterConfigurer) GetWriter(c *cli.Context) io.Writer {
	if c.Bool("quieter") {
		return nil
	}

	switch c.Command.Name {
	case "status":
		return os.Stdout
	case "dashboard":
		return os.Stdout
	case "server":
		return os.Stdout
	}
	return nil

}
