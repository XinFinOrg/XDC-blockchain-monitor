package service

import (
	"fmt"
	"log"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func Hotstuff(config *types.Config, bc *types.Blockchain) error {
	if !bc.Hotstuff {
		return nil
	}

	log.Println("Check Hotstuff Parameters on", bc.Name)
	if err := masternode(config, bc); err != nil {
		return err
	}
	if err := confirmBlock(config, bc); err != nil {
		return err
	}
	if err := checkContiguousRounds(config, bc); err != nil {
		return err
	}
	return nil
}

func checkContiguousRounds(config *types.Config, bc *types.Blockchain) *types.ErrorMonitor {
	nonContiguousCount := 0
	latest := bc.LatestFetchedBlockNumber
	blockCache := bc.BlockCache
	details := ""
	for i := latest - len(blockCache) + 1; i < latest; i++ {
		// Check the difference in Round between consecutive blocks
		next := blockCache[i+1].XDPoSBlockResult
		cur := blockCache[i].XDPoSBlockResult
		difference := next.Round - cur.Round
		if difference != 1 {
			details += fmt.Sprintf("block: %d, round: %d, next block: %d, round: %d\n", cur.Number, cur.Round, next.Number, next.Round)
			nonContiguousCount++
		}
	}

	if float64(nonContiguousCount) > float64(len(blockCache))*config.Rules.ContiguousRounds.Rate {
		e := &types.ErrorMonitor{
			Title:   fmt.Sprintf("more than %d of the blocks in total %d have non-contiguous rounds", nonContiguousCount, len(blockCache)),
			Details: details,
		}
		return e
	}

	return nil
}

func confirmBlock(config *types.Config, bc *types.Blockchain) *types.ErrorMonitor {
	latest := bc.LatestFetchedBlockNumber
	fmt.Println("latest", latest, "Confirmed Rate", int(config.Rules.Confirmed.Rate))
	if !bc.BlockCache[latest-int(config.Rules.Confirmed.Rate)].XDPoSBlockResult.Committed {
		details := ""
		for i := bc.LatestFetchedBlockNumber; i > bc.LatestFetchedBlockNumber-len(bc.BlockCache); i-- {
			block := bc.BlockCache[i].XDPoSBlockResult
			details += fmt.Sprintf("block: %d, round: %d, committed: %t, timestamp: %d\n", block.Number, block.Round, block.Committed, block.Timestamp)
		}
		e := &types.ErrorMonitor{
			Title:   fmt.Sprintf("last %d blocks are not committed", int(config.Rules.Confirmed.Rate)),
			Details: details,
		}
		return e
	}
	return nil
}
