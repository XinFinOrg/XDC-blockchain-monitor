package types

import (
	"encoding/json"
	"math/big"
)

type GetCandidatesResponse struct {
	Jsonrpc string                  `json:"jsonrpc"`
	ID      int                     `json:"id"`
	Result  GetCandidatesResultType `json:"result"`
}

type GetCandidatesResultType struct {
	Candidates map[string]CandidateInfo `json:"candidates"`
	Epoch      int64                    `json:"epoch"`
	Success    bool                     `json:"success"`
}

type CandidateInfo struct {
	Capacity *big.Int `json:"capacity"`
	Status   string   `json:"status"`
}

func (r *GetCandidatesResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
