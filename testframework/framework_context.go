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

package testframework

import (
	log4 "github.com/alecthomas/log4go"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-tool/common"
)

//TestFrameworkContext is the context for test case
type TestFrameworkContext struct {
	Ont       *sdk.OntologySdk //sdk to ontology
	failNowCh chan interface{}
}

//NewTestFrameworkContext return a TestFrameworkContext instance
func NewTestFrameworkContext(ont *sdk.OntologySdk, failNowCh chan interface{}) *TestFrameworkContext {
	return &TestFrameworkContext{
		Ont:       ont,
		failNowCh: failNowCh,
	}
}

//LogInfo log info in test case
func (this *TestFrameworkContext) LogInfo(arg0 interface{}, args ...interface{}) {
	log4.Info(arg0, args...)
}

//LogError log error info  when error occur in test case
func (this *TestFrameworkContext) LogError(arg0 interface{}, args ...interface{}) {
	log4.Error(arg0, args...)
}

//LogWarn log warning info in test case
func (this *TestFrameworkContext) LogWarn(arg0 interface{}, args ...interface{}) {
	log4.Warn(arg0, args...)
}

func (this *TestFrameworkContext) NewAccount() *sdk.Account {
	return sdk.NewAccount()
}

//FailNow will stop test, and skip all haven't not test case
func (this *TestFrameworkContext) FailNow() {
	select {
	case <-this.failNowCh:
	default:
		close(this.failNowCh)
	}
}

func (this *TestFrameworkContext) GetGasPrice() uint64 {
	return common.DefConfig.GasPrice
}

func (this *TestFrameworkContext) GetGasLimit() uint64 {
	return common.DefConfig.GasLimit
}
