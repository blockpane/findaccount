package main

import (
	"fmt"
	findaccount "github.com/blockpane/comosaccounts"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		results, err := findaccount.SearchAccounts(os.Args[1])
		if err != nil {
			log.Println(err)
		}
		if len(results) > 0 {
			fmt.Println(results[0].CsvHeader())
			for _, r := range results {
				fmt.Println(r.ToCsv())
			}
		}
	} else {
		log.Fatalf("Error %s takes one argument, a bech32 encoded cosmos address\n", os.Args[0])
	}
}
