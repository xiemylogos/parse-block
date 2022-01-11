package main

import (
	"flag"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
	vconfig "github.com/ontio/ontology/consensus/vbft/config"
	scom "github.com/ontio/ontology/core/store/common"
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
		usedPubKey := make(map[string]bool)
		for _, bookkeeper := range block.Header.Bookkeepers {
			pubkey := vconfig.PubkeyID(bookkeeper)
			if usedPubKey[pubkey] {
				log.Errorf("duplicate pubkey:%s,height:%d", pubkey, block.Header.Height)
			}
			usedPubKey[pubkey] = true
		}
		blkInfo, err := vconfig.VbftBlock(block.Header)
		if err != nil {
			log.Errorf("VbftBlock panic height:%d", block.Header.Height)
			panic(err)
		}
		log.Infof("usedPubKey:%d,height:%d", len(usedPubKey), block.Header.Height)

		var chainConfigHeight uint32
		prevBlock, err := ontSdk.GetBlockByHeight(height - 1)
		if err != nil {
			log.Errorf("GetBlockByHeight prevHeader panic height:%d", height-1)
			panic(err)
		}
		if blkInfo.NewChainConfig != nil {
			prevBlockInfo, err := vconfig.VbftBlock(prevBlock.Header)
			if err != nil {
				log.Errorf("VbftBlock height:%d,err:%s", prevBlock.Header.Height, err)
				panic(err)
			}
			if prevBlockInfo.NewChainConfig != nil {
				chainConfigHeight = prevBlock.Header.Height
			} else {
				chainConfigHeight = prevBlockInfo.LastConfigBlockNum
			}
		} else {
			chainConfigHeight = blkInfo.LastConfigBlockNum
		}
		chainConfigBlock, err := ontSdk.GetBlockByHeight(chainConfigHeight)
		if err != nil && err != scom.ErrNotFound {
			log.Errorf("NewChainConfig is nil height:%d,err:%s", chainConfigHeight, err)
			panic(err)
		}
		if chainConfigBlock == nil {
			log.Errorf("NewChainConfig is nil height:%d", chainConfigHeight)
			panic(nil)
		}
		chanConfigBlkInfo, err := vconfig.VbftBlock(chainConfigBlock.Header)
		if err != nil {
			log.Errorf("NewChainConfig is nil height:%d", chainConfigHeight)
			panic(err)
		}
		if chanConfigBlkInfo.NewChainConfig == nil {
			log.Errorf("NewChainConfig is nil height:%d", chainConfigHeight)
			panic(nil)
		}
		c := chanConfigBlkInfo.NewChainConfig.C
		if uint32(len(usedPubKey)) < c+1 {  //test net
			//if uint32(len(usedPubKey)) < c+1 && height != 183 && height != 23610 {  //mainnet
			log.Errorf("verify header error:  height:%d,pubkey len:%d,c:%d",
				height, len(usedPubKey), c)
			panic(nil)
		}
		log.Infof("parse block succ height:%d", height)
	}
}
