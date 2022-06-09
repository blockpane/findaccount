package findaccount

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"os"
)

var (
	//go:embed chains/*.json
	chainsFs embed.FS
	infos    = make(map[string]*ChainInfo)
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lshortfile)

	// parse out the embedded json files for RPC endpoints
	for k := range Prefixes {
		f, e := chainsFs.Open("chains/" + k + "-chain.json")
		if e != nil {
			log.Println(e)
			continue
		}
		b, e := io.ReadAll(f)
		if e != nil {
			log.Println(e)
			continue
		}
		_ = f.Close()
		chainInfo := &ChainInfo{}
		e = json.Unmarshal(b, chainInfo)
		if e != nil {
			log.Println(e)
			continue
		}
		if chainInfo != nil && len(chainInfo.Apis.Rpc) > 0 {
			infos[k] = chainInfo
		}
	}

	// add extra known-good RPC servers....
	for k, v := range additional {
		if infos[k] == nil {
			log.Println(k, "is not defined skipping addition of RPC")
			continue
		}
		if infos[k].Apis.Rpc == nil {
			infos[k].Apis.Rpc = make([]Rpc, 0)
		}
		for _, node := range v {
			infos[k].Apis.Rpc = append(infos[k].Apis.Rpc, Rpc{Address: node})
		}
	}

}
