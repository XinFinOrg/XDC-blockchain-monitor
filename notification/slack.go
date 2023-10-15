package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func SendToSlack(config *types.Config, bc *types.Blockchain, message string) {
	// Send notification to Slack
	for _, slackConfig := range config.Notifications.Slack {
		if contains(slackConfig.Services, bc.Name) {
			webhookURL := slackConfig.Url
			message = "*" + message + "*\n" + GetMessageForSlack(bc)
			send(webhookURL, message, slackConfig)
		}
	}
}
func send(webhookURL string, message string, slackConfig types.SlackNotification) error {
	tags := ""
	message = message + "\n"
	for _, v := range slackConfig.Tag {
		if v.Active {
			tags += fmt.Sprintf(" <@%s>", v.UserID)
		}
	}
	payload := map[string]string{
		"text": message + tags,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Failed to send notification to Slack. StatusCode: %d, Response: %s", resp.StatusCode, body)
	}

	return nil
}
