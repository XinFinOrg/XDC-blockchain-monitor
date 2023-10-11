package types

import "encoding/json"

type XDPoSBlockResponse struct {
	Jsonrpc string           `json:"jsonrpc"`
	ID      int              `json:"id"`
	Result  XDPoSBlockResult `json:"result"`
}

type XDPoSBlockResult struct {
	Hash       string `json:"Hash"`
	Round      int    `json:"Round"`
	Number     int    `json:"Number"`
	ParentHash string `json:"ParentHash"`
	Committed  bool   `json:"Committed"`
	Miner      string `json:"Miner"`
	Timestamp  int    `json:"Timestamp"`
	// EncodedRLP string `json:"EncodedRLP"`
	Error string `json:"Error"`
}

func (r *XDPoSBlockResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
