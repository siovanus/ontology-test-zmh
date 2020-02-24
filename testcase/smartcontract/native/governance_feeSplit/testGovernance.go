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

package governance_feeSplit

import (
	"fmt"
	"github.com/ontio/ontology-tool/testframework"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

const (
	PEER_PUBKEY = "02890c587f4e4a6a98b455248eabac04b733580cfe5f11acd648c675543dfbb926"
	INIT_ONT    = 1000000
)

func SimulateAuthorizeForPeerAndWithdraw(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{1000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{2000}
	authorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-1000)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user2, INIT_ONT-2000)
	if !ok {
		return false
	}
	//check total stake
	totalStake, err := getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 1000 {
		ctx.LogError("total stake is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user2.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 2000 {
		ctx.LogError("total stake is error")
		return false
	}

	posList = []uint32{1000}
	unAuthorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{2000}
	unAuthorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-1000)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user2, INIT_ONT-2000)
	if !ok {
		return false
	}
	//check total stake
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 1000 {
		ctx.LogError("total stake is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user2.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 2000 {
		ctx.LogError("total stake is error")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	withdrawList := []uint32{500}
	withdraw(ctx, user1, peerPubkeyList, withdrawList)
	withdrawList = []uint32{1000}
	withdraw(ctx, user2, peerPubkeyList, withdrawList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-500)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user2, INIT_ONT-1000)
	if !ok {
		return false
	}
	//check total stake
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 500 {
		ctx.LogError("total stake is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user2.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 1000 {
		ctx.LogError("total stake is error")
		return false
	}

	withdrawList = []uint32{500}
	withdraw(ctx, user1, peerPubkeyList, withdrawList)
	withdrawList = []uint32{1000}
	withdraw(ctx, user2, peerPubkeyList, withdrawList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user2, INIT_ONT)
	if !ok {
		return false
	}
	//check total stake
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 0 {
		ctx.LogError("total stake is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user2.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 0 {
		ctx.LogError("total stake is error")
		return false
	}

	return true
}

func SimulateUnConsensusToConsensus(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{1000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{2000}
	authorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.CandidateStatus ||
		peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos != 3000 {
		ctx.LogError("peerPoolItem data 1 is wrong!")
		return false
	}

	//select in consensus
	posList = []uint32{100000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{500}
	unAuthorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.ConsensusStatus ||
		peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos != 102500 {
		fmt.Println(peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos)
		ctx.LogError("peerPoolItem data 2 is wrong!")
		return false
	}

	//check authorizeInfo data
	//user1
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user1.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 0 || authorizeInfo.ConsensusPos != 101000 {
		ctx.LogError("authorizeInfo data for user1 is wrong!")
		return false
	}
	//user2
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 0 || authorizeInfo.ConsensusPos != 1500 || authorizeInfo.WithdrawUnfreezePos != 500 {
		ctx.LogError("authorizeInfo data for user2 is wrong!")
		return false
	}
	return true
}

func SimulateRejectCandidate(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	registerCandidate(ctx, user1, PEER_PUBKEY2, 10000)
	waitForBlock(ctx)
	rejectCandidate(ctx, user, PEER_PUBKEY2)
	waitForBlock(ctx)
	//check pperPoolMap
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY2]
	if ok {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	//check authorizeInfo data
	//user1
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY2, user1.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.WithdrawUnfreezePos != 10000 {
		ctx.LogError("authorizeInfo data for user1 is wrong!")
		return false
	}
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-10000)
	if !ok {
		return false
	}
	return true
}

func SimulateUnConsensusToUnConsensus(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{1000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{2000}
	authorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	posList = []uint32{500}
	unAuthorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//select in unconsensus
	posList = []uint32{500}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{500}
	unAuthorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.CandidateStatus ||
		peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos != 2500 {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	//check authorizeInfo data
	//user1
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user1.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 1500 || authorizeInfo.ConsensusPos != 0 {
		ctx.LogError("authorizeInfo data for user1 is wrong!")
		return false
	}
	//user2
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 1000 || authorizeInfo.ConsensusPos != 0 || authorizeInfo.WithdrawUnfreezePos != 1000 {
		ctx.LogError("authorizeInfo data for user2 is wrong!")
		return false
	}
	return true
}

func SimulateConsensusToUnConsensus(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	//select in consensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{200000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//select in unconsensus
	posList = []uint32{199000}
	unAuthorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	posList = []uint32{1000}
	authorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.CandidateStatus ||
		peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos != 2000 {
		fmt.Println(peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos)
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	//check authorizeInfo data
	//user1
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user1.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 1000 || authorizeInfo.ConsensusPos != 0 || authorizeInfo.WithdrawCandidatePos != 199000 {
		ctx.LogError("authorizeInfo data for user1 is wrong!")
		return false
	}
	//user2
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 1000 || authorizeInfo.ConsensusPos != 0 {
		ctx.LogError("authorizeInfo data for user2 is wrong!")
		return false
	}
	return true
}

func SimulateConsensusToConsensus(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	//select in consensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{200000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//select in consensus
	posList = []uint32{100000}
	unAuthorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	posList = []uint32{1000}
	authorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.ConsensusStatus ||
		peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos != 101000 {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	//check authorizeInfo data
	//user1
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user1.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 0 || authorizeInfo.ConsensusPos != 100000 || authorizeInfo.WithdrawCandidatePos != 100000 {
		ctx.LogError("authorizeInfo data for user1 is wrong!")
		return false
	}
	//user2
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getAuthorizeInfo error :%v", err)
	}
	if authorizeInfo.NewPos != 0 || authorizeInfo.CandidatePos != 0 || authorizeInfo.ConsensusPos != 1000 {
		ctx.LogError("authorizeInfo data for user2 is wrong!")
		return false
	}
	return true
}

func SimulateQuitUnConsensus(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	//select in unconsensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{1000}
	authorizeForPeer(ctx, user, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//quit unconsensus
	quitNode(ctx, user, PEER_PUBKEY)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.QuitingStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY]
	if ok {
		ctx.LogError("peer quit failed, peerPoolMap is not deleted!")
		return false
	}
	//check balance
	ok = checkBalance(ctx, user, 1000000000-2*INIT_ONT-10000-1000)
	if !ok {
		return false
	}
	//check authorizeInfo data
	//user
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user.Address)
	if authorizeInfo.WithdrawUnfreezePos != 11000 {
		fmt.Println(authorizeInfo.WithdrawUnfreezePos)
		ctx.LogError("authorizeInfo of user is error!")
		return false
	}
	return true
}

func SimulateQuitConsensus(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	//select in consensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{200000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//quit consensus
	posList = []uint32{100000}
	unAuthorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	quitNode(ctx, user, PEER_PUBKEY)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.QuitConsensusStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.QuitingStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY]
	if ok {
		ctx.LogError("peer quit failed, peerPoolMap is not deleted!")
		return false
	}
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-200000)
	if !ok {
		return false
	}
	//check authorizeInfo data
	//user
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user.Address)
	if authorizeInfo.WithdrawUnfreezePos != 10000 {
		ctx.LogError("authorizeInfo of user is error!")
		return false
	}
	//user1
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user1.Address)
	if authorizeInfo.WithdrawUnfreezePos != 200000 {
		ctx.LogError("authorizeInfo of user1 is error!")
		return false
	}

	//register again
	registerCandidate(ctx, user, PEER_PUBKEY, 10000)
	waitForBlock(ctx)
	approveCandidate(ctx, user, PEER_PUBKEY)
	waitForBlock(ctx)

	//check index
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	peerPoolItem, ok := peerPoolMap.PeerPoolMap[PEER_PUBKEY]
	if !ok {
		ctx.LogError("peer register failed!")
		return false
	}
	if peerPoolItem.Index != 8 {
		ctx.LogError("index error!")
		return false
	}

	//register new peer
	registerCandidate(ctx, user, PEER_PUBKEY2, 10000)
	waitForBlock(ctx)
	approveCandidate(ctx, user, PEER_PUBKEY2)
	waitForBlock(ctx)

	//check index
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	peerPoolItem, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY2]
	if !ok {
		ctx.LogError("peer register failed!")
		return false
	}
	if peerPoolItem.Index != 9 {
		ctx.LogError("index error!")
		return false
	}

	return true
}

