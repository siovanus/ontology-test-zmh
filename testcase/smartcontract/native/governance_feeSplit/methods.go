package governance_feeSplit

import (
	"bytes"
	"encoding/hex"

	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/auth"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var OntIDVersion = byte(0)

const PROMISE_POS = 200000

func registerCandidate(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string, initPos uint32) bool {
	params := &governance.RegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
		InitPos:    initPos,
		Caller:     []byte("did:ont:" + user.Address.ToBase58()),
		KeyNo:      1,
	}
	method := "registerCandidate"
	contractAddress := utils.GovernanceContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("registerCandidate txHash is :", txHash.ToHexString())
	waitForBlock(ctx)

	//let node can be authorized
	changeMaxAuthorization(ctx, user, peerPubkey, PROMISE_POS)
	return true
}

func registerCandidate2Sign(ctx *testframework.TestFrameworkContext, ontid *sdk.Account, user *sdk.Account, peerPubkey string, initPos uint32) bool {
	params := &governance.RegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
		InitPos:    initPos,
		Caller:     []byte("did:ont:" + ontid.Address.ToBase58()),
		KeyNo:      1,
	}
	method := "registerCandidate"
	contractAddress := utils.GovernanceContractAddress
	tx, err := ctx.Ont.Native.NewNativeInvokeTransaction(ctx.GetGasPrice(), ctx.GetGasLimit(), OntIDVersion, contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("NewNativeInvokeTransaction error")
		return false
	}
	err = ctx.Ont.SignToTransaction(tx, user)
	if err != nil {
		ctx.LogError("SignToTransaction error")
		return false
	}
	err = ctx.Ont.SignToTransaction(tx, ontid)
	if err != nil {
		ctx.LogError("SignToTransaction error")
		return false
	}
	txHash, err := ctx.Ont.SendTransaction(tx)
	if err != nil {
		ctx.LogError("SendRawTransaction error", err)
		return false
	}
	ctx.LogInfo("registerCandidate2Sign txHash is :", txHash.ToHexString())
	return true
}

func unRegisterCandidate(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string) bool {
	params := &governance.UnRegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
	}
	method := "unRegisterCandidate"
	contractAddress := utils.GovernanceContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("unRegisterCandidate txHash is :", txHash.ToHexString())
	return true
}

func approveCandidate(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string) bool {
	params := &governance.ApproveCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "approveCandidate"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("approveCandidate txHash is :", txHash.ToHexString())
	return true
}

func approveCandidateMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string) bool {
	params := &governance.ApproveCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "approveCandidate"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("approveCandidateMultiSign txHash is :", txHash.ToHexString())
	return true
}

func rejectCandidate(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string) bool {
	params := &governance.RejectCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "rejectCandidate"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("rejectCandidate txHash is :", txHash.ToHexString())
	return true
}

func rejectCandidateMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string) bool {
	params := &governance.RejectCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "rejectCandidate"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("rejectCandidateMultiSign txHash is :", txHash.ToHexString())
	return true
}

