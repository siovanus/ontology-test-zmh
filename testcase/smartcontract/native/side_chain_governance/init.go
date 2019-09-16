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
	"github.com/ontio/ontology-tool/testframework"
)

func TestGovernanceMethods() {
	testframework.TFramework.RegTestCase("SyncGenesisHeader", SyncGenesisHeader)
	testframework.TFramework.RegTestCase("RegisterSideChain", RegisterSideChain)
	testframework.TFramework.RegTestCase("ApproveRegisterSideChain", ApproveRegisterSideChain)
	testframework.TFramework.RegTestCase("AssetMapping", AssetMapping)
	testframework.TFramework.RegTestCase("ApproveAssetMapping", ApproveAssetMapping)
}