func SimulateBlackConsensusAndWhite(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	//select in consensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{200000}
	authorizeForPeer(ctx, user, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check total stake
	totalStake, err := getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != (10000 + 200000) {
		ctx.LogError("total stake is error")
		return false
	}

	//blacknode
	ok = blackNode(ctx, user, []string{PEER_PUBKEY})
	if !ok {
		return false
	}
	waitForBlock(ctx)

	//check if in blackList
	ok, err = inBlackList(ctx, PEER_PUBKEY)
	if !ok {
		ctx.LogError("peer should in blackList")
		return false
	}
	view, err := getView(ctx)
	if err != nil {
		ctx.LogError("getView error :%v", err)
		return false
	}
	if view != 3 {
		ctx.LogError("error: view changed failed!")
		return false
	}

	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY]
	if ok {
		ctx.LogError("peer quit failed, peerPoolMap is not deleted!")
	}
	//check authorizeInfo data
	//user
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user.Address)
	if authorizeInfo.WithdrawUnfreezePos != 190000 {
		ctx.LogError("authorizeInfo of user is error!")
		return false
	}
	//check balance
	ok = checkBalance(ctx, user, 1000000000-2*INIT_ONT-10000-200000)
	if !ok {
		return false
	}
	//check total stake
	totalStake, err = getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 190000 {
		ctx.LogError("total stake is error")
		return false
	}
	//check penaltyStake
	penaltyStake, err := getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPenaltyStake error :%v", err)
		return false
	}
	if penaltyStake.InitPos != 10000 || penaltyStake.AuthorizePos != (200000-190000) {
		ctx.LogError("penalty stake is error")
		return false
	}

	//whiteNode
	ok = whiteNode(ctx, user, PEER_PUBKEY)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	//check if in blackList
	ok, err = inBlackList(ctx, PEER_PUBKEY)
	if ok {
		ctx.LogError("peer should not in blackList!")
		return false
	}
	//check authorizeInfo data
	//user
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user.Address)
	if authorizeInfo.WithdrawUnfreezePos != 190000 {
		ctx.LogError("authorizeInfo of user is error!")
		return false
	}
	ok = checkBalance(ctx, user, 1000000000-2*INIT_ONT-10000-200000)
	if !ok {
		return false
	}
	//check total stake
	totalStake, err = getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 190000 {
		ctx.LogError("total stake is error")
		return false
	}
	//check penaltyStake
	penaltyStake, err = getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPenaltyStake error :%v", err)
		return false
	}
	if penaltyStake.InitPos != 10000 || penaltyStake.AuthorizePos != (200000-190000) {
		ctx.LogError("penalty stake is error")
		return false
	}
	return true
}