func changeMaxAuthorization(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string, maxAuthorize uint32) bool {
	params := &governance.ChangeMaxAuthorizationParam{
		Address:      user.Address,
		PeerPubkey:   peerPubkey,
		MaxAuthorize: maxAuthorize,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "changeMaxAuthorization"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("changeMaxAuthorization txHash is :", txHash.ToHexString())
	return true
}

func setPeerCost(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string, peerCost uint32) bool {
	params := &governance.SetPeerCostParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		PeerCost:   peerCost,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "setPeerCost"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("setPeerCost txHash is :", txHash.ToHexString())
	return true
}

func addInitPos(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string, pos uint32) bool {
	params := &governance.ChangeInitPosParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		Pos:        pos,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "addInitPos"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("addInitPos txHash is :", txHash.ToHexString())
	return true
}

func reduceInitPos(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string, pos uint32) bool {
	params := &governance.ChangeInitPosParam{
		Address:    user.Address,
		PeerPubkey: peerPubkey,
		Pos:        pos,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "reduceInitPos"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("reduceInitPos txHash is :", txHash.ToHexString())
	return true
}

func authorizeForPeer(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkeyList []string, posList []uint32) bool {
	params := &governance.AuthorizeForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "authorizeForPeer"
	_, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	//ctx.LogInfo("authorizeForPeer txHash is :", txHash.ToHexString())
	return true
}

func unAuthorizeForPeer(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkeyList []string, posList []uint32) bool {
	params := &governance.AuthorizeForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "unAuthorizeForPeer"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("unAuthorizeForPeer txHash is :", txHash.ToHexString())
	return true
}

func withdraw(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkeyList []string, withdrawList []uint32) bool {
	params := &governance.WithdrawParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		WithdrawList:   withdrawList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdraw"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("withdraw txHash is :", txHash.ToHexString())
	return true
}

func withdrawOng(ctx *testframework.TestFrameworkContext, user *sdk.Account) bool {
	params := &governance.WithdrawOngParam{
		Address: user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdrawOng"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("withdrawOng txHash is :", txHash.ToHexString())
	return true
}

type commitDposParam struct {
}

func commitDpos(ctx *testframework.TestFrameworkContext, user *sdk.Account) bool {
	params := &commitDposParam{}
	contractAddress := utils.GovernanceContractAddress
	method := "commitDpos"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("commitDpos txHash is :", txHash.ToHexString())
	return true
}

func commitDposMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "commitDpos"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("commitDposMultiSign txHash is :", txHash.ToHexString())
	return true
}

func quitNode(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string) bool {
	params := &governance.QuitNodeParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "quitNode"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("quitNode txHash is :", txHash.ToHexString())
	return true
}

func blackNode(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkeyList []string) bool {
	params := &governance.BlackNodeParam{
		PeerPubkeyList: peerPubkeyList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "blackNode"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("blackNode txHash is :", txHash.ToHexString())
	return true
}

func blackNodeMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkeyList []string) bool {
	params := &governance.BlackNodeParam{
		PeerPubkeyList: peerPubkeyList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "blackNode"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("blackNodeMultiSign txHash is :", txHash.ToHexString())
	return true
}

func whiteNode(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string) bool {
	params := &governance.WhiteNodeParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "whiteNode"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("whiteNode txHash is :", txHash.ToHexString())
	return true
}

func whiteNodeMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string) bool {
	params := &governance.WhiteNodeParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "whiteNode"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("whiteNodeMultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateConfig(ctx *testframework.TestFrameworkContext, user *sdk.Account, config *governance.Configuration) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateConfig"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{config})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateConfig txHash is :", txHash.ToHexString())
	return true
}

func updateConfigMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, config *governance.Configuration) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateConfig"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{config})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateConfigMultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParam(ctx *testframework.TestFrameworkContext, user *sdk.Account, globalParam *governance.GlobalParam) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{globalParam})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateGlobalParam txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParamMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, globalParam *governance.GlobalParam) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{globalParam})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateGlobalParamMultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParam2(ctx *testframework.TestFrameworkContext, user *sdk.Account, globalParam2 *governance.GlobalParam2) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam2"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{globalParam2})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateGlobalParam2 txHash is :", txHash.ToHexString())
	return true
}

func updateGlobalParam2MultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, globalParam2 *governance.GlobalParam2) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam2"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{globalParam2})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateGlobalParam2MultiSign txHash is :", txHash.ToHexString())
	return true
}

func updateSplitCurve(ctx *testframework.TestFrameworkContext, user *sdk.Account, splitCurve *governance.SplitCurve) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateSplitCurve"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{splitCurve})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateSplitCurve txHash is :", txHash.ToHexString())
	return true
}

func updateSplitCurveMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, splitCurve *governance.SplitCurve) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateSplitCurve"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{splitCurve})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("updateSplitCurveMultiSign txHash is :", txHash.ToHexString())
	return true
}

func setPromisePos(ctx *testframework.TestFrameworkContext, user *sdk.Account, promisePos *governance.PromisePos) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "setPromisePos"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{promisePos})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("setPromisePos txHash is :", txHash.ToHexString())
	return true
}

func setPromisePosMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, promisePos *governance.PromisePos) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "setPromisePos"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{promisePos})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("setPromisePosMultiSign txHash is :", txHash.ToHexString())
	return true
}

func setGasAddressMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address string) bool {
	addr, err := common.AddressFromBase58(address)
	if err != nil {
		ctx.LogError("common.AddressFromBase58 error :", err)
		return false
	}
	params := &governance.GasAddress{
		Address: addr,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "setGasAddress"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("setGasAddressMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferPenalty(ctx *testframework.TestFrameworkContext, user *sdk.Account, peerPubkey string, address common.Address) bool {
	params := &governance.TransferPenaltyParam{
		PeerPubkey: peerPubkey,
		Address:    address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "transferPenalty"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferPenalty txHash is :", txHash.ToHexString())
	return true
}

func withdrawFee(ctx *testframework.TestFrameworkContext, user *sdk.Account) bool {
	params := &governance.WithdrawFeeParam{
		Address: user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdrawFee"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("withdrawFee txHash is :", txHash.ToHexString())
	return true
}

func transferPenaltyMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, peerPubkey string, address common.Address) bool {
	params := &governance.TransferPenaltyParam{
		PeerPubkey: peerPubkey,
		Address:    address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "transferPenalty"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferPenaltyMultiSign txHash is :", txHash.ToHexString())
	return true
}

func multiTransfer(ctx *testframework.TestFrameworkContext, contract common.Address, from []*sdk.Account, to []string, amount []uint64) bool {
	var sts []ont.State
	if len(from) != len(to) || len(from) != len(amount) {
		ctx.LogError("input length error")
		return false
	}
	for i := 0; i < len(from); i++ {
		address, err := common.AddressFromBase58(to[i])
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
			return false
		}
		sts = append(sts, ont.State{
			From:  from[i].Address,
			To:    address,
			Value: amount[i],
		})
	}
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := contract
	method := "transfer"
	tx, err := ctx.Ont.Native.NewNativeInvokeTransaction(ctx.GetGasPrice(), ctx.GetGasLimit(), OntIDVersion, contractAddress, method, []interface{}{transfers})
	if err != nil {
		return false
	}
	for _, singer := range from {
		err = ctx.Ont.SignToTransaction(tx, singer)
		if err != nil {
			return false
		}
	}
	txHash, err := ctx.Ont.SendTransaction(tx)
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
		return false
	}
	ctx.LogInfo("multiTransfer txHash is :", txHash.ToHexString())
	return true
}

func transferOntMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OntContractAddress
	method := "transfer"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferOntMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferOntMultiSignToMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OntContractAddress
	method := "transfer"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferOntMultiSignToMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferOngMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OngContractAddress
	method := "transfer"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferOngMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferOngMultiSignToMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	var sts []ont.State
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	sts = append(sts, ont.State{
		From:  from,
		To:    address,
		Value: amount,
	})
	transfers := ont.Transfers{
		States: sts,
	}
	contractAddress := utils.OngContractAddress
	method := "transfer"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{transfers})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferOngMultiSignToMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferFromOngMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	params := &ont.TransferFrom{
		Sender: from,
		From:   utils.OntContractAddress,
		To:     address,
		Value:  amount,
	}
	contractAddress := utils.OngContractAddress
	method := "transferFrom"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferFromOngMultiSign txHash is :", txHash.ToHexString())
	return true
}

