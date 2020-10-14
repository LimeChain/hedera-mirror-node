package mocks

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/stretchr/testify/mock"
)

type MockAddressBookEntryRepository struct {
	mock.Mock
}

func (m MockAddressBookEntryRepository) Entries() (*types.AddressBookEntries, *rTypes.Error) {
	panic("implement me")
}