func SimulateBlackUnConsensusAndWhite(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	//select in consensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{1000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check total stake
	totalStake, err := getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 10000 {
		ctx.LogError("total stake user is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 1000 {
		ctx.LogError("total stake user1 is error")
		return false
	}

	//blacknode
	ok = blackNode(ctx, user, []string{PEER_PUBKEY})
	if !ok {
		return false
	}
	waitForBlock(ctx)
	//check total stake
	totalStake, err = getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 10000 {
		ctx.LogError("total stake user is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 1000 {
		ctx.LogError("total stake user1 is error")
		return false
	}

	//check if in blackList
	ok, err = inBlackList(ctx, PEER_PUBKEY)
	if !ok {
		ctx.LogError("peer should in blackList")
		return false
	}
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.BlackStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	view, err := getView(ctx)
	if err != nil {
		ctx.LogError("getView error :%v", err)
		return false
	}
	if view != 2 {
		ctx.LogError("error: view changed!")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY]
	if ok {
		ctx.LogError("peer quit failed, peerPoolMap is not deleted!")
	}
	//check total stake
	totalStake, err = getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 0 {
		ctx.LogError("total stake user is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 950 {
		ctx.LogError("total stake user1 is error")
		return false
	}
	//check authorizeInfo data
	//user
	authorizeInfo, err := getAuthorizeInfo(ctx, PEER_PUBKEY, user.Address)
	if err == nil {
		ctx.LogError("authorizeInfo of user is error!")
		return false
	}
	//user1
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user1.Address)
	if authorizeInfo.WithdrawUnfreezePos != 950 {
		ctx.LogError("authorizeInfo of user1 is error!")
		return false
	}
	//check balance
	ok = checkBalance(ctx, user, 1000000000-2*INIT_ONT-10000)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user1, INIT_ONT-1000)
	if !ok {
		return false
	}
	//check penaltyStake
	penaltyStake, err := getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPenaltyStake error :%v", err)
		return false
	}
	if penaltyStake.InitPos != 10000 || penaltyStake.AuthorizePos != 50 {
		ctx.LogError("penalty stake is error")
		return false
	}

	//whiteNode
	ok = whiteNode(ctx, user, PEER_PUBKEY)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	//check if in blackList
	ok, err = inBlackList(ctx, PEER_PUBKEY)
	if ok {
		ctx.LogError("peer should not in blackList!")
		return false
	}

	//check total stake
	totalStake, err = getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 0 {
		ctx.LogError("total stake user is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 950 {
		ctx.LogError("total stake user1 is error")
		return false
	}
	//check authorizeInfo data
	//user
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user.Address)
	if err == nil {
		ctx.LogError("authorizeInfo of user is error!")
		return false
	}
	//user1
	authorizeInfo, err = getAuthorizeInfo(ctx, PEER_PUBKEY, user1.Address)
	if authorizeInfo.WithdrawUnfreezePos != 950 {
		ctx.LogError("authorizeInfo of user1 is error!")
		return false
	}
	//check balance
	ok = checkBalance(ctx, user, 1000000000-2*INIT_ONT-10000)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user1, INIT_ONT-1000)
	if !ok {
		return false
	}
	//check penaltyStake
	penaltyStake, err = getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPenaltyStake error :%v", err)
		return false
	}
	if penaltyStake.InitPos != 10000 || penaltyStake.AuthorizePos != 50 {
		ctx.LogError("penalty stake is error")
		return false
	}
	return true
}

