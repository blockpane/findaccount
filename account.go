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
	Validator  string `json:"is_validator"`
	HasBalance bool   `json:"hasBalance"`
	Coins      string `json:"coins"`
	Error      string `json:"error"`
	Link       string `json:"link"`
}

func (r Result) CsvHeader() string {
	return "chain,address,validator,has balance,coins,error"
}

func (r Result) ToCsv() string {
	return fmt.Sprintf("%s,%s,%q,%v,%s,%s", r.Chain, r.Address, r.Validator, r.HasBalance, r.Coins, r.Error)
}

// SearchAccounts is the entrypoint for performing a search
func SearchAccounts(account string) ([]Result, error) {
	results := make([]Result, 0)

	addrMap, err := ConvertToAccounts(account)
	if err != nil {
		return results, err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(infos))
	for k, v := range infos {
		var link string

		accountsMux.Lock()
		chain, rpcs := k, v
		addr := addrMap[k]
		if len(infos[k].Explorers) > 0 {
			link = infos[k].Explorers[0].Url
		}
		accountsMux.Unlock()

		go func() {
			bal, coins, e := QueryAccount(rpcs, chain, addr)
			errStr := "ok"
			if e != nil {
				errStr = e.Error()
			}
			val, _ := IsValidator(rpcs, chain, addr)
			results = append(results, Result{
				Chain:      chain,
				Address:    addr,
				Validator:  val,
				HasBalance: bal,
				Coins:      coins,
				Error:      errStr,
				Link:       link,
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
