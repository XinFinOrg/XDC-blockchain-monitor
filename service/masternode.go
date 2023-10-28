package service

import (
	"fmt"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

var Rate = 0.9

func masternode(config *types.Config, bc *types.Blockchain) error {
	rpc := bc.RPC1
	// Construct the request payload for EthBlockRequest
	requestPayload := types.EthBlockRequest{
		Jsonrpc: "2.0",
		Method:  "XDPoS_getMasternodesByNumber",
		Params:  []interface{}{},
		ID:      1,
	}
	var masternodesResponse types.GetMasternodesByNumberResponse
	err := sendRequest(rpc, requestPayload, &masternodesResponse)
	if err != nil {
		return err
	}

	// Assuming GetMasternodesByNumberResponse has a field Masternodes which is a slice
	numberOfMasternodes := masternodesResponse.Result.MasternodesLen

	// Check if the number of masternodes is less than Rate of 108
	if numberOfMasternodes < int(Rate*float64(bc.Masternode)) {
		return fmt.Errorf("number of masternodes (%d) are too low, current total is %d", numberOfMasternodes, bc.Masternode)
	}

	return nil
}