func SimulateUpdateConfig(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	//check config
	config, err := getVbftConfig(ctx)
	if err != nil {
		ctx.LogError("getVbftConfig error :%v", err)
		return false
	}
	if config.L != 112 || config.K != 7 || config.C != 2 || config.N != 7 || config.BlockMsgDelay != 10000 || config.HashMsgDelay != 10000 ||
		config.PeerHandshakeTimeout != 10 || config.MaxBlockChangeView != 1000 {
		ctx.LogError("config is error")
		return false
	}

	config = &governance.Configuration{
		N:                    7,
		C:                    2,
		K:                    7,
		L:                    112,
		BlockMsgDelay:        5000,
		HashMsgDelay:         5000,
		PeerHandshakeTimeout: 15,
		MaxBlockChangeView:   10000,
	}
	ok = updateConfig(ctx, user, config)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check config
	config, err = getVbftConfig(ctx)
	if err != nil {
		ctx.LogError("getVbftConfig error :%v", err)
		return false
	}
	if config.L != 112 || config.K != 7 || config.C != 2 || config.N != 7 || config.BlockMsgDelay != 5000 || config.HashMsgDelay != 5000 ||
		config.PeerHandshakeTimeout != 15 || config.MaxBlockChangeView != 10000 {
		ctx.LogError("config is error")
		return false
	}
	return true
}

