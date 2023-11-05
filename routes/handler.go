package routes

import (
	"fmt"
	"net/http"

	"github.com/XinFinOrg/XDC-blockchain-monitor/data"
	"github.com/XinFinOrg/XDC-blockchain-monitor/notification"
	"github.com/gin-gonic/gin"
)

// Define the handler for the /deploy endpoint
func handleDeploy(c *gin.Context) {
	type DeployParams struct {
		Service     string `form:"service"`
		Version     string `form:"version"`
		Environment string `form:"environment"`
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
		payload := notification.BuildDeployMessage(params.Service, params.Version, params.Environment, tags, slackConfig.DeployChannel)
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
