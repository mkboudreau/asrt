package output

import (
	"bytes"
	"io"
	"log"
	"strings"
	"text/template"
)

/*
type Result struct {
	Success   bool                   `json:"ok"`
	Error     error                  `json:"error,omitempty"`
	Expected  string                 `json:"expectation,omitempty"`
	Actual    string                 `json:"actual,omitempty"`
	Url       string                 `json:"url,omitempty"`
	Label     string                 `json:"label,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}
*/

type TemplateResultFormatter struct {
	OutputTemplate *template.Template
	Modifiers      *ResultFormatModifiers
}

func NewTemplateResultFormatter(optionString string, m *ResultFormatModifiers) *TemplateResultFormatter {
	var tmpl *template.Template
	if isOptionStringFile(optionString) {
		tmplFile := extractTemplateFile(optionString)
		if t, err := parseTemplateFile(tmplFile); err != nil {
			log.Printf("Could not extract template file %v from option %v", tmplFile, optionString)
		} else {
			tmpl = t
		}
	} else {
		tmplStr := extractTemplateString(optionString)
		if t, err := parseTemplateString(tmplStr); err != nil {
			log.Printf("Could not extract template %v from option %v", tmplStr, optionString)
		} else {
			tmpl = t
		}
	}
	return &TemplateResultFormatter{OutputTemplate: tmpl, Modifiers: m}
}

func extractTemplateFile(optionString string) string {
	return extractParts(optionString)
}

func extractTemplateString(optionString string) string {
	if !isOptionStringValid(optionString) {
		return ""
	}
	return extractParts(optionString)
}

func extractParts(optionString string) string {
	parts := strings.SplitAfterN(optionString, "=", 2)
	if len(parts) != 2 {
		return ""
	} else {
		return parts[1]
	}
}

func isOptionStringFile(s string) bool {
	return strings.HasPrefix(strings.ToUpper(s), "TEMPLATE-FILE=")
}

func isOptionStringValid(s string) bool {
	return strings.HasPrefix(strings.ToUpper(s), "TEMPLATE=")
}

func parseTemplateString(templateString string) (*template.Template, error) {
	return template.New("TemplateResultFormatter-Output").Parse(templateString)
}

func parseTemplateFile(templateFile string) (*template.Template, error) {
	return template.ParseFiles(templateFile)
}

func (rf *TemplateResultFormatter) Header() io.Reader {
	return strings.NewReader("")
}

func (rf *TemplateResultFormatter) Footer() io.Reader {
	return strings.NewReader("")
}

func (rf *TemplateResultFormatter) RecordSeparator() io.Reader {
	return strings.NewReader("\n")
}

func (rf *TemplateResultFormatter) Reader(result *Result) io.Reader {
	if rf.OutputTemplate == nil {
		return strings.NewReader("")
	}
	buffer := new(bytes.Buffer)
	rf.OutputTemplate.Execute(buffer, result)
	return buffer
}

func (rf *TemplateResultFormatter) AggregateReader(results []*Result) io.Reader {
	if rf.OutputTemplate == nil {
		return strings.NewReader("")
	}
	buffer := new(bytes.Buffer)
	rf.OutputTemplate.Execute(buffer, newAggregateQuietResult(results))
	return buffer
}