func SimulateCommitDPosAuth(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	commitDpos(ctx, user1)
	waitForBlock(ctx)

	view, err := getView(ctx)
	if err != nil {
		ctx.LogError("getView error :%v", err)
		return false
	}
	if view != 1 {
		ctx.LogError("error: view changed!")
		return false
	}
	return true
}

func SimulateUpdateGlobalParam(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	globalParam := &governance.GlobalParam{
		CandidateFee: 500000000000,
		MinInitStake: 10000,
		CandidateNum: 8 * 8,
		PosLimit:     30,
		A:            70,
		B:            30,
		Yita:         7,
		Penalty:      10,
	}
	ok = updateGlobalParam(ctx, user, globalParam)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	//check config
	globalParam, err := getGlobalParam(ctx)
	if err != nil {
		ctx.LogError("getGlobalParam error :%v", err)
		return false
	}
	if globalParam.CandidateFee != 500000000000 || globalParam.MinInitStake != 10000 ||
		globalParam.CandidateNum != (8*8) || globalParam.PosLimit != 30 ||
		globalParam.A != 70 || globalParam.B != 30 || globalParam.Yita != 7 || globalParam.Penalty != 10 {
		ctx.LogError("globalParam is error")
		return false
	}
	return true
}

func SimulateUpdateGlobalParam2(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	globalParam2 := &governance.GlobalParam2{
		MinAuthorizePos:      1000,
		CandidateFeeSplitNum: 40,
	}
	ok = updateGlobalParam2(ctx, user, globalParam2)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	//check config
	globalParam2, err := getGlobalParam2(ctx)
	if err != nil {
		ctx.LogError("getGlobalParam2 error :%v", err)
		return false
	}
	if globalParam2.MinAuthorizePos != 1000 || globalParam2.CandidateFeeSplitNum != 40 {
		ctx.LogError("globalParam2 is error")
		return false
	}
	return true
}

func SimulateUpdateSplitCurve(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	splitCurve := &governance.SplitCurve{
		Yi: []uint32{
			0, 1, 2, 3, 4, 389401, 444491, 493282, 536257, 573866, 606531, 634645, 658574, 678660, 695220, 708550,
			718927, 726606, 731826, 734808, 735759, 734870, 732317, 728265, 722867, 716262, 708583, 699949, 690472, 680254, 669391,
			657969, 646069, 633765, 621124, 608209, 595076, 581778, 568361, 554869, 541342, 527814, 514317, 500882, 487534, 474297,
			461191, 448236, 435447, 422839, 410425, 398217, 386223, 374452, 362910, 351604, 340537, 329713, 319135, 308805, 298723,
			288890, 279306, 269969, 260879, 252033, 243429, 235066, 226939, 219045, 211382, 203945, 196731, 189736, 182955, 176384,
			170018, 163854, 157887, 152113, 146526, 141122, 135896, 130845, 125963, 121246, 116690, 112290, 108041, 103940, 99981,
			96162, 92477, 88923, 85496, 82192, 79006, 75936, 72977, 70126, 67380,
		},
	}
	ok = updateSplitCurve(ctx, user, splitCurve)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	//check config
	splitCurve, err := getSplitCurve(ctx)
	if err != nil {
		ctx.LogError("getSplitCurve error :%v", err)
		return false
	}
	if splitCurve.Yi[1] != 1 || splitCurve.Yi[2] != 2 || splitCurve.Yi[3] != 3 || splitCurve.Yi[4] != 4 {
		ctx.LogError("splitCurve is error")
		return false
	}
	return true
}

func SimulateTransferPenalty(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	//select in consensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{1000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check total stake
	totalStake, err := getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 10000 {
		ctx.LogError("total stake user is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 1000 {
		ctx.LogError("total stake user1 is error")
		return false
	}

	//blacknode
	ok = blackNode(ctx, user, []string{PEER_PUBKEY})
	if !ok {
		return false
	}
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-1000)
	if !ok {
		return false
	}
	//check penaltyStake
	penaltyStake, err := getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPenaltyStake1 error :%v", err)
		return false
	}
	if penaltyStake.InitPos != 10000 || penaltyStake.AuthorizePos != 50 {
		ctx.LogError("penalty stake is error")
		return false
	}

	ok = transferPenalty(ctx, user, PEER_PUBKEY, user1.Address)
	if !ok {
		return false
	}
	waitForBlock(ctx)

	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-1000+10000+50)
	if !ok {
		return false
	}
	//check penaltyStake
	penaltyStake, err = getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		return true
	}
	ctx.LogError("getPenaltyStake error, not deleted")
	return true
}

