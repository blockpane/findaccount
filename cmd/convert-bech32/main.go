package main

import (
	"fmt"
	findaccount "github.com/blockpane/comosaccounts"
	"log"
	"os"
	"sort"
)

func main() {
	if len(os.Args) == 2 {
		accounts, err := findaccount.ConvertToAccounts(os.Args[1])
		if err != nil {
			log.Println(err)
		}
		results := make([]string, len(accounts))
		i := 0
		for k, v := range accounts {
			results[i] = fmt.Sprintf("%-14s: %s\n", k, v)
			i += 1
		}
		sort.Strings(results)
		for _, s := range results {
			fmt.Print(s)
		}
	} else {
		log.Fatalf("Error %s takes one argument, a bech32 encoded cosmos address\n", os.Args[0])
	}
}