func transferFromOngMultiSignToMultiSign(ctx *testframework.TestFrameworkContext, pubKeys []keypair.PublicKey, users []*sdk.Account, address common.Address, amount uint64) bool {
	from, err := types.AddressFromMultiPubKeys(pubKeys, int((5*len(pubKeys)+6)/7))
	if err != nil {
		ctx.LogError("types.AddressFromMultiPubKeys error", err)
	}
	params := &ont.TransferFrom{
		Sender: from,
		From:   utils.OntContractAddress,
		To:     address,
		Value:  amount,
	}
	contractAddress := utils.OngContractAddress
	method := "transferFrom"
	txHash, err := invokeNativeContractWithMultiSign(ctx, ctx.GetGasPrice(), ctx.GetGasLimit(), pubKeys, users, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("transferFromOngMultiSignToMultiSign txHash is :", txHash.ToHexString())
	return true
}

func assignFuncsToRole(ctx *testframework.TestFrameworkContext, user *sdk.Account, contract common.Address, role string, function string) bool {
	params := &auth.FuncsToRoleParam{
		ContractAddr: contract,
		AdminOntID:   []byte("did:ont:" + user.Address.ToBase58()),
		Role:         []byte(role),
		FuncNames:    []string{function},
		KeyNo:        1,
	}
	method := "assignFuncsToRole"
	contractAddress := utils.AuthContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("assignFuncsToRole txHash is :", txHash.ToHexString())
	return true
}

func assignOntIDsToRole(ctx *testframework.TestFrameworkContext, user *sdk.Account, contract common.Address, role string, ontids []string) bool {
	params := &auth.OntIDsToRoleParam{
		ContractAddr: contract,
		AdminOntID:   []byte("did:ont:" + user.Address.ToBase58()),
		Role:         []byte(role),
		Persons:      [][]byte{},
		KeyNo:        1,
	}
	for _, ontid := range ontids {
		params.Persons = append(params.Persons, []byte(ontid))
	}
	contractAddress := utils.AuthContractAddress
	method := "assignOntIDsToRole"
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("assignOntIDsToRole txHash is :", txHash.ToHexString())
	return true
}

type RegIDWithPublicKeyParam struct {
	OntID  []byte
	Pubkey []byte
}

func regIdWithPublicKey(ctx *testframework.TestFrameworkContext, user *sdk.Account) bool {
	params := RegIDWithPublicKeyParam{
		OntID:  []byte("did:ont:" + user.Address.ToBase58()),
		Pubkey: keypair.SerializePublicKey(user.PublicKey),
	}
	method := "regIDWithPublicKey"
	contractAddress := utils.OntIDContractAddress
	txHash, err := ctx.Ont.Native.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error :", err)
		return false
	}
	ctx.LogInfo("RegIDWithPublicKeyParam txHash is :", txHash.ToHexString())
	return true
}

func getVbftConfig(ctx *testframework.TestFrameworkContext) (*governance.Configuration, error) {
	contractAddress := utils.GovernanceContractAddress
	config := new(governance.Configuration)
	key := []byte(governance.VBFT_CONFIG)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := config.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize config error!")
	}
	return config, nil
}

func getPreConfig(ctx *testframework.TestFrameworkContext) (*governance.PreConfig, error) {
	contractAddress := utils.GovernanceContractAddress
	preConfig := new(governance.PreConfig)
	key := []byte(governance.PRE_CONFIG)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := preConfig.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize preConfig error!")
	}
	return preConfig, nil
}

func getGlobalParam(ctx *testframework.TestFrameworkContext) (*governance.GlobalParam, error) {
	contractAddress := utils.GovernanceContractAddress
	globalParam := new(governance.GlobalParam)
	key := []byte(governance.GLOBAL_PARAM)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := globalParam.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize globalParam error!")
	}
	return globalParam, nil
}

func getGlobalParam2(ctx *testframework.TestFrameworkContext) (*governance.GlobalParam2, error) {
	contractAddress := utils.GovernanceContractAddress
	globalParam2 := new(governance.GlobalParam2)
	key := []byte(governance.GLOBAL_PARAM2)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		if err := globalParam2.Deserialization(common.NewZeroCopySource(value)); err != nil {
			return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize globalParam2 error!")
		}
	}
	return globalParam2, nil
}

func getSplitCurve(ctx *testframework.TestFrameworkContext) (*governance.SplitCurve, error) {
	contractAddress := utils.GovernanceContractAddress
	splitCurve := new(governance.SplitCurve)
	key := []byte(governance.SPLIT_CURVE)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := splitCurve.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize splitCurve error!")
	}
	return splitCurve, nil
}

func getGovernanceView(ctx *testframework.TestFrameworkContext) (*governance.GovernanceView, error) {
	contractAddress := utils.GovernanceContractAddress
	governanceView := new(governance.GovernanceView)
	key := []byte(governance.GOVERNANCE_VIEW)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := governanceView.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize governanceView error!")
	}
	return governanceView, nil
}

func getView(ctx *testframework.TestFrameworkContext) (uint32, error) {
	governanceView, err := getGovernanceView(ctx)
	if err != nil {
		return 0, errors.NewDetailErr(err, errors.ErrNoCode, "getGovernanceView error")
	}
	return governanceView.View, nil
}

