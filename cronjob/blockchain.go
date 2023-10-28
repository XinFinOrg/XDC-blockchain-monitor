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
				if elapsedTime := time.Since(lastNotificationTime_FetchBlocks); elapsedTime > time.Hour {
					notification.SendToTelegram(config, bc, err)
					notification.SendToSlack(config, bc, err)
					lastNotificationTime_FetchBlocks = time.Now()
				}
			}
			if err := service.CheckMineTime(config, bc); err != nil {
				log.Println("CheckMineTime Error: ", err, bc.Name)
				if elapsedTime := time.Since(lastNotificationTime_CheckMineTime); elapsedTime > time.Hour {
					log.Println("CheckMineTime Error: ", err, bc.Name)
					notification.SendToTelegram(config, bc, err)
					notification.SendToSlack(config, bc, err)
					lastNotificationTime_CheckMineTime = time.Now()
				}
			}
			if err := service.Hotstuff(config, bc); err != nil {
				log.Println("Hotstuff Error: ", err, bc.Name)
				if elapsedTime := time.Since(lastNotificationTime_Hotstuff); elapsedTime > time.Hour {
					notification.SendToTelegram(config, bc, err)
					notification.SendToSlack(config, bc, err)
					lastNotificationTime_Hotstuff = time.Now()
				}
			}
		})
		c.AddFunc("@every 1h", func() {
			if err := service.FetchEpoch(config, bc); err != nil {
				log.Println("FetchEpoch Error: ", err, bc.Name)
				notification.SendToTelegram(config, bc, err)
				notification.SendToSlack(config, bc, err)
			}
		})
	}

	return c
}