func SimulateOntIDAndAuth(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user3, ok := getAccount3(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	registerCandidate(ctx, user3, PEER_PUBKEY2, 10000)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY2]
	if ok {
		ctx.LogError("peer should not exist")
		return false
	}
	return true
}

func SimulateUnRegisterCandidate(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	registerCandidate(ctx, user1, PEER_PUBKEY2, 10000)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-10000)
	if !ok {
		return false
	}

	unRegisterCandidate(ctx, user1, PEER_PUBKEY2)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY2]
	if ok {
		ctx.LogError("peer should not exist")
		return false
	}
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-10000)
	if !ok {
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY2}
	withdrawList := []uint32{10000}
	withdraw(ctx, user1, peerPubkeyList, withdrawList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT)
	if !ok {
		return false
	}
	return true
}

func SimulateFeeSplit(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	user3, ok := getAccount3(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{2000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{1000}
	authorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//transfer ong to governance contract
	_, err := ctx.Ont.Native.Ong.WithdrawONG(ctx.GetGasPrice(), ctx.GetGasLimit(), user, user, 1000000000000)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.WithdrawONG error:%s", err)
		return false
	}
	waitForBlock(ctx)
	_, err = ctx.Ont.Native.Ong.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), user, user, utils.GovernanceContractAddress, 1000000000000)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.Transfer error:%s", err)
		return false
	}
	waitForBlock(ctx)

	commitDpos(ctx, user)
	waitForBlock(ctx)

	splitFeeAddress, err := getSplitFeeAddress(ctx, user3.Address)
	if err != nil {
		ctx.LogError("getSplitFeeAddress error:%s", err)
		return false
	}
	if splitFeeAddress.Amount != 499999999997 {
		ctx.LogError("splitFeeAddress error, is not 499999999997")
	}

	splitFeeAddress, err = getSplitFeeAddress(ctx, user.Address)
	if err != nil {
		ctx.LogError("getSplitFeeAddress error:%s", err)
		return false
	}
	if splitFeeAddress.Amount != 500000000000 {
		ctx.LogError("splitFeeAddress error, is not 500000000000")
	}

	withdrawFee(ctx, user1)
	withdrawFee(ctx, user2)
	withdrawFee(ctx, user3)
	waitForBlock(ctx)

	ok = checkOngBalance(ctx, user1, 0)
	if !ok {
		return false
	}
	ok = checkOngBalance(ctx, user2, 0)
	if !ok {
		return false
	}
	ok = checkOngBalance(ctx, user3, 499999999997)
	if !ok {
		return false
	}

	return true
}

