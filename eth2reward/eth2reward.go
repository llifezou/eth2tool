package eth2reward

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var scanUrl = "https://etherscan.io/txsBeaconWithdrawal?ps=100&cn=%s&m=&sort=withdrawalIndex&order=desc&p=%d"

type RewardMeta struct {
	Index           string
	BlockNumber     string
	ValidatorIndex  string
	ValidatorPubKey string
	Reward          *big.Int
	Balance         *big.Int
}

func QueryEth2Rewards(recipientAddr string, eth2Rpc string) ([]RewardMeta, error) {
	page := 1

	rewards := make([]RewardMeta, 0)
	pubkeyCache := make(map[string]string)
	balanceCache := make(map[string]*big.Int)
	for {
		response, err := http.Get(fmt.Sprintf(scanUrl, recipientAddr, page))
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		response.Body.Close()

		html := string(body)
		rows := regexp.MustCompile(`<tr>(.*?)</tr>`).FindAllString(html, -1)

		if len(rows) == 0 {
			break
		}

		for _, row := range rows {
			cols := regexp.MustCompile(`<td(.*?)>(.*?)</td>`).FindAllStringSubmatch(row, -1)

			if len(cols) == 8 {
				index := regexp.MustCompile(`<td>(.*?)</td>`).FindStringSubmatch(cols[0][0])[1]
				block := regexp.MustCompile(`'/block/(.*?)'>`).FindStringSubmatch(cols[1][0])[1]
				validatorIndex := regexp.MustCompile(`_blank'>(.*?)</a>`).FindStringSubmatch(cols[5][0])[1]
				value := regexp.MustCompile(`<td>(.*?)<b>`).FindStringSubmatch(cols[7][0])[1] + "." + regexp.MustCompile(`</b>(.*?) ETH</td>`).FindStringSubmatch(cols[7][0])[1]

				reward, err := strconv.ParseFloat(value, 64)
				bigReward := big.NewFloat(reward * 1000000000)
				if err != nil {
					panic(err)
				}
				bigRewardStr := bigReward.String()
				bigIntReward, isOk := big.NewInt(0).SetString(bigRewardStr, 10)
				if !isOk {
					panic(isOk)
				}

				_, ok := pubkeyCache[validatorIndex]
				if eth2Rpc != "" && !ok {
					v, err := GetValidatorInfoOfIndex(validatorIndex, eth2Rpc)
					if err == nil {
						fmt.Println("validator index: ", validatorIndex, " pubkey:", v.Data.Validator.Pubkey)
					} else {
						fmt.Println("GetValidatorInfoOfIndex fail, err: ", err.Error())
					}
					pubkeyCache[validatorIndex] = v.Data.Validator.Pubkey
					balance, isOk := big.NewInt(0).SetString(v.Data.Balance, 10)
					if !isOk {
						panic(isOk)
					}
					balanceCache[validatorIndex] = balance
				}

				rewards = append(rewards, RewardMeta{
					Index:           index,
					BlockNumber:     block,
					ValidatorIndex:  validatorIndex,
					ValidatorPubKey: pubkeyCache[validatorIndex],
					Reward:          bigIntReward,
					Balance:         balanceCache[validatorIndex],
				})
			}
		}

		page += 1
		time.Sleep(1 * time.Second)
	}

	return rewards, nil
}

type ValidatorInfo struct {
	Data struct {
		Index     string `json:"index"`
		Balance   string `json:"balance"`
		Status    string `json:"status"`
		Validator struct {
			Pubkey                     string `json:"pubkey"`
			WithdrawalCredentials      string `json:"withdrawal_credentials"`
			EffectiveBalance           string `json:"effective_balance"`
			Slashed                    bool   `json:"slashed"`
			ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
			ActivationEpoch            string `json:"activation_epoch"`
			ExitEpoch                  string `json:"exit_epoch"`
			WithdrawableEpoch          string `json:"withdrawable_epoch"`
		} `json:"validator"`
	} `json:"data"`
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
}

func GetValidatorInfoOfIndex(index string, eth2Rpc string) (*ValidatorInfo, error) {
	response, err := http.Get(fmt.Sprintf(eth2Rpc+"/eth/v1/beacon/states/head/validators/%s", index))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()

	var vInfo ValidatorInfo
	err = json.Unmarshal(body, &vInfo)
	if err != nil {
		return nil, err
	}

	return &vInfo, nil
}
