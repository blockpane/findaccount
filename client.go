package findaccount

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type ChainInfo struct {
	Apis struct {
		Rpc []Rpc `json:"rpc"`
	} `json:"apis"`
}

type Rpc struct {
	Address string `json:"address"`
}

func QueryAccount(info *ChainInfo, chain, account string) (hasBalance bool, balances string, err error) {

	client := &rpchttp.HTTP{}
	//defer client.Stop()
	ok := false
	for _, endpoint := range info.Apis.Rpc {
		client, err = rpchttp.NewWithTimeout(endpoint.Address, "/websocket", 10)
		if err != nil {
			continue
		}
		_, err = client.Status(context.Background())
		if err != nil {
			//_ = client.Stop()
			continue
		}
		ok = true
		break
	}
	if !ok {
		return false, "", fmt.Errorf("could not connect to any endpoints for %s", chain)
	}
	q := types.QueryBalanceRequest{Address: account}
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
		balResp := types.QueryBalanceResponse{}
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
