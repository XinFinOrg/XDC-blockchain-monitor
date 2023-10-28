package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/XinFinOrg/XDC-blockchain-monitor/data"
	"github.com/XinFinOrg/XDC-blockchain-monitor/notification"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/slack/button-click", handleButtonClick)
	r.GET("/blockCache", getBlockCache)
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	return r
}
func getBlockCache(c *gin.Context) {
	config := data.GetCurrentConfig()
	network := c.Query("network")
	for _, i := range config.Blockchains {
		if i.Name == network {
			c.JSON(http.StatusOK, i.BlockCache)
			return
		}
	}
	c.JSON(http.StatusBadRequest, fmt.Sprintf("Network %s Not Found", network))
}

func handleButtonClick(c *gin.Context) {
	// Parse the URL-encoded payload from Slack
	payloadStr := c.DefaultPostForm("payload", "")
	payloadStr, err := url.QueryUnescape(payloadStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to unescape payload"})
		return
	}

	// Debug: Print the decoded payload
	fmt.Printf("Received Slack payload: %s\n", payloadStr)

	var payload SlackPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON payload"})
		return
	}

	buttonClickedText := payload.Actions[0].Text.Text

	// Determine the action based on the button text
	var actionMessage string
	switch buttonClickedText {
	case "Acknowledge":
		actionMessage = "has Acknowledged this issue"
	case "Ignore":
		actionMessage = "has Ignored this issue"
	default:
		// Handle other button texts or unknown buttons
		fmt.Println("Unknown action", buttonClickedText)
		c.JSON(http.StatusOK, gin.H{"message": "Unknown action"})
		return
	}

	// You can modify this part to include the message timestamp and channel ID
	// of the original message that you want to update.
	// For demonstration purposes, we'll use placeholders.
	messageTimestamp := payload.Message.Ts
	channelID := payload.Channel.ID

	// Construct the updated message
	username := payload.User.Username
	updatedMessage := fmt.Sprintf("<@%s> %s", username, actionMessage)
	fmt.Println(updatedMessage)
	// Create a slice of blocks without the buttons block
	updatedBlocks := []Block{}

	for _, block := range payload.Message.Blocks {
		// Exclude the buttons block by checking its type
		if block.Type != "actions" {
			updatedBlocks = append(updatedBlocks, block)
		}
	}
	resBlock := Block{
		Type: "section",
		Fields: []Field{
			{
				Type: "mrkdwn",
				Text: updatedMessage,
			},
		},
	}
	updatedBlocks = append(updatedBlocks, resBlock)

	// Send the updated message back to Slack
	slackResponse := map[string]interface{}{
		"channel": channelID,
		"text":    updatedMessage,
		"blocks":  updatedBlocks,
		"ts":      messageTimestamp,
	}
	err = notification.Update(slackResponse)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to update payload %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}
	c.JSON(http.StatusOK, slackResponse)
	return
}
