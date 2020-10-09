package mocks

import (
	"fmt"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m MockAccountRepository) RetrieveBalanceAtBlock(addressStr string, consensusEnd int64) (*types.Amount, *rTypes.Error) {
	fmt.Println("Retrieving Balance At Block")
	args := m.Called()
	return args.Get(0).(*types.Amount), nil
}
