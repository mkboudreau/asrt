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
		Usage: "Pretty Print the Output",
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
