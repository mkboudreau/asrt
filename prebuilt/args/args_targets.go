package args

import (
	"fmt"

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
	method := config.CommandMethod(config.GetUpperOrDefault(c.String("method"), config.MethodGet))
	if method == "" {
		return nil, config.ErrInvalidMethod
	}

	for _, u := range c.Args() {
		expectedStatus := config.DefaultHttpStatuses[string(method)]
		t, err := config.NewTarget("", u, expectedStatus)
		if err != nil {
			return nil, fmt.Errorf("could not create target with url %v: %v", u, err)
		}
		t.Timeout = timeout
		t.Method = method
		targets = append(targets, t)
	}
	return targets, nil
}
