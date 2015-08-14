package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SlackIconUrl string

const (
	SlackIconGreenLAN  SlackIconUrl = "https://www.dropbox.com/s/3zyeho7fmbmrhjq/LanGreen.png?dl=1"
	SlackIconYellowLAN              = "https://www.dropbox.com/s/rrugtqf62epjeg5/LanYellow.png?dl=1"
	SlackIconRedLAN                 = "https://www.dropbox.com/s/lol19valntxoa15/LanRed.png?dl=1"

	SlackIconStatusIndicator = "https://www.dropbox.com/s/uolw04kx70jjroi/StatusIndicator.png?dl=1"
	SlackIconTaskManager     = "https://www.dropbox.com/s/92mq1cabokl1dnh/TaskManager.png?dl=1"
)

const (
	DefaultSlackIcon    SlackIconUrl = SlackIconStatusIndicator
	DefaultSlackChannel string       = "#general"
	DefaultSlackUser    string       = "asrt"
)

// SlackWriter implements io.WriteCloser.
// Note: calls to Write only write to an internal buffer.
// Calling Close performs an http.Post
type SlackWriter struct {
	Url     string
	Payload *slackPayload
	buffer  *bytes.Buffer
}

type slackPayload struct {
	Channel string       `json:"channel"`
	User    string       `json:"username"`
	Text    string       `json:"text"`
	Icon    SlackIconUrl `json:"icon_url"`
}

// Creates a new SlackWriter with sensible defaults
func NewSlackWriter(url string) *SlackWriter {
	return &SlackWriter{
		Url: url,
		Payload: &slackPayload{
			Channel: DefaultSlackChannel,
			User:    DefaultSlackUser,
			Icon:    DefaultSlackIcon,
		},
	}
}

// SlackWriter Write method writes only to an internal buffer.
// Caller must Close the writer in order for a the submission to occur
func (slack *SlackWriter) Write(p []byte) (n int, err error) {
	if slack.buffer == nil {
		slack.buffer = new(bytes.Buffer)
	}

	return slack.buffer.Write(p)
}

// Close method does the actual http.Post(...) to Slack
func (slack *SlackWriter) Close() error {
	if buf, err := slack.buildRequestBody(); err != nil {
		return fmt.Errorf("could not close slackwriter: %v", err)
	} else {
		if _, respErr := http.Post(slack.Url, "application/json", buf); respErr != nil {
			return respErr
		}
	}

	slack.reset()
	return nil
}

func (slack *SlackWriter) buildRequestBody() (io.Reader, error) {
	slack.Payload.Text = slack.buffer.String()
	jsonBytes, jsonErr := json.Marshal(&slack.Payload)
	if jsonErr != nil {
		return nil, fmt.Errorf("could not create json request body: %v", jsonErr)
	}

	buf := bytes.NewBuffer(jsonBytes)
	return buf, nil
}

func (slack *SlackWriter) reset() {
	slack.buffer = nil
	slack.Payload.Text = ""
}

func (slack *SlackWriter) SlackChannel(channel string) *SlackWriter {
	slack.Payload.Channel = channel
	return slack
}
func (slack *SlackWriter) SlackUser(user string) *SlackWriter {
	slack.Payload.User = user
	return slack
}
func (slack *SlackWriter) SlackIconUrl(icon SlackIconUrl) *SlackWriter {
	slack.Payload.Icon = icon
	return slack
}
