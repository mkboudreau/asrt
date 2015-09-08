package args

import (
	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
)

func init() {
	config.RegisterTargetConfigurer(new(targetFromArgs))
}

type targetFromArgs struct {
}

func (tc *targetFromArgs) String() string {
	return "CLI Args Target Configurer"
}

func (tc *targetFromArgs) GetCommandFlags() []cli.Flag {
	return []cli.Flag{}
}

func (tc *targetFromArgs) GetTargets(c *cli.Context) ([]*config.Target, error) {
	var targets []*config.Target

	timeout := config.GetTimeDurationConfig(c, "timeout")
	for _, entry := range c.Args() {
		t, err := config.ParseTarget(entry)
		if err != nil {
			return nil, err
		}
		t.Timeout = timeout
		targets = append(targets, t)
	}
	return targets, nil
}
