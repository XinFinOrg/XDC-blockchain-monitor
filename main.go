package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/XinFinOrg/XDC-blockchain-monitor/cronjob"
	"github.com/XinFinOrg/XDC-blockchain-monitor/routes"
	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
	"github.com/joho/godotenv"
)

func main() {
	env()
	config := config()

	// Setup the cronjob using the parsed configuration
	cronManager := cronjob.SetupCron(&config)
	cronManager.Start()

	r := routes.SetupRouter()
	r.Run(":8080")
}
func env() {
	err := godotenv.Load(".env")
	if err != nil {
		// Handle error if the .env file is not found or has errors.
		panic("Error loading .env file")
	}
}

func config() types.Config {
	// Read and parse the config.json
	configData, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config.json: %s", err)
	}

	notificationData, err := os.ReadFile("config-notification.json")
	if err != nil {
		log.Fatalf("Error reading config.json: %s", err)
	}

	var config types.Config
	var notification types.Notification
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("Error parsing config.json: %s", err)
	}

	err = json.Unmarshal(notificationData, &notification)
	if err != nil {
		log.Fatalf("Error parsing config-notification.json: %s", err)
	}

	config.Notifications = &notification

	for i := range config.Blockchains {
		config.Blockchains[i].BlockCache = make(map[int]*types.BlockRPCResult)
		config.Blockchains[i].BlockCacheLock = &sync.Mutex{}
		config.Blockchains[i].FetchBlockNumber = config.Monitors.FetchBlockNumber
	}

	// Read from env
	config.Notifications.Slack[0].Token = os.Getenv("SLACK_BOT_TOKEN")
	config.Notifications.Telegram.Token = os.Getenv("TELEGRAM_BOT_TOKEN")
	return config
}
