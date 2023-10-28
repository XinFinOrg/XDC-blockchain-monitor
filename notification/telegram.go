package notification

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func SendToTelegram(config *types.Config, bc *types.Blockchain, err error) {
	// Send notification to Slack
	message := err.Error()
	teleConfig := config.Notifications.Telegram
	if contains(teleConfig.Services, bc.Name) {
		message = "*" + message + "*\n" + GetMessageForSlack(bc)
		for _, channel := range teleConfig.Channel {
			if channel.Active {
				sendMessage(channel.ChatID, message, config.Notifications.Telegram.Token)
			}
		}
	}

}

func sendMessage(chatID int, message string, token string) error {
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

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
