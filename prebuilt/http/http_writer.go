package http

import (
	"io"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/writer"
)

func init() {
	config.RegisterWriterConfigurer(new(httpWriterConfigurer))
}

type httpWriterConfigurer struct {
}

func (wc *httpWriterConfigurer) String() string {
	return "HTTP Writer Configurer"
}

func (wc *httpWriterConfigurer) GetCommandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "http-url",
			Usage: "http url to send results to. Setting this parameter enables sending result data over http",
		},
		cli.StringFlag{
			Name:  "http-method",
			Usage: "Overrides the default method for http. http-url is required to enable http integration.",
			Value: "POST",
		},
		cli.StringFlag{
			Name:  "http-auth",
			Usage: "If set, and http-url is set, the value here will be passed to the Authorization header. http-url is required to enable http integration.",
		},
	}
}

func (wc *httpWriterConfigurer) GetWriter(c *cli.Context) io.Writer {
	if c.String("http-url") == "" {
		return nil
	}

	w := writer.NewHttpWriter(c.String("http-url"))

	if c.String("http-method") != "" {
		w.HttpMethod(c.String("http-method"))
	}
	if c.String("http-auth") != "" {
		w.HttpAuth(c.String("http-auth"))
	}

	return w
}
