package output

import (
	"errors"
	"fmt"
	"testing"
)

var tabTestSet = []*resultFormatTestCase{
	{
		expect:    "[ok]\t200\twww",
		results:   []*Result{NewResult(true, nil, "200", "www")},
		modifiers: &ResultFormatModifiers{},
		format:    tabFormat,
	},
	{
		expect:    "[!ok]\t404\twww.notfound.com",
		results:   []*Result{NewResult(false, nil, "404", "www.notfound.com")},
		modifiers: &ResultFormatModifiers{},
		format:    tabFormat,
	},
	{
		expect:    "[err]\t500\twww",
		results:   []*Result{NewResult(true, errors.New("hi"), "500", "www")},
		modifiers: &ResultFormatModifiers{},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v\t500\tabc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v\t500\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v\t500\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    "*[ok]*\t201\twww",
		results:   []*Result{NewResult(true, nil, "201", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    "*[err]*\t201\twww",
		results:   []*Result{NewResult(true, errors.New(""), "201", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    "*[!ok]*\t201\twww",
		results:   []*Result{NewResult(false, nil, "201", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    "[ok]\twww",
		results:   []*Result{NewResult(true, nil, "200", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    "[!ok]\twww",
		results:   []*Result{NewResult(false, nil, "200", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    "[err]\twww",
		results:   []*Result{NewResult(false, errors.New("hi"), "200", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    "*[!ok]*\tabc",
		results:   []*Result{NewResult(false, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Markdown: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v\tabc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v\tabc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*\t500\tabc"),
		results:   []*Result{NewResult(true, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*\t202\txyz"),
		results:   []*Result{NewResult(false, nil, "202", "xyz")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*\t300\tabcdef"),
		results:   []*Result{NewResult(false, nil, "300", "abcdef")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*\t301\tabcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "301", "abcd")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*\tabcdxxx"),
		results:   []*Result{NewResult(true, nil, "334", "abcdxxx")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*\tzzz"),
		results:   []*Result{NewResult(false, nil, "333", "zzz")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    tabFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*\tabcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "335", "abcd")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    tabFormat,
	},
}

func TestTabResultFormats(t *testing.T) {
	for _, testcase := range tabTestSet {
		runResultFormatTestCase(t, testcase)
	}
}
