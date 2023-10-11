package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/liam-lai/xinfin-monitor/types"
)

func FetchBlocks(config *types.Config, bc *types.Blockchain) error {
	lock := bc.BlockCacheLock
	lock.Lock()
	defer lock.Unlock()

	RPC := bc.RPC1
	blockCache := bc.BlockCache
	log.Println("fetchBlocks on blockchain", bc.Name, " function called:", RPC)

	latestBlock, _, err := getBlockWithNumber(RPC, "latest")
	if err != nil {
		return err
	}
	if latestBlock == nil || latestBlock.Result.Number == "" {
		return fmt.Errorf("fetch blocks error")
	}

	latestBlockNumber := parseHexToInt(latestBlock.Result.Number)
	startBlock := latestBlockNumber - bc.FetchBlockNumber + 1

	if bc.LatestFetchedBlockNumber == latestBlockNumber {
		return fmt.Errorf("block doesn't generate during monitor period!! Need to check")
	}

	bc.LatestFetchedBlockNumber = latestBlockNumber
	log.Println("set latest block number", latestBlockNumber)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := startBlock; i <= latestBlockNumber; i++ {
		if _, exists := blockCache[i]; !exists {
			wg.Add(1)
			go func(blockNum int) {
				defer wg.Done()

				blockHex := fmt.Sprintf("0x%x", blockNum)
				ethBlock, xdPosBlock, err := getBlockWithNumber(RPC, blockHex)
				if err != nil {
					log.Printf("getBlockWithNumber function called with rpcURL: %s, blockNumber: %d, error %s", RPC, blockNum, err)
					return
				}
				mu.Lock()
				blockCache[blockNum] = &types.BlockRPCResult{
					EthBlockResult: ethBlock.Result,
				}
				if xdPosBlock != nil {
					blockCache[blockNum].XDPoSBlockResult = xdPosBlock.Result
				}
				mu.Unlock()
			}(i)
		}
	}

	wg.Wait()

	// Trim cache to latest 900 blocks
	for k := range blockCache {
		// n, _ := hexToNumber(blockCache[k].EthBlockResult.Number)
		// log.Printf("current block in map blockNumber: %d.", n)

		if k < startBlock {
			delete(blockCache, k)
		}
	}
	return nil
}

func getBlockWithNumber(rpcURL string, blockNum string) (*types.EthBlockResponse, *types.XDPoSBlockResponse, error) {
	// Construct the request payload for EthBlockRequest
	requestPayload := types.EthBlockRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNum, true},
		ID:      1,
	}

	var ethBlockResponse types.EthBlockResponse
	err := sendRequest(rpcURL, requestPayload, &ethBlockResponse)
	if err != nil {
		return nil, nil, err
	}

	if ethBlockResponse.ID == 0 || ethBlockResponse.Result.Miner == "0x0000000000000000000000000000000000000000" {
		return &ethBlockResponse, nil, err
	}

	// Construct the request payload for XDPoSBlock
	xdPosBlockRequest := types.EthBlockRequest{
		Jsonrpc: "2.0",
		Method:  "XDPoS_getV2BlockByNumber",
		Params:  []interface{}{blockNum},
		ID:      1,
	}

	var xdPosBlockResponse types.XDPoSBlockResponse
	err = sendRequest(rpcURL, xdPosBlockRequest, &xdPosBlockResponse)
	return &ethBlockResponse, &xdPosBlockResponse, err
}
