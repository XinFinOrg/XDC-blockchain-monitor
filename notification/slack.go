package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
	"github.com/slack-go/slack"
)

func AlertSendToSlack(config *types.Config, bc *types.Blockchain, err error) {
	// Send notification to Slack
	for _, slackConfig := range config.Notifications.Slack {
		details := GetMessageForSlack(bc)
		title := err.Error()

		buildSlackMessage(bc.Name, title, details, slackConfig)
		if customErr, ok := err.(*types.ErrorMonitor); ok {
			SendDebugMsg(title, customErr.Details, slackConfig)
		}
	}
}

func buildSlackMessage(service string, title string, details string, slackConfig types.SlackNotification) error {
	tags := ""
	details = details + "\n"
	for _, v := range slackConfig.Tag {
		if v.Active && contains(v.Environments, service) {
			tags += fmt.Sprintf(" <@%s>", v.UserID)
		}
	}

	// No one wants to sub this service
	if len(tags) == 0 {
		return nil
	}
	// Create the payload with the attachment.

	payload := buildAlertMessage(title, details+tags, slackConfig.AlertChannel)
	err := Send(payload, slackConfig)
	return err
}
func Send(payload SlackMessage, slackConfig types.SlackNotification) error {

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Define the Slack API endpoint URL
	apiURL := "https://slack.com/api/chat.postMessage"

	// Create an HTTP POST request to send the message
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set the Slack API token in the Authorization header
	req.Header.Set("Authorization", "Bearer "+slackConfig.Token)
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send notification to Slack. StatusCode: %d, Response: %s", resp.StatusCode, body)
	}

	return nil
}

func SendDebugMsg(title string, msg string, slackConfig types.SlackNotification) {

	channelID := slackConfig.AlertChannel // Replace with your channel ID

	api := slack.New(slackConfig.Token)

	// Create a snippet
	fileUploadParameters := slack.FileUploadParameters{
		Channels: []string{channelID},
		Content:  msg,
		Filename: "debug.txt",
		Title:    title,
		Filetype: "text",
	}

	_, err := api.UploadFile(fileUploadParameters)
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	log.Println("Snippet uploaded successfully!")
}

func Update(slackResponse map[string]interface{}) error {
	// Convert payload to JSON
	updateJSON, err := json.Marshal(slackResponse)
	if err != nil {
		return err
	}

	// Define the Slack API endpoint URL for updating messages
	apiURL := "https://slack.com/api/chat.update"

	// Create an HTTP POST request to update the message
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(updateJSON))
	if err != nil {
		return err
	}

	// Set the Slack API token and Content-Type header
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SLACK_BOT_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code and handle any errors
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to update message in Slack. StatusCode: %d, Response: %s", resp.StatusCode, body)
		// Handle the error accordingly
		// You may want to return an error response or log the error.
	}

	return nil
}
