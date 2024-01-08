package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func CheckMineTime(config *types.Config, bc *types.Blockchain) error {
	lock := bc.BlockCacheLock

	lock.Lock()
	defer lock.Unlock()

	blockCache := bc.BlockCache
	targetTime := bc.MineTime
	latest := bc.LatestFetchedBlockNumber
	fetchNum := bc.FetchBlockNumber

	log.Println("Check Mine Time on blocks:", len(blockCache))
	if len(blockCache) == 0 {
		return errors.New("no blocks provided")
	}

	violations := 0
	for i := latest; i > latest-fetchNum+1; i-- {
		prevBlockTimestamp := parseHexToInt(blockCache[i-1].EthBlockResult.Timestamp)
		currBlockTimestamp := parseHexToInt(blockCache[i].EthBlockResult.Timestamp)

		// Get the time difference between this block and the previous block
		timeDiff := currBlockTimestamp - prevBlockTimestamp
		if timeDiff < targetTime-1 || timeDiff > targetTime+1 {
			violations++
		}
	}

	if float64(violations)/float64(len(blockCache)) > config.Rules.Minetime.Rate {
		return fmt.Errorf("more than %d of the blocks have inconsistent mining times", violations)
	}

	// Total Time
	currBlockTimestamp := parseHexToInt(blockCache[latest].EthBlockResult.Timestamp)
	prevBlockTimestamp := parseHexToInt(blockCache[latest-fetchNum+1].EthBlockResult.Timestamp)

	// Get the time difference between this block and the previous block
	timeDiff := currBlockTimestamp - prevBlockTimestamp
	if float64(timeDiff) >= float64(fetchNum*targetTime)*(1+config.Rules.Minetime.Rate) {
		return fmt.Errorf("bc mined too slow total time: %d sec, total blocks: %d", timeDiff, fetchNum)
	}

	return nil
}
