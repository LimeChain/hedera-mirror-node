package services

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/mocks"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exampleBlock() *types.Block {
	return &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}
}

func exampleAmount() *types.Amount {
	return &types.Amount{
		Value: int64(1000),
	}
}

func exampleRequest(withBlockIdentifier bool) *rTypes.AccountBalanceRequest {
	var blockIdentifier *rTypes.PartialBlockIdentifier = nil
	if withBlockIdentifier {
		index := int64(1)
		hash := "0x123"
		blockIdentifier = &rTypes.PartialBlockIdentifier{
			Index: &index,
			Hash:  &hash,
		}
	}
	return &rTypes.AccountBalanceRequest{
		AccountIdentifier: &rTypes.AccountIdentifier{Address: "0.0.1"},
		BlockIdentifier:   blockIdentifier,
	}
}

func expectedAccountBalanceResponse() *rTypes.AccountBalanceResponse {
	return &rTypes.AccountBalanceResponse{
		BlockIdentifier: &rTypes.BlockIdentifier{
			Index: 1,
			Hash:  "0x123jsjs",
		},
		Balances: []*rTypes.Amount{
			{
				Value:    "1000",
				Currency: config.CurrencyHbar,
			},
		},
	}
}

func TestAccountBalance(t *testing.T) {
	// mocks:
	mocks.Setup()
	mocks.MBlockRepository.On("RetrieveLatest").Return(exampleBlock(), mocks.NilError)
	mocks.MAccountRepository.On("RetrieveBalanceAtBlock").Return(exampleAmount(), mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	accountService := NewAccountAPIService(commons, mocks.MAccountRepository)

	// when:
	actualResult, e := accountService.AccountBalance(nil, exampleRequest(false))

	// then:
	assert.Equal(t, expectedAccountBalanceResponse(), actualResult)
	assert.Nil(t, e)
}

func TestAccountBalanceFullData(t *testing.T) {
	// mocks:
	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(exampleBlock(), mocks.NilError)
	mocks.MAccountRepository.On("RetrieveBalanceAtBlock").Return(exampleAmount(), mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	accountService := NewAccountAPIService(commons, mocks.MAccountRepository)

	// when:
	actualResult, e := accountService.AccountBalance(nil, exampleRequest(true))

	//then:
	assert.Equal(t, expectedAccountBalanceResponse(), actualResult)
	assert.Nil(t, e)
	mocks.MBlockRepository.AssertNotCalled(t, "RetrieveLatest")
}

func TestAccountBalanceThrows(t *testing.T) {
	// mocks:
	mocks.Setup()

	mocks.MBlockRepository.On("RetrieveLatest").Return(mocks.NilBlock, &rTypes.Error{})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	accountService := NewAccountAPIService(commons, mocks.MAccountRepository)

	// when:
	actualResult, e := accountService.AccountBalance(nil, exampleRequest(false))

	// then:
	assert.Nil(t, actualResult)
	assert.NotNil(t, e)
}
