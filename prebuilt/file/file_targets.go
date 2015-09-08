package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
)

func init() {
	config.RegisterTargetConfigurer(new(targetFromFile))
}

type targetFromFile struct {
}

func (tc *targetFromFile) String() string {
	return "File Target Configurer"
}

func (tc *targetFromFile) GetCommandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Use file with list of URLs, HTTP Methods, and optional HTTP Headers",
			Value: "",
		},
	}
}

func (tc *targetFromFile) GetTargets(c *cli.Context) ([]*config.Target, error) {
	var targets []*config.Target

	filename := c.String("file")
	if filename == "" {
		return targets, nil
	}

	timeout := config.GetTimeDurationConfig(c, "timeout")

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, " ")

		t, tErr := config.ParseTarget(line)
		if tErr != nil {
			return nil, fmt.Errorf("could not create target from string %v: %v", line, tErr)
		}

		t.Timeout = timeout
		targets = append(targets, t)
	}

	return targets, nil
}
