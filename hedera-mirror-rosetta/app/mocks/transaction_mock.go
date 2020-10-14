package mocks

import (
	"fmt"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m MockTransactionRepository) FindByHashInBlock(identifier string, consensusStart int64, consensusEnd int64) (*types.Transaction, *rTypes.Error) {
	panic("implement me")
}

func (m MockTransactionRepository) FindBetween(start int64, end int64) ([]*types.Transaction, *rTypes.Error) {
	fmt.Println("Finding bеtween")
	args := m.Called()
	return args.Get(0).([]*types.Transaction), args.Get(1).(*rTypes.Error)
}

func (m MockTransactionRepository) Types() map[int]string {
	panic("implement me")
}

func (m MockTransactionRepository) TypesAsArray() []string {
	fmt.Println("Types As Array")
	args := m.Called()
	return args.Get(0).([]string)
}

func (m MockTransactionRepository) Statuses() map[int]string {
	fmt.Println("Statuses")
	args := m.Called()
	return args.Get(0).(map[int]string)
}
