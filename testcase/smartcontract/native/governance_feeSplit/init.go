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
	"github.com/ontio/ontology-tool/testframework"
)

func TestGovernanceMethods() {
	testframework.TFramework.RegTestCase("RegIdWithPublicKey", RegIdWithPublicKey)
	testframework.TFramework.RegTestCase("AssignFuncsToRole", AssignFuncsToRole)
	testframework.TFramework.RegTestCase("AssignFuncsToRoleAny", AssignFuncsToRoleAny)
	testframework.TFramework.RegTestCase("AssignOntIDsToRole", AssignOntIDsToRole)
	testframework.TFramework.RegTestCase("AssignOntIDsToRoleAny", AssignOntIDsToRoleAny)
	testframework.TFramework.RegTestCase("RegisterCandidate", RegisterCandidate)
	testframework.TFramework.RegTestCase("RegisterCandidate2Sign", RegisterCandidate2Sign)
	testframework.TFramework.RegTestCase("UnRegisterCandidate", UnRegisterCandidate)
	testframework.TFramework.RegTestCase("ApproveCandidate", ApproveCandidate)
	testframework.TFramework.RegTestCase("RejectCandidate", RejectCandidate)
	testframework.TFramework.RegTestCase("ChangeMaxAuthorization", ChangeMaxAuthorization)
	testframework.TFramework.RegTestCase("SetPeerCost", SetPeerCost)
	testframework.TFramework.RegTestCase("AddInitPos", AddInitPos)
	testframework.TFramework.RegTestCase("ReduceInitPos", ReduceInitPos)
	testframework.TFramework.RegTestCase("AuthorizeForPeer", AuthorizeForPeer)
	testframework.TFramework.RegTestCase("UnAuthorizeForPeer", UnAuthorizeForPeer)
	testframework.TFramework.RegTestCase("Withdraw", Withdraw)
	testframework.TFramework.RegTestCase("QuitNode", QuitNode)
	testframework.TFramework.RegTestCase("BlackNode", BlackNode)
	testframework.TFramework.RegTestCase("WhiteNode", WhiteNode)
	testframework.TFramework.RegTestCase("CommitDpos", CommitDpos)
	testframework.TFramework.RegTestCase("UpdateConfig", UpdateConfig)
	testframework.TFramework.RegTestCase("UpdateGlobalParam", UpdateGlobalParam)
	testframework.TFramework.RegTestCase("UpdateGlobalParam2", UpdateGlobalParam2)
	testframework.TFramework.RegTestCase("UpdateSplitCurve", UpdateSplitCurve)
	testframework.TFramework.RegTestCase("TransferPenalty", TransferPenalty)
	testframework.TFramework.RegTestCase("SetPromisePos", SetPromisePos)
	testframework.TFramework.RegTestCase("SetGasAddress", SetGasAddress)
	testframework.TFramework.RegTestCase("GetGasAddress", GetGasAddress)
	testframework.TFramework.RegTestCase("GetVbftConfig", GetVbftConfig)
	testframework.TFramework.RegTestCase("GetPreConfig", GetPreConfig)
	testframework.TFramework.RegTestCase("GetGlobalParam", GetGlobalParam)
	testframework.TFramework.RegTestCase("GetGlobalParam2", GetGlobalParam2)
	testframework.TFramework.RegTestCase("GetSplitCurve", GetSplitCurve)
	testframework.TFramework.RegTestCase("GetGovernanceView", GetGovernanceView)
	testframework.TFramework.RegTestCase("GetPeerPoolItem", GetPeerPoolItem)
	testframework.TFramework.RegTestCase("GetPeerPoolMap", GetPeerPoolMap)
	testframework.TFramework.RegTestCase("GetAuthorizeInfo", GetAuthorizeInfo)
	testframework.TFramework.RegTestCase("GetTotalStake", GetTotalStake)
	testframework.TFramework.RegTestCase("GetPenaltyStake", GetPenaltyStake)
	testframework.TFramework.RegTestCase("GetAttributes", GetAttributes)
	testframework.TFramework.RegTestCase("GetSplitFee", GetSplitFee)
	testframework.TFramework.RegTestCase("GetSplitFeeAddress", GetSplitFeeAddress)
	testframework.TFramework.RegTestCase("GetPromisePos", GetPromisePos)
	testframework.TFramework.RegTestCase("InBlackList", InBlackList)
	testframework.TFramework.RegTestCase("WithdrawOng", WithdrawOng)
	testframework.TFramework.RegTestCase("Vrf", Vrf)
	testframework.TFramework.RegTestCase("MultiTransferOnt", MultiTransferOnt)
	testframework.TFramework.RegTestCase("MultiTransferOng", MultiTransferOng)
	testframework.TFramework.RegTestCase("TransferOntMultiSign", TransferOntMultiSign)
	testframework.TFramework.RegTestCase("TransferOngMultiSign", TransferOngMultiSign)
	testframework.TFramework.RegTestCase("TransferFromOngMultiSign", TransferFromOngMultiSign)
	testframework.TFramework.RegTestCase("TransferOntMultiSignAddress", TransferOntMultiSignAddress)
	testframework.TFramework.RegTestCase("TransferOngMultiSignAddress", TransferOngMultiSignAddress)
	testframework.TFramework.RegTestCase("TransferFromOngMultiSignAddress", TransferFromOngMultiSignAddress)
	testframework.TFramework.RegTestCase("GetAddressMultiSign", GetAddressMultiSign)
	testframework.TFramework.RegTestCase("TransferOntMultiSignToMultiSign", TransferOntMultiSignToMultiSign)
	testframework.TFramework.RegTestCase("TransferOngMultiSignToMultiSign", TransferOngMultiSignToMultiSign)
	testframework.TFramework.RegTestCase("TransferFromOngMultiSignToMultiSign", TransferFromOngMultiSignToMultiSign)
	testframework.TFramework.RegTestCase("GetVbftInfo", GetVbftInfo)
}
