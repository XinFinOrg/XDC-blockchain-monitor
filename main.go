package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/liam-lai/xinfin-monitor/cronjob"
	"github.com/liam-lai/xinfin-monitor/routes"
	"github.com/liam-lai/xinfin-monitor/types"
)

func main() {
	config := config()

	// Setup the cronjob using the parsed configuration
	cronManager := cronjob.SetupCron(&config)
	cronManager.Start()

	r := routes.SetupRouter()
	r.Run(":8080")
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
	return config
}
