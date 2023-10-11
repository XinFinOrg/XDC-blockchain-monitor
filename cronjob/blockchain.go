// cronjob/jobs.go

package cronjob

import (
	"log"

	"github.com/liam-lai/xinfin-monitor/notification"
	"github.com/liam-lai/xinfin-monitor/service"
	"github.com/liam-lai/xinfin-monitor/types"

	"github.com/robfig/cron/v3"
)

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
				notification.SendToTelegram(config, bc, err.Error())
				notification.SendToSlack(config, bc, err.Error())
			}
			if err := service.CheckMineTime(config, bc); err != nil {
				log.Println("CheckMineTime Error: ", err, bc.Name)
				notification.SendToTelegram(config, bc, err.Error())
				notification.SendToSlack(config, bc, err.Error())
			}
			if err := service.Hotstuff(config, bc); err != nil {
				log.Println("Hotstuff Error: ", err, bc.Name)
				notification.SendToTelegram(config, bc, err.Error())
				notification.SendToSlack(config, bc, err.Error())
			}
		})
		c.AddFunc("@every 1h", func() {
			if err := service.FetchEpoch(config, bc); err != nil {
				log.Println("FetchEpoch Error: ", err, bc.Name)
				notification.SendToTelegram(config, bc, err.Error())
				notification.SendToSlack(config, bc, err.Error())
			}
		})
	}

	return c
}
