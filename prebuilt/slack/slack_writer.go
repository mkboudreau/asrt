package slack

import (
	"io"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/writer"
)

func init() {
	config.RegisterWriterConfigurer(new(slackWriterConfigurer))
}

type slackWriterConfigurer struct {
}

func (wc *slackWriterConfigurer) String() string {
	return "Slack Writer Configurer"
}

func (wc *slackWriterConfigurer) GetCommandFlags() []cli.Flag {
	return []cli.Flag{
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
}

func (wc *slackWriterConfigurer) GetWriter(c *cli.Context) io.Writer {
	if c.String("slack-url") == "" {
		return nil
	}

	w := writer.NewSlackWriter(c.String("slack-url"))

	if c.String("slack-channel") != "" {
		w.SlackChannel(c.String("slack-channel"))
	}
	if c.String("slack-user") != "" {
		w.SlackUser(c.String("slack-user"))
	}
	if c.String("slack-icon") != "" {
		icon := c.String("slack-icon")
		w.SlackIconUrl(writer.SlackIconUrl(icon))
	}

	return w
}
