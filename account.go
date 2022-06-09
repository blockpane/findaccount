package findaccount

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"log"
	"sort"
	"sync"
)

var accountsMux sync.Mutex

type Result struct {
	Chain      string `json:"chain"`
	Address    string `json:"address"`
	HasBalance bool   `json:"hasBalance"`
	Coins      string `json:"coins"`
	Error      string `json:"error"`
}

func (r Result) CsvHeader() string {
	return "chain,address,has balance,coins,error"
}

func (r Result) ToCsv() string {
	return fmt.Sprintf("%s,%s,%v,%s,%s", r.Chain, r.Address, r.HasBalance, r.Coins, r.Error)
}

func SearchAccounts(account string) ([]Result, error) {
	results := make([]Result, 0)

	addrMap, err := ConvertToAccounts(account)
	if err != nil {
		return results, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(infos))
	for k, v := range infos {

		accountsMux.Lock()
		chain, rpcs := k, v
		addr := addrMap[k]
		accountsMux.Unlock()

		go func() {
			bal, coins, e := QueryAccount(rpcs, chain, addr)
			errStr := "ok"
			if e != nil {
				errStr = e.Error()
			}
			results = append(results, Result{
				Chain:      chain,
				Address:    addr,
				HasBalance: bal,
				Coins:      coins,
				Error:      errStr,
			})
			wg.Done()
		}()
	}
	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return sort.StringsAreSorted([]string{results[i].Chain, results[j].Chain})
	})

	return results, err
}

func ConvertToAccounts(s string) (map[string]string, error) {
	accounts := make(map[string]string)
	_, b64, err := bech32.DecodeAndConvert(s)

	if err != nil {
		return nil, err
	}

	for k, v := range Prefixes {
		addr, e := bech32.ConvertAndEncode(v, b64)
		if e != nil {
			log.Println(k, e)
		}
		accounts[k] = addr
	}

	return accounts, nil
}
