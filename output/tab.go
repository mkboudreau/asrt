package output

import (
	"fmt"
)

const (
	fmtTabRaw string = "%v\t%v\t%v\n"

	// three sections: status text, expected value, url
	// for each section:
	// structure: color before, data, color after, tab
	fmtTabPretty = "%v%v%v\t%v%v%v\t%v%v%v\n"
)

func WriteToTab(result *Result) {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	fmt.Printf(fmtTabRaw, status, result.Expected, result.Url)
}

func WriteToTabPretty(result *Result) {
	statusColor := colorGreen
	statusText := statusTextOk
	if !result.Success {
		statusColor = colorRed
		statusText = statusTextNotOk
	}
	if result.Error != nil {
		statusColor = colorRed
		statusText = statusTextError
	}

	bStatus := statusColor
	aStatus := colorReset
	bExpected := ""
	aExpected := ""
	bUrl := ""
	aUrl := ""

	fmt.Printf(fmtTabPretty, bStatus, statusText, aStatus, bExpected, result.Expected, aExpected, bUrl, result.Url, aUrl)
}
