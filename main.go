package main

import (
	"flag"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
	vconfig "github.com/ontio/ontology/consensus/vbft/config"
	"github.com/ontio/ontology/core/types"
)

var (
	ServerConfig string
)

func init() {
	flag.StringVar(&ServerConfig, "cfg", "./config.json", "Config of parse block server")
}

func main() {
	cfg, err := NewSvrConfig(ServerConfig)
	if err != nil {
		panic(err)
	}
	ontSdk := ontology_go_sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(cfg.RpcAddr)
	curentHeight := cfg.BlockHeight
	nodeIds := make(map[string]bool,0)
	for _,id := range cfg.NodeId {
		nodeIds[id] = true
	}
	for height := curentHeight; height > 0; height-- {
		if height == 0 {
			log.Info("current height:%d", height)
			return
		}
		block, err := ontSdk.GetBlockByHeight(height)
		if err != nil {
			log.Errorf("GetBlockByHeight panic height:%d", height)
			panic(err)
		}
		for _, bookkeeper := range block.Header.Bookkeepers {
			pubkey := vconfig.PubkeyID(bookkeeper)
			address := types.AddressFromPubKey(pubkey)
			if nodeIds[address.ToBase58()] {
				log.Info("block height:%d",block.Header.Height)
				panic(nil)
			}
		}
	}
}
