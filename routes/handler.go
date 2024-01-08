package routes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/XinFinOrg/XDC-blockchain-monitor/data"
	"github.com/XinFinOrg/XDC-blockchain-monitor/notification"
	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
	"github.com/gin-gonic/gin"
)

// Define the handler for the /deploy endpoint
func handleDeploy(c *gin.Context) {
	type DeployParams struct {
		Service     string `form:"service"`
		Version     string `form:"version"`
		Environment string `form:"environment"`
		Comment     string `form:"comment"`
	}

	var params DeployParams

	// Bind the query parameters to the struct
	if c.ShouldBindQuery(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	config := data.GetCurrentConfig()
	for _, slackConfig := range config.Notifications.Slack {
		tags := ""
		for _, v := range slackConfig.Tag {
			if v.Active && contains(v.Environments, params.Environment) {
				tags += fmt.Sprintf(" <@%s>", v.UserID)
			}
		}
		payload := notification.BuildDeployMessage(params.Service, params.Version, params.Environment, params.Comment, tags, slackConfig.DeployChannel)
		// Send the notification to Slack
		err := notification.Send(payload, slackConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification to Slack"})
			return
		}
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Deployment notification sent successfully"})
}

func handleSnooze(bc string) {
	config := data.GetCurrentConfig()

	for _, v := range config.Blockchains {
		if v.Name == bc {
			// Launch a goroutine to handle the snooze
			go func(blockchain *types.Blockchain) {
				blockchain.Active = false
				// Sleep for 3 hours
				log.Println("Snooze period start for blockchain:", blockchain.Name)
				time.Sleep(3 * time.Hour)
				// After 3 hours, update the Active status to true
				blockchain.Active = true

				// Optionally, you can perform additional actions after snooze period
				// For example, sending a notification or logging the event
				log.Println("Snooze period ended for blockchain:", blockchain.Name)
			}(v)
		}
	}
}
