package notification

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/liam-lai/xinfin-monitor/types"
)

const (
	BotToken = "" // TODO READ From .env
)

func SendToTelegram(config *types.Config, bc *types.Blockchain, message string) {
	// Send notification to Slack
	for _, teleConfig := range config.Notifications.Telegram {
		if contains(teleConfig.Services, bc.Name) {
			message = "*" + message + "*\n" + GetMessageForSlack(bc)
			for _, channel := range teleConfig.Channel {
				if channel.Active {
					sendMessage(channel.ChatID, message)
				}
			}
		}
	}
}

func sendMessage(chatID int, message string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BotToken)

	data := url.Values{}
	data.Set("chat_id", strconv.Itoa(chatID))
	data.Set("text", message)

	response, err := http.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	/*
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(body)) // You can process the response if needed
	*/
	return nil
}
