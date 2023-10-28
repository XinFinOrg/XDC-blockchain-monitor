package main

import (
	"github.com/XinFinOrg/XDC-blockchain-monitor/cronjob"
	"github.com/XinFinOrg/XDC-blockchain-monitor/data"
	"github.com/XinFinOrg/XDC-blockchain-monitor/routes"
	"github.com/joho/godotenv"
)

func main() {
	env()
	config := data.Config()

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
