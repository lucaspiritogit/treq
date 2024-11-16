package response

import (
	"encoding/json"

	"github.com/rivo/tview"
)

func GetResponseTextView() *tview.TextView {
	responseTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetScrollable(true)
	responseTextView.SetBorder(true).SetTitle("Response")

	return responseTextView
}

func FormatJSON(input string) (string, error) {
	var prettyJSON any
	err := json.Unmarshal([]byte(input), &prettyJSON)
	if err != nil {
		return "", err
	}

	prettyJSONBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		return "", err
	}

	return string(prettyJSONBytes), nil
}
