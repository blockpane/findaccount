# findaccount

Tool for identifying what IBC chains an account exists on. Give it an address and it will query public RPC nodes and print out CSV with info.

**Important** it is necessary to retrieve the JSON files from the [Cosmos chain-registry](https://github.com/cosmos/chain-registry). To do this, run the [chains/fetch.sh](chains/fetch.sh) script.

Build: 

```
bash ./chains/fetch.sh
go get ./...
go build -o findaccounts cmd/findaccounts/main.go
```

![example](example.png)

Example:

```
findaccounts juno1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twfn0ja8 |grep true
cerberus,cerberus1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twrxvq0s,"",true,balance:<denom:"ucrbrus" amount:"514436665011420" > ,ok
chihuahua,chihuahua1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twu5p8me,"",true,balance:<denom:"uhuahua" amount:"15375994400" > ,ok
comdex,comdex1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twcwwtrv,"",true,balance:<denom:"ucmdx" amount:"300000000" > ,ok
cosmoshub,cosmos1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twlpvf6m,"",true,balance:<denom:"uatom" amount:"37256755969" > ,ok
dig,dig1aeh8gqu9wr4u8ev6edlgfq03rcy6v5tw849zcq,"",true,balance:<denom:"udig" amount:"116934" > ,ok
evmos,evmos1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twaqa8qn,"",true,balance:<denom:"ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518" amount:"5000" > ,ok
galaxy,galaxy1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twr72f3f,"",true,balance:<denom:"uglx" amount:"660000000000" > ,ok
gravitybridge,gravity1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twm373ln,"",true,balance:<denom:"ugraviton" amount:"4287" > ,ok
juno,juno1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twfn0ja8,"",true,balance:<denom:"ibc/008BFD000A10BCE5F0D4DD819AE1C1EC2942396062DABDD6AE64A655ABC7085B" amount:"686021124" > ,ok
kichain,ki1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twwvax70,"",true,balance:<denom:"uxki" amount:"6586450747" > ,ok
likecoin,like1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twvasteq,"",true,balance:<denom:"nanolike" amount:"4990540034853" > ,ok
meme,meme1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twp767a3,"",true,balance:<denom:"umeme" amount:"191311162413" > ,ok
osmosis,osmo1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twh6levf,"",true,balance:<denom:"uosmo" amount:"119849309021" > ,ok
stargaze,stars1aeh8gqu9wr4u8ev6edlgfq03rcy6v5twtam532,"",true,balance:<denom:"ustars" amount:"493715660" > ,ok
```