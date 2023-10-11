package types

import "encoding/json"

type GetMasternodesByNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Number         int      `json:"Number"`
		Round          int      `json:"Round"`
		MasternodesLen int      `json:"MasternodesLen"`
		Masternodes    []string `json:"Masternodes"`
		PenaltyLen     int      `json:"PenaltyLen"`
		Penalty        []string `json:"Penalty"`
	} `json:"result"`
}

func (r *GetMasternodesByNumberResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
