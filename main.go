package main

import (
	"fmt"
	"github.com/llifezou/eth2tool/eth2reward"
	"github.com/spf13/cobra"
	"math/big"
)

func main() {
	_ = eth2RewardsCmd.Execute()
}

var (
	recipientAddr string
	eth2Rpc       string

	eth2RewardsCmd = &cobra.Command{
		Use:   "eth2-rewards",
		Short: "eth2 tool",
		Long:  "query eth2 rewards tool",
		Run: func(cmd *cobra.Command, args []string) {
			rewards, err := eth2reward.QueryEth2Rewards(recipientAddr, eth2Rpc)
			if err != nil {
				fmt.Println(err)
				return
			}
			totalRewards := big.NewInt(0)
			validators := make(map[string]struct{}, 0)
			for _, r := range rewards {
				fmt.Println("index", r.Index, " BlockNumber", r.BlockNumber, " ValidatorIndex", r.ValidatorIndex, " ValidatorPubKey", r.ValidatorPubKey, " value", r.Reward.String())
				totalRewards = big.NewInt(0).Add(totalRewards, r.Reward)
				validators[r.ValidatorIndex] = struct{}{}
			}

			fmt.Println("validator counts: ", len(validators))
			fmt.Printf("total Rewards: %s Gwei, %s Wei \n", totalRewards.String(), big.NewInt(0).Mul(totalRewards, big.NewInt(1000000000)).String())
			return
		},
	}
)

func init() {
	eth2RewardsCmd.Flags().StringVar(&recipientAddr, "recipientAddr", "", "eth2 recipient addr (execution layer address)")
	eth2RewardsCmd.Flags().StringVar(&eth2Rpc, "eth2Rpc", "", "If present, the pubkey will be queried")
}
