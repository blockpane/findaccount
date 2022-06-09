package main

import (
	"encoding/json"
	"flag"
	"fmt"
	findaccount "github.com/blockpane/comosaccounts"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"io/fs"
	"log"
	"net/http"
	"net/netip"
	"sort"
)

func main() {
	var port int
	var xForwarded string
	var useXForwarded bool

	flag.IntVar(&port, "p", 8080, "http port to listen on")
	flag.StringVar(&xForwarded, "h", "X-Forwarded-For", "optional: trusted X-Forwarded-For Header")
	flag.BoolVar(&useXForwarded, "x", false, "Use the X-Forwarded-For header for logs (behind a reverse proxy)")
	flag.Parse()

	invalidRequest := []byte(`{"error":"invalid request"}`)
	invalidResponse := []byte(`"error":"unknown server error"`)

	http.HandleFunc("/q", func(writer http.ResponseWriter, request *http.Request) {
		remoteIp := request.RemoteAddr

		log := func(msg string) {
			log.Printf("%s: %s", request.RemoteAddr, msg)
		}

		remoteIp = request.Header.Get(xForwarded)
		if useXForwarded {
			_, e := netip.ParseAddr(remoteIp)
			if e != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				_, _ = writer.Write(invalidResponse)
				log(fmt.Sprintf("invalid value in %s header: %q", xForwarded, request.Header.Get(xForwarded)))
				return
			}
		}

		addr := request.URL.Query()["addr"]
		if addr == nil || len(addr) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(invalidRequest)
			return
		}

		// ensure a valid addr before continuing
		_, _, err := bech32.DecodeAndConvert(addr[0])
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(invalidRequest)
			log(fmt.Sprintf("could not decode bech32 address %q", addr[0]))
			return
		}

		result, err := findaccount.SearchAccounts(addr[0])
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write(invalidResponse)
			log("server error: " + err.Error())
			return
		}

		// re-sort so chains with a balance are on the top of the list
		sort.Slice(result, func(i, j int) bool {
			return result[i].HasBalance == true && result[j].HasBalance == false
		})

		body, err := json.Marshal(result)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write(invalidResponse)
			log("could not serialize results: " + err.Error())
			return
		}

		_, _ = writer.Write(body)
	})
	rootDir, err := fs.Sub(findaccount.StaticFs, "static")
	if err != nil {
		log.Fatalln(err)
	}
	http.Handle("/", http.FileServer(http.FS(rootDir)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