func SimulateFeeSplit2(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	user3, ok := getAccount3(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	setPeerCost(ctx, user, PEER_PUBKEY, 10)
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{2000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	posList = []uint32{1000}
	authorizeForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//transfer ong to governance contract
	_, err := ctx.Ont.Native.Ong.WithdrawONG(ctx.GetGasPrice(), ctx.GetGasLimit(), user, user, 1000000000000)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.WithdrawONG error:%s", err)
		return false
	}
	waitForBlock(ctx)
	_, err = ctx.Ont.Native.Ong.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), user, user, utils.GovernanceContractAddress, 1000000000000)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.Transfer error:%s", err)
		return false
	}
	waitForBlock(ctx)

	commitDpos(ctx, user)
	waitForBlock(ctx)

	splitFeeAddress, err := getSplitFeeAddress(ctx, user3.Address)
	if err != nil {
		ctx.LogError("getSplitFeeAddress error:%s", err)
		return false
	}
	if splitFeeAddress.Amount != 499999999997 {
		ctx.LogError("splitFeeAddress error, is not 499999999997, is %v", splitFeeAddress.Amount)
	}

	splitFeeAddress, err = getSplitFeeAddress(ctx, user2.Address)
	if err != nil {
		ctx.LogError("getSplitFeeAddress error:%s", err)
		return false
	}
	if splitFeeAddress.Amount != 150000000000 {
		ctx.LogError("splitFeeAddress error, is not 150000000000, is %v", splitFeeAddress.Amount)
	}

	splitFeeAddress, err = getSplitFeeAddress(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getSplitFeeAddress error:%s", err)
		return false
	}
	if splitFeeAddress.Amount != 300000000000 {
		ctx.LogError("splitFeeAddress error, is not 300000000000, is %v", splitFeeAddress.Amount)
	}

	splitFeeAddress, err = getSplitFeeAddress(ctx, user.Address)
	if err != nil {
		ctx.LogError("getSplitFeeAddress error:%s", err)
		return false
	}
	if splitFeeAddress.Amount != 50000000000 {
		ctx.LogError("splitFeeAddress error, is not 50000000000, is %v", splitFeeAddress.Amount)
	}

	withdrawFee(ctx, user1)
	withdrawFee(ctx, user2)
	withdrawFee(ctx, user3)
	waitForBlock(ctx)

	ok = checkOngBalance(ctx, user1, 300000000000)
	if !ok {
		return false
	}
	ok = checkOngBalance(ctx, user2, 150000000000)
	if !ok {
		return false
	}
	ok = checkOngBalance(ctx, user3, 499999999997)
	if !ok {
		return false
	}

	return true
}

func SimulateChangeInitPos(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	reduceInitPos(ctx, user, PEER_PUBKEY, 500)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].InitPos != 10000 {
		ctx.LogError("reduceInitPos1 should be failed")
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint32{300000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos != 0 {
		ctx.LogError("authorizeForPeer should be failed")
		return false
	}

	addInitPos(ctx, user, PEER_PUBKEY, 10000)
	waitForBlock(ctx)
	changeMaxAuthorization(ctx, user, PEER_PUBKEY, 400000)
	waitForBlock(ctx)
	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint32{300000}
	authorizeForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].TotalPos != 300000 {
		ctx.LogError("authorizeForPeer should be success")
		return false
	}

	reduceInitPos(ctx, user, PEER_PUBKEY, 10000)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].InitPos != 20000 {
		ctx.LogError("reduceInitPos2 should be failed")
		return false
	}

	reduceInitPos(ctx, user1, PEER_PUBKEY, 5000)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].InitPos != 20000 {
		ctx.LogError("reduceInitPos3 should be failed")
		return false
	}

	reduceInitPos(ctx, user, PEER_PUBKEY, 5000)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].InitPos != 15000 {
		ctx.LogError("reduceInitPos3 should be success")
		return false
	}
	commitDpos(ctx, user)
	return true
}

func SimulatePromisePos(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	promisePos := &governance.PromisePos{
		PeerPubkey: PEER_PUBKEY,
		PromisePos: 5000,
	}
	setPromisePos(ctx, user1, promisePos)
	waitForBlock(ctx)
	promisePos, err := getPromisePos(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPromisePos error :%v", err)
		return false
	}
	if promisePos.PromisePos != 10000 {
		ctx.LogError("setPromisePos should be failed")
		return false
	}

	promisePos = &governance.PromisePos{
		PeerPubkey: PEER_PUBKEY,
		PromisePos: 5000,
	}
	setPromisePos(ctx, user, promisePos)
	waitForBlock(ctx)
	promisePos, err = getPromisePos(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPromisePos error :%v", err)
		return false
	}
	if promisePos.PromisePos != 5000 {
		ctx.LogError("setPromisePos should be success, promisePos.PromisePos is %v, should be 5000", promisePos.PromisePos)
		return false
	}

	reduceInitPos(ctx, user, PEER_PUBKEY, 5000)
	waitForBlock(ctx)
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].InitPos != 5000 {
		ctx.LogError("reduceInitPos3 should be success")
		return false
	}

	return true
}

