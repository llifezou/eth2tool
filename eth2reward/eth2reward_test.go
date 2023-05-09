package eth2reward

import (
	"fmt"
	"math/big"
	"testing"
)

func TestQueryEth2Rewards(t *testing.T) {
	rewards, err := QueryEth2Rewards("0x23715A59BEd8A94AA4FFebE8B7f1125b84FE970a", "")
	if err != nil {
		t.Fatal(err)
	}

	totalRewards := big.NewInt(0)
	validators := make(map[string]struct{}, 0)
	for _, r := range rewards {
		fmt.Println("index", r.Index, "BlockNumber", r.BlockNumber, "ValidatorIndex", r.ValidatorIndex, "ValidatorPubKey", r.ValidatorPubKey, "value", r.Reward.String())
		totalRewards = big.NewInt(0).Add(totalRewards, r.Reward)
		validators[r.ValidatorIndex] = struct{}{}
	}

	fmt.Println("validator counts: ", len(validators))
	fmt.Printf("total Rewards: %s Gwei, %s Wei \n", totalRewards.String(), big.NewInt(0).Mul(totalRewards, big.NewInt(1000000000)).String())
}