func getPeerPoolMap(ctx *testframework.TestFrameworkContext) (*governance.PeerPoolMap, error) {
	contractAddress := utils.GovernanceContractAddress
	view, err := getView(ctx)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getView error")
	}
	peerPoolMap := &governance.PeerPoolMap{
		PeerPoolMap: make(map[string]*governance.PeerPoolItem),
	}
	viewBytes, err := governance.GetUint32Bytes(view)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "GetUint32Bytes, get viewBytes error!")
	}
	key := ConcatKey([]byte(governance.PEER_POOL), viewBytes)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := peerPoolMap.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize peerPoolMap error!")
	}
	return peerPoolMap, nil
}

func getAuthorizeInfo(ctx *testframework.TestFrameworkContext, peerPubkey string, address common.Address) (*governance.AuthorizeInfo, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	authorizeInfo := new(governance.AuthorizeInfo)
	key := ConcatKey([]byte(governance.AUTHORIZE_INFO_POOL), peerPubkeyPrefix, address[:])
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := authorizeInfo.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize authorizeInfo error!")
	}
	return authorizeInfo, nil
}

func inBlackList(ctx *testframework.TestFrameworkContext, peerPubkey string) (bool, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	key := ConcatKey([]byte(governance.BLACK_LIST), peerPubkeyPrefix)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		return true, nil
	}
	return false, nil
}

func getTotalStake(ctx *testframework.TestFrameworkContext, address common.Address) (*governance.TotalStake, error) {
	contractAddress := utils.GovernanceContractAddress
	totalStake := new(governance.TotalStake)
	key := ConcatKey([]byte(governance.TOTAL_STAKE), address[:])
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := totalStake.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize totalStake error!")
	}
	return totalStake, nil
}

func getPenaltyStake(ctx *testframework.TestFrameworkContext, peerPubkey string) (*governance.PenaltyStake, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	penaltyStake := new(governance.PenaltyStake)
	key := ConcatKey([]byte(governance.PENALTY_STAKE), peerPubkeyPrefix)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := penaltyStake.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize penaltyStake error!")
	}
	return penaltyStake, nil
}

func getAttributes(ctx *testframework.TestFrameworkContext, peerPubkey string) (*governance.PeerAttributes, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	peerAttributes := new(governance.PeerAttributes)
	key := ConcatKey([]byte(governance.PEER_ATTRIBUTES), peerPubkeyPrefix)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		if err := peerAttributes.Deserialization(common.NewZeroCopySource(value)); err != nil {
			return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize peerAttributes error!")
		}
	}
	return peerAttributes, nil
}

func getSplitFeeAddress(ctx *testframework.TestFrameworkContext, address common.Address) (*governance.SplitFeeAddress, error) {
	contractAddress := utils.GovernanceContractAddress
	splitFeeAddress := new(governance.SplitFeeAddress)
	key := ConcatKey([]byte(governance.SPLIT_FEE_ADDRESS), address[:])
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := splitFeeAddress.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize splitFeeAddress error!")
	}
	return splitFeeAddress, nil
}

func getSplitFee(ctx *testframework.TestFrameworkContext) (uint64, error) {
	contractAddress := utils.GovernanceContractAddress
	key := ConcatKey([]byte(governance.SPLIT_FEE))
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return 0, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	splitFee, err := serialization.ReadUint64(bytes.NewBuffer(value))
	if err != nil {
		return 0, errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadUint64, deserialize splitFee error!")
	}
	return splitFee, nil
}

func getPromisePos(ctx *testframework.TestFrameworkContext, peerPubkey string) (*governance.PromisePos, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	key := ConcatKey([]byte(governance.PROMISE_POS), peerPubkeyPrefix)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	promisePos := new(governance.PromisePos)
	if err := promisePos.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize promisePos error!")
	}
	return promisePos, nil
}

func getGasAddress(ctx *testframework.TestFrameworkContext) (*governance.GasAddress, error) {
	contractAddress := utils.GovernanceContractAddress
	key := []byte(governance.GAS_ADDRESS)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	gasAddress := new(governance.GasAddress)
	if err := gasAddress.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize gasAddress error!")
	}
	return gasAddress, nil
}