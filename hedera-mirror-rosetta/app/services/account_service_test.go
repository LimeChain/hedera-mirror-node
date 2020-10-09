package services

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type BlockRepoMock struct {
	mock.Mock
}

func (brm *BlockRepoMock) RetrieveLatest() (*types.Block, *rTypes.Error) {
	return &types.Block{}, nil
}

func (brm *BlockRepoMock) FindByIndex(index int64) (*types.Block, *rTypes.Error) {
	panic("implement me")
}

func (brm *BlockRepoMock) FindByHash(hash string) (*types.Block, *rTypes.Error) {
	panic("implement me")
}

func (brm *BlockRepoMock) FindByIdentifier(index int64, hash string) (*types.Block, *rTypes.Error) {
	panic("implement me")
}

func (brm *BlockRepoMock) RetrieveGenesis() (*types.Block, *rTypes.Error) {
	panic("implement me")
}

type TransactionRepoMock struct{}

type AccountRepoMock struct {
	mock.Mock
}

func (a AccountRepoMock) RetrieveBalanceAtBlock(addressStr string, consensusEnd int64) (*types.Amount, *rTypes.Error) {
	return nil, nil
}

func (t TransactionRepoMock) FindByHashInBlock(identifier string, consensusStart int64, consensusEnd int64) (*types.Transaction, *rTypes.Error) {
	panic("implement me")
}

func (t TransactionRepoMock) FindBetween(start int64, end int64) ([]*types.Transaction, *rTypes.Error) {
	panic("implement me")
}

func (t TransactionRepoMock) Types() map[int]string {
	panic("implement me")
}

func (t TransactionRepoMock) TypesAsArray() []string {
	panic("implement me")
}

func (t TransactionRepoMock) Statuses() map[int]string {
	panic("implement me")
}

func TestAccountBalance(t *testing.T) {
	req := &rTypes.AccountBalanceRequest{
		AccountIdentifier: &rTypes.AccountIdentifier{
			Address:    "0.0.1",
			SubAccount: nil,
			Metadata:   nil,
		},
		BlockIdentifier: nil,
	}

	brm := &BlockRepoMock{}
	brm.On("RetrieveLatest").Return(&types.Block{
		Index: 1,
		Hash:  "123",
	})
	commons := NewCommons(brm, &TransactionRepoMock{})
	arm := &AccountRepoMock{}
	arm.On("RetrieveBalanceAtBlock").Return(&types.Amount{
		Value: 1000,
	})
	as := NewAccountAPIService(commons, arm)

	abr, _ := as.AccountBalance(nil, req)
	expectedAbr := &rTypes.AccountBalanceResponse{
		BlockIdentifier: &rTypes.BlockIdentifier{
			Index: 1,
			Hash:  "0x123",
		},
		Balances: []*rTypes.Amount{
			{Value: "1000", Currency: config.CurrencyHbar},
		},
	}
	assert.Equal(t, expectedAbr, abr)
}
