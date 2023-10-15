package routes

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	TeamID   string `json:"team_id"`
}

type Container struct {
	Type        string `json:"type"`
	MessageTS   string `json:"message_ts"`
	ChannelID   string `json:"channel_id"`
	IsEphemeral bool   `json:"is_ephemeral"`
}

type Team struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Text struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

type Block struct {
	Type     string   `json:"type"`
	BlockID  string   `json:"block_id"`
	Text     *Text    `json:"text,omitempty"`
	Elements []Button `json:"elements,omitempty"`
	Fields   []Field  `json:"fields,omitempty"`
}

type Field struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Verbatim bool   `json:"verbatim"`
}

type Button struct {
	Type     string `json:"type"`
	ActionID string `json:"action_id"`
	Text     Text   `json:"text"`
	Style    string `json:"style"`
	Value    string `json:"value"`
}

type Message struct {
	BotID  string  `json:"bot_id"`
	Type   string  `json:"type"`
	Text   string  `json:"text"`
	User   string  `json:"user"`
	Ts     string  `json:"ts"`
	AppID  string  `json:"app_id"`
	Blocks []Block `json:"blocks"`
	Team   string  `json:"team"`
}

type State struct {
	Values map[string]interface{} `json:"values"`
}

type SlackPayload struct {
	Type                string      `json:"type"`
	User                User        `json:"user"`
	APIAppID            string      `json:"api_app_id"`
	Token               string      `json:"token"`
	Container           Container   `json:"container"`
	TriggerID           string      `json:"trigger_id"`
	Team                Team        `json:"team"`
	Enterprise          interface{} `json:"enterprise"`
	IsEnterpriseInstall bool        `json:"is_enterprise_install"`
	Channel             Channel     `json:"channel"`
	Message             Message     `json:"message"`
	State               State       `json:"state"`
	ResponseURL         string      `json:"response_url"`
	Actions             []Button    `json:"actions"`
}
