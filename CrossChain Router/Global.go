package main

var (
	ChainId           string
	ChainPrivateKey   []byte
	ChainAddress      string
	ChainPort         string
	ChainType         string
	ChainCalcResoure  string
	RelayChainAddress string
	localPort         string
	LocalKey          []byte
	Keys              map[string][]byte // only used in relay chain
	TmpUUID           string
	UUIDList          map[string]struct{} //  only used in relay chain
)
