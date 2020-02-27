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
	"encoding/json"
	"io/ioutil"
	"time"

	asdk "github.com/ontio/multi-chain-go-sdk"
	"github.com/ontio/ontology-crypto/keypair"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/testframework"
)

type SyncGenesisHeaderParam struct {
	Path     []string
	ChainRpc string
}

func SyncGenesisHeader(ctx *testframework.TestFrameworkContext) bool {
	data, err := ioutil.ReadFile("./side_chain_params/SyncGenesisHeader.json")
	if err != nil {
		ctx.LogError("ioutil.ReadFile failed %v", err)
		return false
	}
	SyncGenesisHeaderParam := new(SyncGenesisHeaderParam)
	err = json.Unmarshal(data, SyncGenesisHeaderParam)
	if err != nil {
		ctx.LogError("json.Unmarshal failed %v", err)
		return false
	}

	var users []*sdk.Account
	var pubKeys []keypair.PublicKey
	time.Sleep(1 * time.Second)
	for _, path := range SyncGenesisHeaderParam.Path {
		user, ok := getAccountByPassword(ctx, path)
		if !ok {
			return false
		}
		users = append(users, user)
		pubKeys = append(pubKeys, user.PublicKey)
	}

	aSdk := asdk.NewMultiChainSdk()
	aSdk.NewRpcClient().SetAddress(SyncGenesisHeaderParam.ChainRpc)
	genesisBlock, err := aSdk.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("get side chain genesis block error: %s", err)
		return false
	}
	genesisBlockHeader := genesisBlock.Header.ToArray()

	ok := syncGenesisHeader(ctx, pubKeys, users, genesisBlockHeader)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	return true
}
