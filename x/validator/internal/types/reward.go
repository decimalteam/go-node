package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const lastBlock = 43702611
const firstReward = 333
const lastReward = 68

func GetRewardForBlock(blockHeight uint64) sdk.Int {
	if blockHeight > lastBlock {
		return sdk.NewInt(0)
	}

	if blockHeight == lastBlock {
		return sdk.NewInt(lastReward)
	}

	reward := sdk.NewInt(firstReward)
	reward = reward.Sub(sdk.NewInt(int64(blockHeight / 200000)))

	if reward.LT(sdk.NewInt(1)) {
		return sdk.NewInt(1)
	}

	return reward
}
