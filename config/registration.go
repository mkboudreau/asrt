package config

import (
	"io"

	"github.com/codegangsta/cli"
)

var (
	registeredWriters           = make(map[io.Writer]bool)
	registeredTargetConfigurers = make(map[TargetConfigurer]bool)
)

// TargetConfigurer ...
type TargetConfigurer interface {
	GetTargets(c *cli.Context) []*Target
}

// RegisterTargetConfigurer ...
func RegisterTargetConfigurer(tc TargetConfigurer) {
	if tc != nil {
		registeredTargetConfigurers[tc] = true
	}
}
