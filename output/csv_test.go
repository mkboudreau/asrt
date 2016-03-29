package output

import (
	"errors"
	"fmt"
	"testing"
)

var csvTestSet = []*resultFormatTestCase{
	{
		expect:    "[ok],200,200,n/a,www",
		results:   []*Result{NewResult(true, nil, "200", "200", "www", "")},
		modifiers: &ResultFormatModifiers{},
		format:    csvFormat,
	},
	{
		expect:    "[!ok],404,400,n/a,www.notfound.com",
		results:   []*Result{NewResult(false, nil, "404", "400", "www.notfound.com", "")},
		modifiers: &ResultFormatModifiers{},
		format:    csvFormat,
	},
	{
		expect:    "[err],500,n/a,n/a,www",
		results:   []*Result{NewResult(false, errors.New("hi"), "500", "", "www", "")},
		modifiers: &ResultFormatModifiers{},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v,500,500,n/a,abc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "500", "abc", "")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v,404,500,n/a,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "404", "500", "abc", "")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v,500,n/a,n/a,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "", "abc", "")},
		modifiers: &ResultFormatModifiers{Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    "*[ok]*,201,201,n/a,www",
		results:   []*Result{NewResult(true, nil, "201", "201", "www", "")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    "*[err]*,201,201,n/a,www",
		results:   []*Result{NewResult(false, errors.New(""), "201", "201", "www", "")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    "*[!ok]*,200,201,n/a,www",
		results:   []*Result{NewResult(false, nil, "200", "201", "www", "")},
		modifiers: &ResultFormatModifiers{Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    "[ok],www",
		results:   []*Result{NewResult(true, nil, "200", "200", "www", "")},
		modifiers: &ResultFormatModifiers{NoHeader: true},
		format:    csvFormat,
	},
	{
		expect:    "[!ok],www",
		results:   []*Result{NewResult(false, nil, "201", "200", "www", "")},
		modifiers: &ResultFormatModifiers{NoHeader: true},
		format:    csvFormat,
	},
	{
		expect:    "[err],www",
		results:   []*Result{NewResult(false, errors.New("hi"), "200", "", "www", "")},
		modifiers: &ResultFormatModifiers{NoHeader: true},
		format:    csvFormat,
	},
	{
		expect:    "*[!ok]*,abc",
		results:   []*Result{NewResult(false, nil, "400", "500", "abc", "")},
		modifiers: &ResultFormatModifiers{NoHeader: true, Markdown: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[ok]%v,abc", colorGreen, colorReset),
		results:   []*Result{NewResult(true, nil, "500", "500", "abc", "")},
		modifiers: &ResultFormatModifiers{NoHeader: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[!ok]%v,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, nil, "400", "500", "abc", "")},
		modifiers: &ResultFormatModifiers{NoHeader: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("%v[err]%v,abc", colorRed, colorReset),
		results:   []*Result{NewResult(false, fmt.Errorf("HELLO"), "500", "", "abc", "")},
		modifiers: &ResultFormatModifiers{NoHeader: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*,500,500,n/a,abc"),
		results:   []*Result{NewResult(true, nil, "500", "500", "abc", "")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*,400,202,n/a,xyz"),
		results:   []*Result{NewResult(false, nil, "400", "202", "xyz", "")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*,400,300,n/a,abcdef"),
		results:   []*Result{NewResult(false, nil, "400", "300", "abcdef", "")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*,301,n/a,n/a,abcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "301", "", "abcd", "")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[ok]*,abcdxxx"),
		results:   []*Result{NewResult(true, nil, "334", "334", "abcdxxx", "")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, NoHeader: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[!ok]*,zzz"),
		results:   []*Result{NewResult(false, nil, "300", "333", "zzz", "")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, NoHeader: true},
		format:    csvFormat,
	},
	{
		expect:    fmt.Sprintf("*[err]*,abcd"),
		results:   []*Result{NewResult(false, fmt.Errorf("ABC"), "335", "", "abcd", "")},
		modifiers: &ResultFormatModifiers{Markdown: true, Pretty: true, NoHeader: true},
		format:    csvFormat,
	},
}

func TestCsvResultFormats(t *testing.T) {
	for _, testcase := range csvTestSet {
		runResultFormatTestCase(t, testcase)
	}
}
