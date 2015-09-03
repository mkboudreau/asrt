package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
)

var StatusFlags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "DEBUG",
		Name:   "debug, d",
		Usage:  "Enable debug logging mode",
	},
	cli.StringFlag{
		Name:  "timeout, t",
		Usage: "Timeout for client to wait for connections. 0 = no timeout. Format is Golang time.Duration.",
		Value: "0",
	},
	cli.StringFlag{
		Name:  "file, f",
		Usage: "Use file with list of URLs, HTTP Methods, and optional HTTP Headers",
		Value: "",
	},
	cli.BoolFlag{
		Name:  "failures-only",
		Usage: "Only report on failures. This is useful for long running jobs, especially when coupled with notification to slack or something similar.",
	},
	cli.StringFlag{
		Name:  "format, fmt",
		Usage: fmt.Sprint("Output format. Valid values:", config.ValidFormats),
		Value: "tab",
	},
	cli.BoolFlag{
		Name:  "pretty, p",
		Usage: "Pretty Print the Output. Note: this is mutually exclusive with the pretty option. markdown always take precedence over pretty",
	},
	cli.BoolFlag{
		Name:  "markdown, md",
		Usage: "Markdown Output. Note: this is mutually exclusive with the pretty option. markdown always take precedence over pretty",
	},
	cli.BoolFlag{
		Name:  "aggregate, a",
		Usage: "Aggregate all results into a single result",
	},
	cli.BoolFlag{
		Name:  "quiet, q",
		Usage: "Quiet results into just the statuses. Usually useful in aggregate -qa",
	},
	cli.BoolFlag{
		Name:  "quieter, qq",
		Usage: "Turns off standard output",
	},
	cli.IntFlag{
		Name:  "workers, w",
		Usage: "Number of workers/goroutines to use to hit the sites",
		Value: 1,
	},
	cli.StringFlag{
		Name:  "method, m",
		Usage: fmt.Sprint("Use HTTP Method for all URLs on command line. Does not affect file inputs. Valid values:", config.ValidMethods),
		Value: "GET",
	},
	cli.StringFlag{
		Name:  "slack-url",
		Usage: "Slack incoming webhook URL. Setting this parameter enables slack notifications",
	},
	cli.StringFlag{
		Name:  "slack-channel",
		Usage: "Overrides the default channel for slack notifications. slack-url is required to enable slack integration.",
	},
	cli.StringFlag{
		Name:  "slack-user",
		Usage: "Overrides the username this application posts as for slack notifications. slack-url is required to enable slack integration.",
	},
	cli.StringFlag{
		Name:  "slack-icon",
		Usage: "Overrides the icon used for this application for slack notifications. slack-url is required to enable slack integration.",
	},
}

var DashboardFlags = append(StatusFlags,
	cli.StringFlag{
		Name:  "rate, r",
		Usage: "Rate between refreshes of statuses. Only effective for dashboard settings. 0 = no refresh. Format is Golang time.Duration.",
		Value: "30s",
	})

var ServerFlags = append(DashboardFlags,
	cli.StringFlag{
		Name:  "port",
		Usage: "Port to listen on",
		Value: "7070",
	})

var Commands = []cli.Command{
	{
		Name:        "status",
		Usage:       "Print simple status lines for the API list",
		Description: "Argument is one ore more URLs if a file is not provided.",
		Action:      cmdStatus,
		Flags:       StatusFlags,
	},
	{
		Name:        "dashboard",
		Usage:       "Print a dashboard that refreshes for the API list",
		Description: "Argument is one ore more URLs if a file is not provided.",
		Action:      cmdDashboard,
		Flags:       DashboardFlags,
	},
	{
		Name:        "server",
		Usage:       "Listen on a port for requests",
		Description: "Argument is one ore more URLs if a file is not provided.",
		Action:      cmdServer,
		Flags:       ServerFlags,
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Printf(
		"%s: '%s' is not an %s command. See '%s --help'.\n",
		c.App.Name,
		command,
		c.App.Name,
		c.App.Name,
	)
}
