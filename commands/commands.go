package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var validFormats = []string{"CSV", "TAB", "JSON"}
var validMethods = []string{"GET", "PUT", "POST", "DELETE", "HEAD", "PATCH"}

var statusFlags = []cli.Flag{
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
	cli.StringFlag{
		Name:  "format, fmt",
		Usage: fmt.Sprint("Output format. Valid values:", validFormats),
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
	cli.IntFlag{
		Name:  "workers, w",
		Usage: "Number of workers/goroutines to use to hit the sites",
		Value: 1,
	},
	cli.StringFlag{
		Name:  "method, m",
		Usage: fmt.Sprint("Use HTTP Method for all URLs on command line. Does not affect file inputs. Valid values:", validMethods),
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

var dashboardFlags = append(statusFlags,
	cli.StringFlag{
		Name:  "rate, r",
		Usage: "Rate between refreshes of statuses. Only effective for dashboard settings. 0 = no refresh. Format is Golang time.Duration.",
		Value: "30s",
	})

var Commands = []cli.Command{
	{
		Name:        "status",
		Usage:       "Print simple status lines for the API list",
		Description: "Argument is one ore more URLs if a file is not provided.",
		Action:      cmdStatus,
		Flags:       statusFlags,
	},
	{
		Name:        "dashboard",
		Usage:       "Print a dashboard that refreshes for the API list",
		Description: "Argument is one ore more URLs if a file is not provided.",
		Action:      cmdDashboard,
		Flags:       dashboardFlags,
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
