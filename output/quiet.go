package output

import (
	"fmt"
)

const (
	fmtQuietRaw string = "%v\n"

	// three sections: status text, expected value, url
	// for each section:
	// structure: color before, data, color after, tab
	fmtQuietPretty = "%v%v%v\n"
)

func WriteToQuiet(result *Result) {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	fmt.Printf(fmtQuietRaw, status)
}

func WriteToQuietPretty(result *Result) {
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

	fmt.Printf(fmtQuietPretty, bStatus, statusText, aStatus)
}
