#!/bin/bash

cd $(dirname $0)

rm ./*.json

while read chain; do
  wget -O ${chain}-chain.json https://github.com/cosmos/chain-registry/raw/master/${chain}/chain.json
done << EOF
agoric
akash
arkh
assetmantle
axelar
bandchain
bitcanna
bitsong
bostrom
carbon
cerberus
cheqd
chihuahua
comdex
cosmoshub
crescent
cronos
cryptoorgchain
decentr
desmos
dig
echelon
emoney
evmos
fetchhub
firmachain
galaxy
genesisl1
gravitybridge
impacthub
injective
irisnet
juno
kava
kichain
konstellation
likecoin
logos
lumnetwork
meme
microtick
mythos
nomic
octa
odin
oraichain
osmosis
panacea
persistence
provenance
regen
rizon
secretnetwork
sentinel
shentu
sifchain
sommelier
stargaze
starname
terra
thorchain
umee
vidulum
EOF