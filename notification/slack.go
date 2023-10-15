package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func SendToSlack(config *types.Config, bc *types.Blockchain, title string) {
	// Send notification to Slack
	for _, slackConfig := range config.Notifications.Slack {
		if contains(slackConfig.Services, bc.Name) {
			details := GetMessageForSlack(bc)
			send(title, details, slackConfig)
		}
	}
}

func send(title string, details string, slackConfig types.SlackNotification) error {
	tags := ""
	details = details + "\n"
	for _, v := range slackConfig.Tag {
		if v.Active {
			tags += fmt.Sprintf(" <@%s>", v.UserID)
		}
	}
	// Create the payload with the attachment.

	payload := buildMessage(title, details+tags, slackConfig.Channel)

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

	// Parse the JSON response
	var slackResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&slackResponse)
	if err != nil {
		return err
	}

	// Check if the request was successful
	if slackResponse["ok"].(bool) {
		// Request was successful, you can access other fields in slackResponse as needed
		ts := slackResponse["ts"].(string)
		channel := slackResponse["channel"].(string)
		message := slackResponse["message"].(map[string]interface{})
		//text := message["text"].(string)
		//user := message["user"].(string)

		// Handle the response data...
		fmt.Printf("Message sent successfully. Timestamp: %s, Channel: %s, Message: %v\n", ts, channel, message)
	} else {
		// Request failed, handle the error...
		fmt.Println("Slack API error:", slackResponse["error"].(string))
	}

	return nil
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
	} else {
		// Message updated successfully
		fmt.Println("Message updated successfully")
	}
	return nil
}
