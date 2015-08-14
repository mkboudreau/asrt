package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/output"
)

var (
	ErrInvalidTargets       error = errors.New("Must specify url targets in file or as arguments to command.")
	ErrInvalidTimeoutFormat       = errors.New("Timeout must conform to time.Duration format.")
	ErrInvalidRateFormat          = errors.New("Rate must conform to time.Duration format.")
	ErrInvalidMethod              = errors.New(fmt.Sprintf("Method unknown. Valid methods are %d", validMethods))
	ErrInvalidFormat              = errors.New(fmt.Sprintf("Output format unknown. Valid formats are %d", validFormats))
)

var DefaultHttpStatuses = map[string]int{
	string(methodGet):    200,
	string(methodPost):   201,
	string(methodPut):    200,
	string(methodDelete): 200,
	string(methodHead):   200,
	string(methodPatch):  200,
}

const (
	methodGet    commandMethod = "GET"
	methodPost                 = "POST"
	methodPut                  = "PUT"
	methodDelete               = "DELETE"
	methodHead                 = "HEAD"
	methodPatch                = "PATCH"
)

const (
	formatCSV  outputFormat = "CSV"
	formatTAB               = "TAB"
	formatJSON              = "JSON"
)

type configuration struct {
	context         *cli.Context
	CommandName     string
	Rate            time.Duration
	Output          outputFormat
	Pretty          bool
	AggregateOutput bool
	Quiet           bool
	Markdown        bool
	Workers         int
	Targets         []*target
}

type target struct {
	Method         commandMethod
	Timeout        time.Duration
	ExpectedStatus int
	URL            string
	Headers        map[string]string
}

type commandMethod string
type outputFormat string

func getConfiguration(c *cli.Context) (*configuration, error) {
	config := &configuration{context: c, CommandName: c.Command.Name, Targets: make([]*target, 0)}

	config.Pretty = c.Bool("pretty")
	config.AggregateOutput = c.Bool("aggregate")
	config.Quiet = c.Bool("quiet")
	config.Markdown = c.Bool("markdown")
	config.Workers = c.Int("workers")
	file := c.String("file")

	config.Output = outputFormat(getUpperOrDefault(c.String("format"), formatTAB))
	if config.Output == "" {
		return nil, ErrInvalidFormat
	}

	config.Rate = getTimeDurationConfig(c, "rate")

	if file == "" && len(c.Args()) == 0 {
		//cli.ShowCommandHelp(c, "status")
		return nil, ErrInvalidTargets
	}

	if newTargets, err := buildTargetsFromArgs(c); err != nil {
		return nil, err
	} else {
		config.Targets = append(config.Targets, newTargets...)
	}
	if file != "" {
		if newTargets, err := buildTargetsFromFile(c); err != nil {
			return nil, err
		} else {
			config.Targets = append(config.Targets, newTargets...)
		}
	}

	return config, nil
}

func getUpperOrDefault(val string, def interface{}) string {
	v := strings.ToUpper(val)
	if v == "" {
		if tmp, ok := def.(string); !ok {
			return ""
		} else {
			v = tmp
		}
	}
	return v
}
func getTimeDurationConfig(c *cli.Context, key string) time.Duration {
	v := c.String(key)
	d, err := time.ParseDuration(v)
	if err != nil {
		d = 0 * time.Second
		log.Printf("Could not parse duration %v, defaulting to %v", v, d)
	}
	return d
}

func (config *configuration) ResultFormatter() output.ResultFormatter {
	modifiers := &output.ResultFormatModifiers{
		Pretty:    config.Pretty,
		Aggregate: config.AggregateOutput,
		Quiet:     config.Quiet,
		Markdown:  config.Markdown,
	}

	switch {
	case config.Output == formatJSON:
		return output.NewJsonResultFormatter(modifiers)
	case config.Output == formatCSV:
		return output.NewCsvResultFormatter(modifiers)
	case config.Output == formatTAB:
		return output.NewTabResultFormatter(modifiers)
	}

	return nil
}

func (config *configuration) Writer() io.Writer {
	writers := make([]io.Writer, 0)

	main := config.getMainWriter()
	if main != nil {
		writers = append(writers, main)
	}

	slack := config.getSlackWriter()
	if slack != nil {
		writers = append(writers, slack)
	}

	return newProxyWriteCloser(writers...)
}

type proxyWriteCloser struct {
	writers []io.Writer
	writer  io.Writer
}

func newProxyWriteCloser(writers ...io.Writer) *proxyWriteCloser {
	return &proxyWriteCloser{
		writers: writers,
		writer:  io.MultiWriter(writers...),
	}
}

func (pwc *proxyWriteCloser) Write(p []byte) (n int, err error) {
	return pwc.writer.Write(p)
}

func (pwc *proxyWriteCloser) Close() error {
	for _, w := range pwc.writers {
		if c, ok := w.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (config *configuration) getMainWriter() io.Writer {
	switch config.context.Command.Name {
	case "status":
		return os.Stdout
	case "dashboard":
		return os.Stdout
	case "server":
	}
	return nil
}

func (config *configuration) getSlackWriter() io.Writer {
	if config.context.String("slack-url") == "" {
		return nil
	}

	w := output.NewSlackWriter(config.context.String("slack-url"))

	if config.context.String("slack-channel") != "" {
		w.SlackChannel(config.context.String("slack-channel"))
	}
	if config.context.String("slack-user") != "" {
		w.SlackUser(config.context.String("slack-user"))
	}
	if config.context.String("slack-icon") != "" {
		icon := config.context.String("slack-icon")
		w.SlackIconUrl(output.SlackIconUrl(icon))
	}

	fmt.Println(w.Payload)

	return w
}

func buildTargetsFromFile(c *cli.Context) ([]*target, error) {
	filename := c.String("file")
	targets := make([]*target, 0)

	timeout := getTimeDurationConfig(c, "timeout")

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

		method := commandMethod(getUpperOrDefault(parts[0], methodGet))
		if method == "" {
			return nil, ErrInvalidMethod
		}

		urlString := parts[1]
		var t *target
		var tErr error

		if len(parts) > 2 {
			expectedStatus, sError := strconv.Atoi(parts[2])
			if sError != nil && sError != io.EOF {
				return nil, fmt.Errorf("Found invalid line in file [%v]. Could not parse expected status %v: error %v", line, parts[2], sError)
			}
			t, tErr = newTargetWithExpectation(urlString, expectedStatus)
		} else {

			expectedStatus := DefaultHttpStatuses[string(method)]
			t, tErr = newTargetWithExpectation(urlString, expectedStatus)
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

func buildTargetsFromArgs(c *cli.Context) ([]*target, error) {

	timeout := getTimeDurationConfig(c, "timeout")
	method := commandMethod(getUpperOrDefault(c.String("method"), methodGet))
	if method == "" {
		return nil, ErrInvalidMethod
	}

	targets := make([]*target, 0)
	for _, u := range c.Args() {
		expectedStatus := DefaultHttpStatuses[string(method)]
		t, err := newTargetWithExpectation(u, expectedStatus)
		if err != nil {
			return nil, fmt.Errorf("could not create target with url %v: %v", u, err)
		}
		t.Timeout = timeout
		t.Method = method
		targets = append(targets, t)
	}
	return targets, nil
}

func newTargetWithExpectation(urlString string, expectedStatus int) (*target, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	t := &target{
		URL:            u.String(),
		ExpectedStatus: expectedStatus,
	}

	return t, nil
}
