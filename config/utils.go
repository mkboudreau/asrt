package config

import (
	"log"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

func GetOutputFormatOrDefault(val string, defaultFormat OutputFormat) OutputFormat {
	v := strings.ToUpper(val)
	if strings.HasPrefix(v, string(FormatCSV)) {
		return FormatCSV
	} else if strings.HasPrefix(v, string(FormatJSON)) {
		return FormatJSON
	} else if strings.HasPrefix(v, string(FormatTAB)) {
		return FormatTAB
	} else if strings.HasPrefix(v, string(FormatTEMPLATE)) {
		return FormatTEMPLATE
	} else {
		return defaultFormat
	}
}

func GetMarkdownOptionOrDefault(val string, defaultBoolean bool) bool {
	v := strings.ToUpper(val)
	if strings.Contains(v, "-MD") {
		return true
	} else {
		return defaultBoolean
	}
}

func GetPrettyOptionOrDefault(val string, defaultBoolean bool) bool {
	v := strings.ToUpper(val)
	if GetMarkdownOptionOrDefault(val, false) {
		return false //markdown takes precedence over color options
	} else if strings.Contains(v, "-NO-COLOR") {
		return false
	} else if strings.Contains(v, "-COMPACT") {
		return false
	} else {
		return defaultBoolean
	}
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
