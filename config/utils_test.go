package config

import (
	"testing"
)

var outputFormatTestCases = []struct {
	Text                          string
	DefaultFormat, ExpectedFormat OutputFormat
}{
	{"csv", FormatJSON, FormatCSV},
	{"csv-", FormatJSON, FormatCSV},
	{"csv-md", FormatJSON, FormatCSV},
	{"csv-no-color", FormatJSON, FormatCSV},
	{"csv-compact", FormatJSON, FormatCSV},
	{"CSV-compact", FormatJSON, FormatCSV},
	{"cSv-compact", FormatJSON, FormatCSV},
	{"cSsv-compact", FormatJSON, FormatJSON},

	{"tab", FormatJSON, FormatTAB},
	{"tab-", FormatJSON, FormatTAB},
	{"tab-md", FormatJSON, FormatTAB},
	{"tab-compact", FormatJSON, FormatTAB},
	{"TAB", FormatJSON, FormatTAB},
	{"taB", FormatJSON, FormatTAB},
	{"taab", FormatJSON, FormatJSON},

	{"json", FormatTAB, FormatJSON},
	{"json-", FormatTAB, FormatJSON},
	{"json-md", FormatTAB, FormatJSON},
	{"json-compact", FormatTAB, FormatJSON},
	{"json-no-color", FormatTAB, FormatJSON},
	{"JSON", FormatTAB, FormatJSON},
	{"JSON-md", FormatTAB, FormatJSON},
	{"JSSON-md", FormatTAB, FormatTAB},

	{"template", FormatTAB, FormatTEMPLATE},
	{"template-", FormatTAB, FormatTEMPLATE},
	{"template-md", FormatTAB, FormatTEMPLATE},
	{"template-compact", FormatTAB, FormatTEMPLATE},
	{"TeMpLate", FormatTAB, FormatTEMPLATE},
	{"TEMPLATE", FormatTAB, FormatTEMPLATE},
	{"teemplate", FormatTAB, FormatTAB},

	{"", FormatJSON, FormatJSON},
	{"abc", FormatJSON, FormatJSON},
}

func TestGetOutputFormatOrDefault(t *testing.T) {
	for _, test := range outputFormatTestCases {
		actual := GetOutputFormatOrDefault(test.Text, test.DefaultFormat)
		if actual != test.ExpectedFormat {
			t.Errorf("Expected %v, but got %v with text %v", test.ExpectedFormat, actual, test.Text)
		}
	}
}

var markdownOptionTestCases = []struct {
	Text              string
	Default, Expected bool
}{
	{"csv", true, true},
	{"csv", false, false},
	{"csv-", true, true},
	{"csv-", false, false},
	{"csv-m", true, true},
	{"csv-m", false, false},
	{"tab", true, true},
	{"tab", false, false},
	{"tab-m", true, true},
	{"tab-m", false, false},
	{"tab-pretty", true, true},
	{"tab-pretty", false, false},

	{"csv-md", false, true},
	{"csv-md", true, true},
	{"csv-MD", false, true},
	{"csv-MD", true, true},
	{"CSV-MD", false, true},
	{"CSV-MD", true, true},
	{"tAB-MD", false, true},
	{"tab-MD", false, true},
	{"tab-MD", true, true},
	{"json-md", false, true},
	{"json-md", true, true},
}

func TestGetMarkdownOptionOrDefault(t *testing.T) {
	for _, test := range markdownOptionTestCases {
		actual := GetMarkdownOptionOrDefault(test.Text, test.Default)
		if actual != test.Expected {
			t.Errorf("Expected %v, but got %v with text %v", test.Expected, actual, test.Text)
		}
	}
}

var prettyOptionTestCases = []struct {
	Text              string
	Default, Expected bool
}{
	{"csv", true, true},
	{"csv", false, false},
	{"csv-", true, true},
	{"csv-", false, false},
	{"csv-m", true, true},
	{"csv-m", false, false},
	{"tab", true, true},
	{"tab", false, false},
	{"tab-m", true, true},
	{"tab-m", false, false},

	{"csv-no-color", false, false},
	{"csv-no-color", true, false},
	{"csv-compact", true, false},
	{"csv-compact", true, false},
	{"csv-COMPACT", true, false},
	{"csv-COMPACT", true, false},
	{"csv-md", true, false},
	{"csv-MD", false, false},
	{"csv-MD", true, false},
	{"CSV-MD", false, false},
	{"CSV-MD", true, false},
	{"tAB-MD", false, false},
	{"tab-MD", false, false},
	{"tab-MD", true, false},
	{"json-md", false, false},
	{"json-md", true, false},
}

func TestGetPrettyOptionOrDefault(t *testing.T) {
	for _, test := range prettyOptionTestCases {
		actual := GetPrettyOptionOrDefault(test.Text, test.Default)
		if actual != test.Expected {
			t.Errorf("Expected %v, but got %v with text %v", test.Expected, actual, test.Text)
		}
	}
}
