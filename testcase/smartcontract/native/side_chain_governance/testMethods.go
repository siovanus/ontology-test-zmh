/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package side_chain_governance

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	sdk "github.com/ontio/multi-chain-go-sdk"
	"github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/errors"
	"github.com/ontio/multi-chain/native/service/governance/node_manager"
	"github.com/ontio/multi-chain/native/service/utils"
	osdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
)

type BlackChainParam struct {
	Path    []string
	ChainID uint64
}

func BlackChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/BlackChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	blackChainParam := new(BlackChainParam)
	err = json.Unmarshal(data, blackChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range blackChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Ccm.BlackChain(blackChainParam.ChainID, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Ccm.BlackChain error: %v", err)
		return false
	}
	ctx.LogInfo("BlackChain txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func WhiteChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/WhiteChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	blackChainParam := new(BlackChainParam)
	err = json.Unmarshal(data, blackChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range blackChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Ccm.WhiteChain(blackChainParam.ChainID, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Ccm.WhiteChain error: %v", err)
		return false
	}
	ctx.LogInfo("WhiteChain txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type SyncGenesisHeaderParam struct {
	Path     []string
	ChainID  uint64
	ChainRpc string
}

func SyncGenesisHeader(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/SyncGenesisHeader.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	syncGenesisHeaderParam := new(SyncGenesisHeaderParam)
	err = json.Unmarshal(data, syncGenesisHeaderParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range syncGenesisHeaderParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	sideSdk := osdk.NewOntologySdk()
	sideSdk.NewRpcClient().SetAddress(syncGenesisHeaderParam.ChainRpc)
	genesisBlock, err := sideSdk.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("get side chain genesis block error: %s", err)
		return false
	}
	genesisBlockHeader := genesisBlock.Header.ToArray()

	txHash, err := ctx.Ont.Native.Hs.SyncGenesisHeader(syncGenesisHeaderParam.ChainID, genesisBlockHeader, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Hs.SyncGenesisHeader error: %v", err)
		return false
	}
	ctx.LogInfo("SyncGenesisHeader txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type SideChainParam struct {
	Path         string
	Chainid      uint64
	Router       uint64
	Name         string
	BlocksToWait uint64
	CCMCAddress  string
}

func RegisterSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RegisterSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	sideChainParam := new(SideChainParam)
	err = json.Unmarshal(data, sideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, sideChainParam.Path)
	if !ok {
		return false
	}
	CCMCAddress, err := hex.DecodeString(sideChainParam.CCMCAddress)
	if err != nil {
		ctx.LogError("hex.DecodeString error %v", err)
		return false
	}
	txHash, err := ctx.Ont.Native.Scm.RegisterSideChain(user.Address, sideChainParam.Chainid,
		sideChainParam.Router, sideChainParam.Name, sideChainParam.BlocksToWait,
		CCMCAddress, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.RegisterSideChain error: %v", err)
		return false
	}
	ctx.LogInfo("RegisterSideChain txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func UpdateSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/UpdateSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	sideChainParam := new(SideChainParam)
	err = json.Unmarshal(data, sideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, sideChainParam.Path)
	if !ok {
		return false
	}
	CCMCAddress, err := hex.DecodeString(sideChainParam.CCMCAddress)
	if err != nil {
		ctx.LogError("hex.DecodeString error %v", err)
		return false
	}
	txHash, err := ctx.Ont.Native.Scm.UpdateSideChain(user.Address, sideChainParam.Chainid,
		sideChainParam.Router, sideChainParam.Name, sideChainParam.BlocksToWait, CCMCAddress, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Scm.UpdateSideChain error: %v", err)
		return false
	}
	ctx.LogInfo("UpdateSideChain txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type ApproveSideChainParam struct {
	Path    []string
	Chainid uint64
}

func ApproveRegisterSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveRegisterSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveSideChainParam := new(ApproveSideChainParam)
	err = json.Unmarshal(data, approveSideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range approveSideChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Scm.ApproveRegisterSideChain(approveSideChainParam.Chainid, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Scm.ApproveRegisterSideChain error: %v", err)
			return false
		}
		ctx.LogInfo("ApproveRegisterSideChain txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

func ApproveUpdateSideChain(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveUpdateSideChain.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveSideChainParam := new(ApproveSideChainParam)
	err = json.Unmarshal(data, approveSideChainParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range approveSideChainParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Scm.ApproveUpdateSideChain(approveSideChainParam.Chainid, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Scm.ApproveUpdateSideChain error: %v", err)
			return false
		}
		ctx.LogInfo("ApproveUpdateSideChain txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

type RegisterPeerParam struct {
	PeerPubkey string
	Path       string
}

func RegisterCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RegisterCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	registerPeerParam := new(RegisterPeerParam)
	err = json.Unmarshal(data, registerPeerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, registerPeerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Nm.RegisterCandidate(registerPeerParam.PeerPubkey, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.RegisterCandidate error: %v", err)
		return false
	}
	ctx.LogInfo("RegisterCandidate txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type PeerParam2 struct {
	PeerPubkey string
	Path       string
}

func UnRegisterCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/UnRegisterCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam2)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, peerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Nm.UnRegisterCandidate(peerParam.PeerPubkey, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.UnRegisterCandidate error: %v", err)
		return false
	}
	ctx.LogInfo("UnRegisterCandidate txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func QuitNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/QuitNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam2)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	user, ok := getAccountByPassword(ctx, peerParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Nm.QuitNode(peerParam.PeerPubkey, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.QuitNode error: %v", err)
		return false
	}
	ctx.LogInfo("QuitNode txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type PeerParam struct {
	PeerPubkey string
	Path       []string
}

func ApproveCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range peerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Nm.ApproveCandidate(peerParam.PeerPubkey, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Nm.ApproveCandidate error: %v", err)
			return false
		}
		ctx.LogInfo("ApproveCandidate txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

func RejectCandidate(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RejectCandidate.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range peerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Nm.RejectCandidate(peerParam.PeerPubkey, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Nm.RejectCandidate error: %v", err)
			return false
		}
		ctx.LogInfo("RejectCandidate txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

type PeerListParam struct {
	PeerPubkeyList []string
	Path           []string
}

func BlackNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/BlackNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerListParam := new(PeerListParam)
	err = json.Unmarshal(data, peerListParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range peerListParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Nm.BlackNode(peerListParam.PeerPubkeyList, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Nm.BlackNode error: %v", err)
			return false
		}
		ctx.LogInfo("BlackNode txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

func WhiteNode(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/WhiteNode.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	peerParam := new(PeerParam)
	err = json.Unmarshal(data, peerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range peerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Nm.WhiteNode(peerParam.PeerPubkey, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Nm.WhiteNode error: %v", err)
			return false
		}
		ctx.LogInfo("WhiteNode txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

type Configuration struct {
	BlockMsgDelay        uint32
	HashMsgDelay         uint32
	PeerHandshakeTimeout uint32
	MaxBlockChangeView   uint32
	Path                 []string
}

func UpdateConfig(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/UpdateConfig.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	configuration := new(Configuration)
	err = json.Unmarshal(data, configuration)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range configuration.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Nm.UpdateConfig(configuration.BlockMsgDelay, configuration.HashMsgDelay,
		configuration.PeerHandshakeTimeout, configuration.MaxBlockChangeView, users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.UpdateConfig error: %v", err)
		return false
	}
	ctx.LogInfo("CommitDpos txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type RelayerListParam struct {
	AddressList []string
	Path        string
}

func RegisterRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RegisterRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerListParam := new(RelayerListParam)
	err = json.Unmarshal(data, relayerListParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	addressList := make([]common.Address, 0)
	for _, addr := range relayerListParam.AddressList {
		address, err := common.AddressFromBase58(addr)
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
			return false
		}
		addressList = append(addressList, address)
	}

	user, ok := getAccountByPassword(ctx, relayerListParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.RegisterRelayer(addressList, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.RegisterRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("RegisterRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

func RemoveRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/RemoveRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	relayerListParam := new(RelayerListParam)
	err = json.Unmarshal(data, relayerListParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	addressList := make([]common.Address, 0)
	for _, addr := range relayerListParam.AddressList {
		address, err := common.AddressFromBase58(addr)
		if err != nil {
			ctx.LogError("common.AddressFromBase58 failed %v", err)
			return false
		}
		addressList = append(addressList, address)
	}

	user, ok := getAccountByPassword(ctx, relayerListParam.Path)
	if !ok {
		return false
	}
	txHash, err := ctx.Ont.Native.Rm.RemoveRelayer(addressList, user)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Rm.RemoveRelayer error: %v", err)
		return false
	}
	ctx.LogInfo("RemoveRelayer txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}

type ApproveRelayerParam struct {
	ID   uint64
	Path []string
}

func ApproveRegisterRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveRegisterRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveRelayerParam := new(ApproveRelayerParam)
	err = json.Unmarshal(data, approveRelayerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range approveRelayerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Rm.ApproveRegisterRelayer(approveRelayerParam.ID, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Rm.ApproveRegisterRelayer error: %v", err)
			return false
		}
		ctx.LogInfo("ApproveRegisterRelayer txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

func ApproveRemoveRelayer(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/ApproveRemoveRelayer.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	approveRelayerParam := new(ApproveRelayerParam)
	err = json.Unmarshal(data, approveRelayerParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	time.Sleep(1 * time.Second)
	for _, path := range approveRelayerParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		txHash, err := ctx.Ont.Native.Rm.ApproveRemoveRelayer(approveRelayerParam.ID, user)
		if err != nil {
			ctx.LogError("ctx.Ont.Native.Rm.ApproveRemoveRelayer error: %v", err)
			return false
		}
		ctx.LogInfo("ApproveRemoveRelayer txHash is: %v", txHash.ToHexString())
	}
	waitForBlock(ctx)
	return true
}

func GetConfig(ctx *testframework.TestFrameworkContext) bool {
	config, err := getConfig(ctx)
	if err != nil {
		ctx.LogError("getConfig failed %v", err)
		return false
	}

	fmt.Println("config.BlockMsgDelay is:", config.BlockMsgDelay)
	fmt.Println("config.HashMsgDelay is:", config.HashMsgDelay)
	fmt.Println("config.PeerHandshakeTimeout is:", config.PeerHandshakeTimeout)
	fmt.Println("config.MaxBlockChangeView is:", config.MaxBlockChangeView)
	return true
}

func getConfig(ctx *testframework.TestFrameworkContext) (*node_manager.Configuration, error) {
	contractAddress := utils.NodeManagerContractAddress
	config := new(node_manager.Configuration)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), []byte(node_manager.VBFT_CONFIG))
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := config.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize config error!")
	}
	return config, nil
}

func GetPeerPoolMap(ctx *testframework.TestFrameworkContext) bool {
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap failed %v", err)
		return false
	}

	for _, v := range peerPoolMap.PeerPoolMap {
		fmt.Println("###########################################")
		fmt.Println("peerPoolItem.Index is:", v.Index)
		fmt.Println("peerPoolItem.PeerPubkey is:", v.PeerPubkey)
		fmt.Println("peerPoolItem.Address is:", v.Address.ToBase58())
		fmt.Println("peerPoolItem.Status is:", v.Status)
	}
	return true
}

func GetGovernanceView(ctx *testframework.TestFrameworkContext) bool {
	governanceView, err := getGovernanceView(ctx)
	if err != nil {
		ctx.LogError("getGovernanceView failed %v", err)
		return false
	}
	fmt.Println("governanceView.View is:", governanceView.View)
	fmt.Println("governanceView.TxHash is:", governanceView.TxHash)
	fmt.Println("governanceView.Height is:", governanceView.Height)
	return true
}

func getGovernanceView(ctx *testframework.TestFrameworkContext) (*node_manager.GovernanceView, error) {
	contractAddress := utils.NodeManagerContractAddress
	governanceView := new(node_manager.GovernanceView)
	key := []byte(node_manager.GOVERNANCE_VIEW)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, fmt.Errorf("getStorage error")
	}
	if err := governanceView.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, fmt.Errorf("deserialize, deserialize governanceView error!")
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

func getPeerPoolMap(ctx *testframework.TestFrameworkContext) (*node_manager.PeerPoolMap, error) {
	contractAddress := utils.NodeManagerContractAddress
	view, err := getView(ctx)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getView error")
	}
	peerPoolMap := &node_manager.PeerPoolMap{
		PeerPoolMap: make(map[string]*node_manager.PeerPoolItem),
	}
	viewBytes := utils.GetUint32Bytes(view)
	key := ConcatKey([]byte(node_manager.PEER_POOL), viewBytes)
	value, err := ctx.Ont.GetStorage(contractAddress.ToHexString(), key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := peerPoolMap.Deserialization(common.NewZeroCopySource(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize peerPoolMap error!")
	}
	return peerPoolMap, nil
}

type CommitDposParam struct {
	Path []string
}

func CommitDpos(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/CommitDpos.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	commitDposParam := new(CommitDposParam)
	err = json.Unmarshal(data, commitDposParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	time.Sleep(1 * time.Second)
	for _, path := range commitDposParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
	}

	txHash, err := ctx.Ont.Native.Nm.CommitDpos(users)
	if err != nil {
		ctx.LogError("ctx.Ont.Native.Nm.UpdateConfig error: %v", err)
		return false
	}
	ctx.LogInfo("CommitDpos txHash is: %v", txHash.ToHexString())
	waitForBlock(ctx)
	return true
}
