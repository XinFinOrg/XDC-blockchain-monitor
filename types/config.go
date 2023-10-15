package types

import "sync"

type Config struct {
	Monitors      Monitor       `json:"monitor"`
	Rules         Rule          `json:"rules"`
	Blockchains   []*Blockchain `json:"blockchain"`
	Contracts     []*Contract   `json:"contract"`
	Notifications *Notification `json:"notification"`
}

type Monitor struct {
	FetchBlockNumber int
}

type Rule struct {
	Masternode       RuleStatus
	Epoch            RuleStatus
	Minetime         RuleStatus
	ContiguousRounds RuleStatus
	Confirmed        RuleStatus
}

type RuleStatus struct {
	Rate   float64
	Active bool
}

type Blockchain struct {
	Name       string `json:"name"`
	Stats      string `json:"stats"`
	RPC1       string `json:"rpc1"`
	RPC2       string `json:"rpc2,omitempty"`
	MineTime   int    `json:"mineTime"`
	Masternode int    `json:"masternode"`
	Test       string `json:"test" slack:"ignore"`
	Active     bool   `json:"active" slack:"ignore"`
	Hotstuff   bool   `json:"hotstuff" slack:"ignore"`

	BlockCache               map[int]*BlockRPCResult `slack:"ignore"`
	BlockCacheLock           *sync.Mutex             `slack:"ignore"`
	FetchBlockNumber         int                     `slack:"ignore"`
	LatestFetchedBlockNumber int
	LatestFetchedEpochNumber int
}

type Contract struct {
	Name string `json:"name"`
	Test string `json:"test"`
}

type Notification struct {
	Slack    []SlackNotification  `json:"slack"`
	Telegram TelegramNotification `json:"telegram"`
}

type SlackNotification struct {
	Services []string       `json:"services"`
	Token    string         `json:"token"`
	Channel  string         `json:"channel"`
	Tag      []SlackUserTag `json:"tag"`
}

type SlackUserTag struct {
	Name   string `json:"name"`
	UserID string `json:"userid"`
	Active bool   `json:"active"`
}

type TelegramNotification struct {
	Token    string   `json:"token"`
	Services []string `json:"services"`
	Channel  []struct {
		Name   string `json:"name"`
		ChatID int    `json:"chatid"`
		Active bool   `json:"active"`
	} `json:"channel"`
	Tag string `json:"tag"`
}
