package service

import (
	"fmt"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func FetchEpoch(config *types.Config, bc *types.Blockchain) error {
	rpc := bc.RPC1
	// Construct the request payload for EthBlockRequest
	requestPayload := types.EthBlockRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getCandidates",
		Params:  []interface{}{"latest"},
		ID:      1,
	}
	var candidatesResponse types.GetCandidatesResponse
	err := sendRequest(rpc, requestPayload, &candidatesResponse)
	if err != nil {
		return err
	}

	if bc.LatestFetchedEpochNumber == int(candidatesResponse.Result.Epoch) {
		return fmt.Errorf("epoch number doesn't update during 2 fetch period")
	}
	bc.LatestFetchedEpochNumber = int(candidatesResponse.Result.Epoch)

	return nil
}
