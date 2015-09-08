package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
		line = strings.Replace(line, "\t", " ", 10)
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, " ")
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			return nil, fmt.Errorf("Found invalid line in file [%v]. Must have method and url", line)
		}

		method := config.CommandMethod(config.GetUpperOrDefault(parts[0], config.MethodGet))
		if method == "" {
			return nil, config.ErrInvalidMethod
		}

		urlString := parts[1]
		var t *config.Target
		var tErr error

		if len(parts) > 2 {
			expectedStatus, sError := strconv.Atoi(parts[2])
			if sError != nil && sError != io.EOF {
				return nil, fmt.Errorf("Found invalid line in file [%v]. Could not parse expected status %v: error %v", line, parts[2], sError)
			}
			t, tErr = config.NewTarget("", urlString, expectedStatus)
		} else {

			expectedStatus := config.DefaultHttpStatuses[string(method)]
			t, tErr = config.NewTarget("", urlString, expectedStatus)
		}
		if tErr != nil {
			return nil, fmt.Errorf("could not create target with url %v: %v", urlString, tErr)
		}

		t.Timeout = timeout
		t.Method = method
		targets = append(targets, t)
	}

	return targets, nil
}
