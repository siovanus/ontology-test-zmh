package native

import (
	"github.com/ontio/ontology-test/testcase/smartcontract/native/governance_feeSplit"
)

func TestNative() {
	governance_feeSplit.TestGovernanceMethods()
	//governance_feeSplit.TestGovernanceContract()
	//governance_feeSplit.TestGovernanceContractError()
}
