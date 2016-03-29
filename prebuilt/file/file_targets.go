package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

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
		cli.BoolFlag{
			Name:  "no-environment",
			Usage: "Turn off environment substitution in target files.",
		},
	}
}

func (tc *targetFromFile) GetTargets(c *cli.Context) ([]*config.Target, error) {
	filename := c.String("file")
	if filename == "" {
		return []*config.Target{}, nil
	}

	timeout := config.GetTimeDurationConfig(c, "timeout")
	turnOffEnvExpansion := c.Bool("no-environment")

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return tc.targetsFromReaderWithTimeout(file, timeout, turnOffEnvExpansion)
}

func (tc *targetFromFile) targetsFromReaderWithTimeout(reader io.Reader, timeout time.Duration, noEnv bool) ([]*config.Target, error) {
	var targets []*config.Target
	r := bufio.NewReader(reader)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		line = strings.Trim(line, "\n\t ")

		if strings.HasPrefix(line, "#") {
			// ignore
			continue
		} else if len(line) == 0 {
			// ignore
			continue
		}

		if !noEnv {
			line = os.ExpandEnv(line)
		}

		t, tErr := config.ParseTarget(line)
		if tErr != nil {
			return nil, fmt.Errorf("could not create target from string %v: %v", line, tErr)
		}

		t.Timeout = timeout
		targets = append(targets, t)
	}

	return targets, nil
}
