package output

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	fmtJsonRaw string = "%v,%v,%v\n"

	// three sections: status text, expected value, url
	// for each section:
	// structure: color before, data, color after, comma
	fmtJsonPretty = "%v%v%v,%v%v%v,%v%v%v\n"
)

func WriteToJson(result *Result) {
	b, err := json.Marshal(result)
	if err != nil {
		log.Println("could not marshall json error:", err)
	} else {
		fmt.Println(string(b))
	}
}

func WriteToJsonPretty(result *Result) {
	b, err := json.MarshalIndent(&result, "", "\t")
	if err != nil {
		log.Println("could not marshall json error:", err)
	} else {
		fmt.Println(string(b))
	}
	/*
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

		fmt.Printf(fmtJsonPretty, bStatus, statusText, aStatus, bExpected, result.Expected, aExpected, bUrl, result.Url, aUrl)
	*/
}
