package mocks

import (
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
	panic("implement me")
}

func (m MockTransactionRepository) Types() map[int]string {
	panic("implement me")
}

func (m MockTransactionRepository) TypesAsArray() []string {
	panic("implement me")
}

func (m MockTransactionRepository) Statuses() map[int]string {
	panic("implement me")
}
