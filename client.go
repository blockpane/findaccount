package findaccount

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"regexp"
	"strings"
)

type ChainInfo struct {
	Apis struct {
		Rpc []Rpc `json:"rpc"`
	} `json:"apis"`
}

type Rpc struct {
	Address string `json:"address"`
}

var portRex = regexp.MustCompile(`.*:\d+$`)
var protoRex = regexp.MustCompile(`^\w+://`)

func getClient(info *ChainInfo, chain string) (*rpchttp.HTTP, error) {
	client := &rpchttp.HTTP{}
	var err error
	ok := false
	for i := range info.Apis.Rpc {
		endpoint := info.Apis.Rpc[len(info.Apis.Rpc)-1-i]
		endpoint.Address = strings.TrimRight(endpoint.Address, "/")
		var unknown bool

		if !portRex.MatchString(endpoint.Address) {
			switch protoRex.FindString(endpoint.Address) {
			case "https://":
				endpoint.Address = endpoint.Address + ":443"
			case "http://":
				endpoint.Address = endpoint.Address + ":80"
			case "tcp://":
				endpoint.Address = endpoint.Address + ":26657"
			default:
				unknown = true
			}
		}
		if unknown {
			continue
		}
		client, err = rpchttp.NewWithTimeout(endpoint.Address, "/websocket", 10)
		if err != nil {
			continue
		}
		status, e := client.Status(context.Background())
		if e != nil || status.SyncInfo.CatchingUp {
			continue
		}
		ok = true
		break
	}
	if !ok {
		err = fmt.Errorf("could not connect to any endpoints for %s", chain)
	}
	return client, err
}

func IsValidator(info *ChainInfo, chain, account string) (validator string, err error) {
	client, err := getClient(info, chain)
	if err != nil {
		return
	}
	// Check if the account is also a validator
	_, b64, err := bech32.DecodeAndConvert(account)
	if err != nil {
		return
	}
	accountsMux.Lock()
	prefix := Prefixes[chain]
	accountsMux.Unlock()
	addr, _ := bech32.ConvertAndEncode(prefix+"valoper", b64)
	valQ := staketypes.QueryValidatorRequest{ValidatorAddr: addr}
	valQuery, err := valQ.Marshal()
	if err != nil {
		return
	}
	valResult, err := client.ABCIQuery(context.Background(), "/cosmos.staking.v1beta1.Query/Validator", valQuery)
	if err != nil {
		return
	}
	if len(valResult.Response.Value) > 0 {
		valResp := staketypes.QueryValidatorResponse{}
		err = valResp.Unmarshal(valResult.Response.Value)
		if err != nil {
			return
		}
		validator = valResp.Validator.GetMoniker()
		//fmt.Println(valResp)

	}
	return
}

func QueryAccount(info *ChainInfo, chain, account string) (hasBalance bool, balances string, err error) {

	client, err := getClient(info, chain)

	if err != nil {
		return false, "", err
	}
	q := banktypes.QueryBalanceRequest{Address: account}
	var query []byte
	query, err = q.Marshal()
	if err != nil {
		return
	}
	result, err := client.ABCIQuery(context.Background(), "/cosmos.bank.v1beta1.Query/AllBalances", query)
	if err != nil {
		return
	}

	if len(result.Response.Value) > 0 {
		balResp := banktypes.QueryBalanceResponse{}
		err = balResp.Unmarshal(result.Response.Value)
		if err != nil {
			return
		}
		balances = balResp.String()
		if len(balances) > 0 {
			hasBalance = true
		}
	}

	return
}
