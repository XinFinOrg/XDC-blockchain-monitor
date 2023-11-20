// cronjob/jobs.go

package cronjob

import (
	"log"
	"time"

	"github.com/XinFinOrg/XDC-blockchain-monitor/notification"
	"github.com/XinFinOrg/XDC-blockchain-monitor/service"
	"github.com/XinFinOrg/XDC-blockchain-monitor/types"

	"github.com/robfig/cron/v3"
)

// Track the time of the last notification
var lastNotificationTime_FetchBlocks time.Time
var lastNotificationTime_CheckMineTime time.Time
var lastNotificationTime_Hotstuff time.Time

// Define counters for each type of error
var errorCountFetchBlocks int
var errorCountCheckMineTime int
var errorCountHotstuff int

func SetupCron(config *types.Config) *cron.Cron {
	c := cron.New(cron.WithSeconds())

	for _, bc := range config.Blockchains {
		// Here, you capture the bc variable to avoid Go's loop variable capture behavior
		bc := bc
		if !bc.Active {
			continue
		}

		c.AddFunc("@every 30s", func() {
			if err := service.FetchBlocks(config, bc); err != nil {
				log.Println("Fetch Blocks Error: ", err, bc.Name)
				errorCountFetchBlocks++
				if errorCountFetchBlocks >= 3 {
					if elapsedTime := time.Since(lastNotificationTime_FetchBlocks); elapsedTime > time.Hour {
						notification.SendToTelegram(config, bc, err)
						notification.AlertSendToSlack(config, bc, err)
						lastNotificationTime_FetchBlocks = time.Now()
						errorCountFetchBlocks = 0
					}
				}
			} else {
				errorCountFetchBlocks = 0
			}
			if err := service.CheckMineTime(config, bc); err != nil {
				log.Println("CheckMineTime Error: ", err, bc.Name)
				errorCountCheckMineTime++
				if errorCountCheckMineTime >= 3 {
					if elapsedTime := time.Since(lastNotificationTime_CheckMineTime); elapsedTime > time.Hour {
						notification.SendToTelegram(config, bc, err)
						notification.AlertSendToSlack(config, bc, err)
						lastNotificationTime_CheckMineTime = time.Now()
						errorCountCheckMineTime = 0
					}
				}
			} else {
				errorCountCheckMineTime = 0
			}
			if err := service.Hotstuff(config, bc); err != nil {
				log.Println("Hotstuff Error: ", err, bc.Name)
				errorCountHotstuff++
				if errorCountHotstuff >= 3 {
					if elapsedTime := time.Since(lastNotificationTime_Hotstuff); elapsedTime > time.Hour {
						notification.SendToTelegram(config, bc, err)
						notification.AlertSendToSlack(config, bc, err)
						lastNotificationTime_Hotstuff = time.Now()
						errorCountHotstuff = 0
					}
				}
			} else {
				errorCountHotstuff = 0
			}
		})
		c.AddFunc("@every 1h", func() {
			if err := service.FetchEpoch(config, bc); err != nil {
				log.Println("FetchEpoch Error: ", err, bc.Name)
				notification.SendToTelegram(config, bc, err)
				notification.AlertSendToSlack(config, bc, err)
			}
		})
	}

	return c
}
