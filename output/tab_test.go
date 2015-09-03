package output

import (
	"errors"
	"fmt"
	"testing"
)

var tabTestSet = []*resultFormatTestCase{
	{
		expect:    "[ok]\t200\t200\twww",
		results:   []*Result{NewResult(true, nil, "200", "200", "www")},
		modifiers: &ResultFormatModifiers{},
		format:    tabFormat,
	},
	{
		expect:    "[!ok]\t404\t406\twww.notfound.com",
		results:   []*Result{NewResult(false, nil, "404", "406", "www.notfound.com")},
		modifiers: &ResultFormatModifiers{},
		format:    tabFormat,
	},
	{
		expect:    "[err]\t500\tn/a\twww",
		results:   []*Result{NewResult(false, errors.New("hi"), "500", "", "www")},
		modifiers: &ResultFormatModifiers{},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v\t500\t500\tabc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "500", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v\t500\t406\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "500", "406", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v\t500\tn/a\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    "*[ok]*\t201\t201\twww",
		results:   []*Result{NewResult(true, nil, "201", "201", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    "*[err]*\t201\tn/a\twww",
		results:   []*Result{NewResult(false, errors.New(""), "201", "", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    "*[!ok]*\t201\t406\twww",
		results:   []*Result{NewResult(false, nil, "201", "406", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    "[ok]\twww",
		results:   []*Result{NewResult(true, nil, "200", "200", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    "[!ok]\twww",
		results:   []*Result{NewResult(false, nil, "200", "406", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    "[err]\twww",
		results:   []*Result{NewResult(false, errors.New("hi"), "200", "", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    "*[!ok]*\tabc",
		results:   []*Result{NewResult(false, nil, "500", "406", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v\tabc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "500", "406", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*\t500\t500\tabc"),
		results:   []*Result{NewResult(true, nil, "500", "500", "abc")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*\t202\t406\txyz"),
		results:   []*Result{NewResult(false, nil, "202", "406", "xyz")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*\t300\t406\tabcdef"),
		results:   []*Result{NewResult(false, nil, "300", "406", "abcdef")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*\t301\tn/a\tabcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "301", "", "abcd")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*\tabcdxxx"),
		results:   []*Result{NewResult(true, nil, "334", "334", "abcdxxx")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*\tzzz"),
		results:   []*Result{NewResult(false, nil, "333", "406", "zzz")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*\tabcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "335", "", "abcd")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    tabFormat,
	},
}

func TestTabResultFormats(t *testing.T) {
	for _, testcase := range tabTestSet {
		runResultFormatTestCase(t, testcase)
	}
}
