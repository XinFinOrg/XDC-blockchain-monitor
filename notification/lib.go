package notification

import (
	"fmt"
	"reflect"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func BuildDeployMessage(service, version, environment, tags, channel string) SlackMessage {
	title := ":rocket: Deployment: " + environment
	fmt.Println("tags", tags)
	message := SlackMessage{
		Channel: channel,
		Text:    title,
		Blocks: []Block{
			{
				Type: "header",
				Text: &Text{
					Type:  "plain_text",
					Text:  title,
					Emoji: true,
				},
			},
			{
				Type: "section",
				Fields: []Field{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("service: *%s* \nenvironment: *%s* \n version: `%s`", service, environment, version),
					},
					{
						Type: "mrkdwn",
						Text: tags,
					},
				},
			},
		},
	}
	return message
}

func buildAlertMessage(title string, msg string, channel string) SlackMessage {
	message := SlackMessage{
		Channel: channel,
		Text:    title,
		Blocks: []Block{
			{
				Type: "header",
				Text: &Text{
					Type:  "plain_text",
					Text:  title,
					Emoji: true,
				},
			},
			{
				Type: "section",
				Fields: []Field{
					{
						Type: "mrkdwn",
						Text: msg,
					},
				},
			},
			{
				Type: "actions",
				Elements: []Element{
					{
						Type: "button",
						Text: Text{
							Type:  "plain_text",
							Text:  "Acknowledge",
							Emoji: true,
						},
						Style: "primary",
						Value: "acknowledge_button_click",
					},
					{
						Type: "button",
						Text: Text{
							Type:  "plain_text",
							Text:  "Ignore",
							Emoji: true,
						},
						Style: "danger",
						Value: "ignore_button_click",
					},
				},
			},
		},
	}

	return message
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetMessageForSlack(bc *types.Blockchain) string {
	v := reflect.ValueOf(bc)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var message string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		fieldName := fieldType.Name

		// Check if the field has the slack:"ignore" tag
		if tag, ok := fieldType.Tag.Lookup("slack"); ok && tag == "ignore" {
			continue
		}

		if !isZero(field) {
			message += fmt.Sprintf("%s: %v\n", fieldName, field.Interface())
		}
	}

	return message
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Ptr, reflect.Slice:
		return v.IsNil()
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	}
	// Consider other fields non-zero
	return false
}
