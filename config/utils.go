package config

import (
	"log"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

func GetUpperOrDefault(val string, def interface{}) string {
	v := strings.ToUpper(val)
	if v == "" {
		if tmp, ok := def.(string); !ok {
			return ""
		} else {
			v = tmp
		}
	}
	return v
}

func GetTimeDurationConfig(c *cli.Context, key string) time.Duration {
	v := c.String(key)
	d, err := time.ParseDuration(v)
	if err != nil {
		d = 0 * time.Second
		log.Printf("Could not parse duration %v, defaulting to %v", v, d)
	}
	return d
}
