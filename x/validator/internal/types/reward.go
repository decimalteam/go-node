package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	lastBlock   = 46_656_000
	firstReward = 50

	firstIncrease = 5
)

func GetRewardForBlock(blockHeight uint64) sdk.Int {
	if blockHeight > lastBlock {
		return sdk.NewInt(0)
	}

	reward := sdk.NewInt(firstReward)
	rewardIncrease := sdk.NewInt(firstIncrease)

	rewardIncrease = rewardIncrease.Add(sdk.NewInt(int64(blockHeight/5184000) * 12))
	reward = reward.Add(sdk.NewInt(int64(blockHeight / 432000)).Mul(rewardIncrease))

	return reward
}
