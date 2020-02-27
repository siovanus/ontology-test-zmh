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
	"fmt"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
	scommon "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/password"
)

func getAccountByPassword(ctx *testframework.TestFrameworkContext, path string) (*sdk.Account, bool) {
	wallet, err := ctx.Ont.OpenWallet(path)
	if err != nil {
		ctx.LogError("open wallet error:%s", err)
		return nil, false
	}
	pwd, err := password.GetPassword()
	if err != nil {
		ctx.LogError("getPassword error:%s", err)
		return nil, false
	}
	user, err := wallet.GetDefaultAccount(pwd)
	if err != nil {
		ctx.LogError("getDefaultAccount error:%s", err)
		return nil, false
	}
	return user, true
}

func invokeNativeContractWithMultiSign(
	ctx *testframework.TestFrameworkContext,
	gasPrice,
	gasLimit uint64,
	pubKeys []keypair.PublicKey,
	singers []*sdk.Account,
	cversion byte,
	contractAddress scommon.Address,
	method string,
	params []interface{},
) (scommon.Uint256, error) {
	tx, err := ctx.Ont.Native.NewNativeInvokeTransaction(gasPrice, gasLimit, cversion, contractAddress, method, params)
	if err != nil {
		return scommon.UINT256_EMPTY, err
	}
	for _, singer := range singers {
		err = ctx.Ont.MultiSignToTransaction(tx, uint16((5*len(pubKeys)+6)/7), pubKeys, singer)
		if err != nil {
			return scommon.UINT256_EMPTY, err
		}
	}
	return ctx.Ont.SendTransaction(tx)
}

func waitForBlock(ctx *testframework.TestFrameworkContext) bool {
	_, err := ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false
	}
	return true
}

func ConcatKey(args ...[]byte) []byte {
	temp := []byte{}
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func getEvent(ctx *testframework.TestFrameworkContext, txHash scommon.Uint256) bool {
	_, err := ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error: %s", err)
		return false
	}
	events, err := ctx.Ont.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		ctx.LogError("GetSmartContractEvent error: %s", err)
		return false
	}

	if events.State == 0 {
		ctx.LogWarn("ontio contract invoke failed, state:0")
		return false
	}

	if len(events.Notify) > 0 {
		states := events.Notify[0].States
		ctx.LogInfo("result is : %+v", states)
		return true
	} else {
		return false
	}

}

func getAddressByHexString(hexString string) (scommon.Address, error) {
	contractByte, err := hex.DecodeString(hexString)
	if err != nil {
		return scommon.Address{}, fmt.Errorf("hex.DecodeString failed %v", err)
	}
	contractAddress, err := scommon.AddressParseFromBytes(scommon.ToArrayReverse(contractByte))
	if err != nil {
		return scommon.Address{}, fmt.Errorf("common.AddressParseFromBytes failed %v", err)
	}
	return contractAddress, nil
}
