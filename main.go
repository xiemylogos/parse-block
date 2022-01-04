package parse_block

import (
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
	vconfig "github.com/ontio/ontology/consensus/vbft/config"
)

func main() {
	ontSdk := ontology_go_sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress("rpc")
	height,err := ontSdk.GetCurrentBlockHeight()
	if err != nil {
		panic(err)
	}
	block,err := ontSdk.GetBlockByHeight(height)
	if err != nil {
		panic(err)
	}
	usedPubKey := make(map[string]bool)
	for _, bookkeeper := range block.Header.Bookkeepers {
		pubkey := vconfig.PubkeyID(bookkeeper)
		if usedPubKey[pubkey]{
			log.Errorf("duplicate pubkey:%s,height:%d",pubkey,block.Header.Height)
		}
		usedPubKey[pubkey] = true
	}
	log.Info("usedPubKey:%d,height:%d",len(usedPubKey),block.Header.Height)
	return
}