func SimulateSetPeerCost(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	setPeerCost(ctx, user1, PEER_PUBKEY, 50)
	waitForBlock(ctx)
	peerAttributes, err := getAttributes(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getAttributes failed %v", err)
		return false
	}
	if peerAttributes.T2PeerCost != 100 || peerAttributes.T1PeerCost != 100 || peerAttributes.TPeerCost != 100 {
		ctx.LogError("peerAttributes1 error")
		return false
	}

	setPeerCost(ctx, user, PEER_PUBKEY, 50)
	waitForBlock(ctx)
	peerAttributes, err = getAttributes(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getAttributes failed %v", err)
		return false
	}
	if peerAttributes.T2PeerCost != 50 || peerAttributes.T1PeerCost != 100 || peerAttributes.TPeerCost != 100 {
		ctx.LogError("peerAttributes2 error")
		return false
	}

	setPeerCost(ctx, user, PEER_PUBKEY, 40)
	waitForBlock(ctx)
	peerAttributes, err = getAttributes(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getAttributes failed %v", err)
		return false
	}
	if peerAttributes.T2PeerCost != 40 || peerAttributes.T1PeerCost != 100 || peerAttributes.TPeerCost != 100 {
		ctx.LogError("peerAttributes3 error")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)

	setPeerCost(ctx, user, PEER_PUBKEY, 50)
	waitForBlock(ctx)
	peerAttributes, err = getAttributes(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getAttributes failed %v", err)
		return false
	}
	if peerAttributes.T2PeerCost != 50 || peerAttributes.T1PeerCost != 40 || peerAttributes.TPeerCost != 100 {
		fmt.Println(peerAttributes.T2PeerCost)
		fmt.Println(peerAttributes.T1PeerCost)
		fmt.Println(peerAttributes.TPeerCost)
		ctx.LogError("peerAttributes4 error")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)

	setPeerCost(ctx, user, PEER_PUBKEY, 60)
	waitForBlock(ctx)
	peerAttributes, err = getAttributes(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getAttributes failed %v", err)
		return false
	}
	if peerAttributes.T2PeerCost != 60 || peerAttributes.T1PeerCost != 50 || peerAttributes.TPeerCost != 40 {
		ctx.LogError("peerAttributes5 error")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)

	peerAttributes, err = getAttributes(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getAttributes failed %v", err)
		return false
	}
	if peerAttributes.T2PeerCost != 60 || peerAttributes.T1PeerCost != 60 || peerAttributes.TPeerCost != 50 {
		ctx.LogError("peerAttributes6 error")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)

	peerAttributes, err = getAttributes(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getAttributes failed %v", err)
		return false
	}
	if peerAttributes.T2PeerCost != 60 || peerAttributes.T1PeerCost != 60 || peerAttributes.TPeerCost != 60 {
		ctx.LogError("peerAttributes7 error")
		return false
	}
	return true
}

func SimulateAddConsensusPeer(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	//check consensus peer
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	count := 0
	for _, v := range peerPoolMap.PeerPoolMap {
		if v.Status == governance.ConsensusStatus {
			count += 1
		}
	}
	if count != 7 {
		ctx.LogError("consensus node num 1 is wrong")
		return false
	}

	//update vbft config
	config := &governance.Configuration{
		N:                    8,
		C:                    2,
		K:                    8,
		L:                    128,
		BlockMsgDelay:        5000,
		HashMsgDelay:         5000,
		PeerHandshakeTimeout: 15,
		MaxBlockChangeView:   10000,
	}
	ok = updateConfig(ctx, user, config)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check consensus peer
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	count = 0
	for _, v := range peerPoolMap.PeerPoolMap {
		if v.Status == governance.ConsensusStatus {
			count += 1
		}
	}
	if count != 8 {
		ctx.LogError("consensus node num 2 is wrong")
		return false
	}

	return true
}
