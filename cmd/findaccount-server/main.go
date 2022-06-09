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
			//writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(invalidRequest)
			return
		}

		// ensure a valid addr before continuing
		_, _, err := bech32.DecodeAndConvert(addr[0])
		if err != nil {
			//writer.WriteHeader(http.StatusBadRequest)
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

		for i := range result {
			if result[i].Error == "ok" {
				result[i].Error = ""
			}
		}

		body, err := json.Marshal(result)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write(invalidResponse)
			log("could not serialize results: " + err.Error())
			return
		}

		_, _ = writer.Write(body)
	})

	http.Handle("/", &CacheHandler{})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// CacheHandler implements the Handler interface with a very long Cache-Control set on responses
type CacheHandler struct{}

func (ch CacheHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rootDir, err := fs.Sub(findaccount.StaticFs, "static")
	if err != nil {
		log.Fatalln(err)
	}
	writer.Header().Set("Cache-Control", "public, max-age=86400")
	http.FileServer(http.FS(rootDir)).ServeHTTP(writer, request)
}
