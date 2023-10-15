package service

import (
	"fmt"
	"log"

	"github.com/XinFinOrg/XDC-blockchain-monitor/types"
)

func Hotstuff(config *types.Config, bc *types.Blockchain) error {
	log.Println("Check Hotstuff Parameters on", bc.Name)
	if !bc.Hotstuff {
		return nil
	}
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

func checkContiguousRounds(config *types.Config, bc *types.Blockchain) error {
	nonContiguousCount := 0
	latest := bc.LatestFetchedBlockNumber
	blockCache := bc.BlockCache

	for i := latest - len(blockCache) + 1; i < latest; i++ {
		// Check the difference in Round between consecutive blocks
		difference := blockCache[i+1].XDPoSBlockResult.Round - blockCache[i].XDPoSBlockResult.Round
		if difference != 1 {
			nonContiguousCount++
		}
	}

	if nonContiguousCount > len(blockCache)*int(config.Rules.ContiguousRounds.Rate) {
		return fmt.Errorf("more than %d of the blocks in total %d have non-contiguous rounds", nonContiguousCount, len(blockCache))
	}

	return nil
}

func confirmBlock(config *types.Config, bc *types.Blockchain) error {
	latest := bc.LatestFetchedBlockNumber
	if !bc.BlockCache[latest-int(config.Rules.Confirmed.Rate)].XDPoSBlockResult.Committed {
		return fmt.Errorf("last 10 blocks is not committed")
	}
	return nil
}
