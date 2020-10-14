package mocks

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
)

var MBlockRepository *MockBlockRepository
var MTransactionRepository *MockTransactionRepository
var MAccountRepository *MockAccountRepository
var MAddressBookEntryRepository *MockAddressBookEntryRepository

var NilBlock *types.Block = nil
var NilError *rTypes.Error = nil
var NilTransaction *rTypes.Transaction = nil

func Setup() {
	MBlockRepository = &MockBlockRepository{}
	MTransactionRepository = &MockTransactionRepository{}
	MAccountRepository = &MockAccountRepository{}
	MAddressBookEntryRepository = &MockAddressBookEntryRepository{}
}
