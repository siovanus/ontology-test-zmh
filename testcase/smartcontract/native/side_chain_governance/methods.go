package side_chain_governance

import (
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/smartcontract/service/native/cross_chain/header_sync"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

func syncGenesisHeader(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account,
	genesisBlockHeader []byte) bool {
	params := &header_sync.SyncGenesisHeaderParam{
		GenesisHeader: genesisBlockHeader,
	}
	contractAddress := utils.HeaderSyncContractAddress
	method := "syncGenesisHeader"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys,
		users, OntIDVersion, contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContractWithMultiSign error :", err)
	}
	ctx.LogInfo("syncGenesisHeader txHash is :", txHash.ToHexString())
	return true
}
