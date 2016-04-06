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

var ValidFormats = []string{"csv", "csv-md", "csv-no-color", "tab", "tab-md", "tab-no-color", "json", "json-compact", "template={{...}}"}
var ValidMethods = []string{"GET", "PUT", "POST", "DELETE", "HEAD", "PATCH"}

var DefaultHttpStatuses = map[string]int{
	string(MethodGet):    200,
	string(MethodPost):   201,
	string(MethodPut):    200,
	string(MethodDelete): 200,
	string(MethodHead):   200,
	string(MethodPatch):  200,
}

type OutputFormat string

const (
	FormatCSV      OutputFormat = "CSV"
	FormatTAB                   = "TAB"
	FormatJSON                  = "JSON"
	FormatTEMPLATE              = "TEMPLATE"
)

const FormatOptionTemplatePrefix string = "TEMPLATE="

type Configuration struct {
	context         *cli.Context
	CommandName     string
	Rate            time.Duration
	FormatString    string
	Output          OutputFormat
	Pretty          bool
	AggregateOutput bool
	Quiet           bool
	NoHeader        bool
	Markdown        bool
	FailuresOnly    bool
	StateChangeOnly bool
	Workers         int
	Targets         []*Target
}

func GetConfiguration(c *cli.Context) (*Configuration, error) {
	config := &Configuration{context: c, CommandName: c.Command.Name, Targets: make([]*Target, 0)}

	config.AggregateOutput = c.Bool("aggregate")
	config.NoHeader = c.Bool("no-header")
	config.Quiet = c.Bool("quiet")
	config.Workers = c.Int("workers")
	config.FailuresOnly = c.Bool("failures-only")
	config.StateChangeOnly = c.Bool("state-change-only")

	config.FormatString = c.String("format")
	config.Output = GetOutputFormatOrDefault(config.FormatString, FormatTAB)
	config.Pretty = GetPrettyOptionOrDefault(config.FormatString, true)
	config.Markdown = GetMarkdownOptionOrDefault(config.FormatString, false)

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

func (config *Configuration) ResultFormatter() output.ResultFormatter {
	modifiers := &output.ResultFormatModifiers{
		Pretty:    config.Pretty,
		Aggregate: config.AggregateOutput,
		NoHeader:  config.NoHeader,
		Markdown:  config.Markdown,
	}

	switch {
	case config.Output == FormatJSON:
		return output.NewJsonResultFormatter(modifiers)
	case config.Output == FormatCSV:
		return output.NewCsvResultFormatter(modifiers)
	case config.Output == FormatTAB:
		return output.NewTabResultFormatter(modifiers)
	case config.Output == FormatTEMPLATE:
		return output.NewTemplateResultFormatter(config.FormatString, modifiers)
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
