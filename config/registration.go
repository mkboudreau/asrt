package config

import (
	"io"

	"github.com/codegangsta/cli"
)

var (
	registeredWriterConfigurers = make(map[WriterConfigurer]bool)
	registeredTargetConfigurers = make(map[TargetConfigurer]bool)
)

// Configurer ...
type Configurer interface {
	GetCommandFlags() []cli.Flag
}

// TargetConfigurer ...
type TargetConfigurer interface {
	Configurer
	GetTargets(c *cli.Context) ([]*Target, error)
}

// RegisterTargetConfigurer ...
func RegisterTargetConfigurer(tc TargetConfigurer) {
	if tc != nil {
		registeredTargetConfigurers[tc] = true
	}
}

// WriterConfigurer ...
type WriterConfigurer interface {
	Configurer
	GetWriter(c *cli.Context) io.Writer
}

// RegisterWriterConfigurer ...
func RegisterWriterConfigurer(wc WriterConfigurer) {
	if wc != nil {
		registeredWriterConfigurers[wc] = true
	}
}

func GetAllConfigurers() []Configurer {
	var configurers []Configurer
	for k, v := range registeredWriterConfigurers {
		if v {
			configurers = append(configurers, k)
		}
	}
	for k, v := range registeredTargetConfigurers {
		if v {
			configurers = append(configurers, k)
		}
	}
	return configurers
}
func GetWriterConfigurers() []WriterConfigurer {
	var configurers []WriterConfigurer
	for k, v := range registeredWriterConfigurers {
		if v {
			configurers = append(configurers, k)
		}
	}
	return configurers
}

func GetTargetConfigurers() []TargetConfigurer {
	var configurers []TargetConfigurer
	for k, v := range registeredTargetConfigurers {
		if v {
			configurers = append(configurers, k)
		}
	}
	return configurers
}
