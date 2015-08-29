package output

import (
	"errors"
	"fmt"
	"testing"
)

var csvTestSet = []*resultFormatTestCase{
	{
		expect:    "[ok],200,www",
		results:   []*Result{NewResult(true, nil, "200", "www")},
		modifiers: &ResultFormatModifiers{},
		format:    csvFormat,
	},
	{
		expect:    "[!ok],404,www.notfound.com",
		results:   []*Result{NewResult(false, nil, "404", "www.notfound.com")},
		modifiers: &ResultFormatModifiers{},
		format:    csvFormat,
	},
	{
		expect:    "[err],500,www",
		results:   []*Result{NewResult(true, errors.New("hi"), "500", "www")},
		modifiers: &ResultFormatModifiers{},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v,500,abc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v,500,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v,500,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "abc")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    "*[ok]*,201,www",
		results:   []*Result{NewResult(true, nil, "201", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    "*[err]*,201,www",
		results:   []*Result{NewResult(true, errors.New(""), "201", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    "*[!ok]*,201,www",
		results:   []*Result{NewResult(false, nil, "201", "www")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    "[ok],www",
		results:   []*Result{NewResult(true, nil, "200", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    csvFormat,
	},
	{
		expect:    "[!ok],www",
		results:   []*Result{NewResult(false, nil, "200", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    csvFormat,
	},
	{
		expect:    "[err],www",
		results:   []*Result{NewResult(false, errors.New("hi"), "200", "www")},
		modifiers: &ResultFormatModifiers{Quiet: true},
		format:    csvFormat,
	},
	{
		expect:    "*[!ok]*,abc",
		results:   []*Result{NewResult(false, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v,abc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "abc")},
		modifiers: &ResultFormatModifiers{Quiet: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*,500,abc"),
		results:   []*Result{NewResult(true, nil, "500", "abc")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*,202,xyz"),
		results:   []*Result{NewResult(false, nil, "202", "xyz")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*,300,abcdef"),
		results:   []*Result{NewResult(false, nil, "300", "abcdef")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*,301,abcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "301", "abcd")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*,abcdxxx"),
		results:   []*Result{NewResult(true, nil, "334", "abcdxxx")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*,zzz"),
		results:   []*Result{NewResult(false, nil, "333", "zzz")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*,abcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "335", "abcd")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, Quiet: true},
		format:    csvFormat,
	},
}

func TestCsvResultFormats(t *testing.T) {
	for _, testcase := range csvTestSet {
		runResultFormatTestCase(t, testcase)
	}
}
