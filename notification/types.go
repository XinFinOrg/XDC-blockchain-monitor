package notification

// SlackMessage represents the structure of a Slack message payload.
type SlackMessage struct {
	Channel string  `json:"channel"`
	Text    string  `json:"text"`
	Blocks  []Block `json:"blocks"`
}

// Block represents a block within the Slack message payload.
type Block struct {
	Type     string    `json:"type"`
	Text     *Text     `json:"text,omitempty"`
	Fields   []Field   `json:"fields,omitempty"`
	Elements []Element `json:"elements,omitempty"`
}

// Text represents the text structure within a block.
type Text struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

// Field represents a field structure within a section block.
type Field struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Element represents an element within an actions block.
type Element struct {
	Type  string `json:"type"`
	Text  Text   `json:"text"`
	Style string `json:"style"`
	Value string `json:"value"`
}
