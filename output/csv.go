package output

import (
	"fmt"
)

const (
	fmtCsvRaw string = "%v,%v,%v\n"

	// three sections: status text, expected value, url
	// for each section:
	// structure: color before, data, color after, comma
	fmtCsvPretty = "%v%v%v,%v%v%v,%v%v%v\n"
)

func WriteToCsv(result *Result) {
	status := result.Success
	if result.Error != nil {
		status = false
	}

	fmt.Printf(fmtCsvRaw, status, result.Expected, result.Url)
}

func WriteToCsvPretty(result *Result) {
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

	fmt.Printf(fmtCsvPretty, bStatus, statusText, aStatus, bExpected, result.Expected, aExpected, bUrl, result.Url, aUrl)
}
