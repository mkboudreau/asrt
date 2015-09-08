package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/output"
)

var (
	ErrInvalidTargets       error = errors.New("Must specify at least one target configurer that will return valid targets, such as a file or cli arguments.")
	ErrInvalidTimeoutFormat       = errors.New("Timeout must conform to time.Duration format.")
	ErrInvalidRateFormat          = errors.New("Rate must conform to time.Duration format.")
	ErrInvalidMethod              = errors.New(fmt.Sprintf("Method unknown. Valid methods are %v", ValidMethods))
	ErrInvalidFormat              = errors.New(fmt.Sprintf("Output format unknown. Valid formats are %v", ValidFormats))
)

var ValidFormats = []string{"CSV", "TAB", "JSON"}
var ValidMethods = []string{"GET", "PUT", "POST", "DELETE", "HEAD", "PATCH"}

var DefaultHttpStatuses = map[string]int{
	string(MethodGet):    200,
	string(MethodPost):   201,
	string(MethodPut):    200,
	string(MethodDelete): 200,
	string(MethodHead):   200,
	string(MethodPatch):  200,
}

const (
	FormatCSV  OutputFormat = "CSV"
	FormatTAB               = "TAB"
	FormatJSON              = "JSON"
)

type Configuration struct {
	context         *cli.Context
	CommandName     string
	Rate            time.Duration
	Output          OutputFormat
	Pretty          bool
	AggregateOutput bool
	Quiet           bool
	Quieter         bool
	Markdown        bool
	FailuresOnly    bool
	Workers         int
	Targets         []*Target
}

type OutputFormat string

func GetConfiguration(c *cli.Context) (*Configuration, error) {
	config := &Configuration{context: c, CommandName: c.Command.Name, Targets: make([]*Target, 0)}

	config.Pretty = c.Bool("pretty")
	config.AggregateOutput = c.Bool("aggregate")
	config.Quiet = c.Bool("quiet")
	config.Quieter = c.Bool("quieter")
	config.Markdown = c.Bool("markdown")
	config.Workers = c.Int("workers")
	config.FailuresOnly = c.Bool("failures-only")

	config.Output = OutputFormat(GetUpperOrDefault(c.String("format"), FormatTAB))
	if !validateOutput(string(config.Output)) {
		return nil, ErrInvalidFormat
	}

	config.Rate = GetTimeDurationConfig(c, "rate")

	newTargets, err := getRegisteredConfigureredTargets(c)
	if err != nil {
		return nil, err
	} else if len(newTargets) == 0 {
		return nil, ErrInvalidTargets
	}

	config.Targets = newTargets

	return config, nil
}

func validateOutput(output string) bool {
	switch {
	case output == string(FormatJSON):
		return true
	case output == string(FormatCSV):
		return true
	case output == string(FormatTAB):
		return true
	case output == "":
		return true
	default:
		return false
	}
}

func (config *Configuration) ResultFormatter() output.ResultFormatter {
	modifiers := &output.ResultFormatModifiers{
		Pretty:    config.Pretty,
		Aggregate: config.AggregateOutput,
		Quiet:     config.Quiet,
		Markdown:  config.Markdown,
	}

	switch {
	case config.Output == FormatJSON:
		return output.NewJsonResultFormatter(modifiers)
	case config.Output == FormatCSV:
		return output.NewCsvResultFormatter(modifiers)
	case config.Output == FormatTAB:
		return output.NewTabResultFormatter(modifiers)
	}

	return nil
}

func (config *Configuration) Writer() io.Writer {
	writers := getRegisteredWriters(config.context)
	return newProxyWriteCloser(writers...)
}

func (config *Configuration) WriterWithWriters(writers ...io.Writer) io.Writer {
	configWriters := getRegisteredWriters(config.context)
	configWriters = append(configWriters, writers...)
	return newProxyWriteCloser(configWriters...)
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
	typeOfStdout := reflect.TypeOf(os.Stdout).String()
	for _, w := range pwc.writers {
		if reflect.TypeOf(w).String() == typeOfStdout {
			continue
		}
		if c, ok := w.(io.Closer); ok {
			if err := c.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func getRegisteredWriters(c *cli.Context) []io.Writer {
	var writers []io.Writer

	configurers := GetWriterConfigurers()

	for _, wc := range configurers {
		w := wc.GetWriter(c)
		if w != nil {
			writers = append(writers, w)
		}
	}

	return writers
}

func getRegisteredConfigureredTargets(c *cli.Context) ([]*Target, error) {
	var targets []*Target

	configurers := GetTargetConfigurers()
	for _, tc := range configurers {
		newTargets, err := tc.GetTargets(c)
		if err != nil {
			return nil, err
		}
		targets = append(targets, newTargets...)
	}

	return targets, nil
}
