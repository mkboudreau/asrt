package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
)

var commonFlags = []cli.Flag{
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
}

var statusFlags = append(commonFlags, getRegisteredConfigurers()...)

var dashboardFlags = append(statusFlags,
	cli.StringFlag{
		Name:  "rate, r",
		Usage: "Rate between refreshes of statuses. Only effective for dashboard settings. 0 = no refresh. Format is Golang time.Duration.",
		Value: "30s",
	})

var serverFlags = append(dashboardFlags,
	cli.StringFlag{
		Name:  "port",
		Usage: "Port to listen on",
		Value: "7070",
	})

func getRegisteredConfigurers() []cli.Flag {
	var flags []cli.Flag
	for _, c := range config.GetAllConfigurers() {
		flags = append(flags, c.GetCommandFlags()...)
	}
	return flags
}

var commands = []cli.Command{
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
	{
		Name:        "server",
		Usage:       "Listen on a port for requests",
		Description: "Argument is one ore more URLs if a file is not provided.",
		Action:      cmdServer,
		Flags:       serverFlags,
	},
}

func GetCommands() []cli.Command {
	return commands
}

func GetCommand(commandName string) cli.Command {
	cmdarr := GetCommands()
	for _, cmd := range cmdarr {
		if cmd.Name == commandName {
			return cmd
		}
	}
	return cli.Command{}
}

func GetFlagsForCommand(commandName string) []cli.Flag {
	cmd := GetCommand(commandName)
	return cmd.Flags
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
