package services

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/repositories/mocks"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
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
		ParentIndex:         1,
		ParentHash:          "0xparenthash",
	}
}

func exampleAmount() *types.Amount {
	return &types.Amount{
		Value: int64(1000),
	}
}

func exampleRequest() *rTypes.AccountBalanceRequest {
	return &rTypes.AccountBalanceRequest{
		AccountIdentifier: &rTypes.AccountIdentifier{Address: "0.0.1"},
		BlockIdentifier:   nil,
	}
}

func TestAccountBalance(t *testing.T) {
	// mocks:
	mockBlockRepository := &mocks.MockBlockRepository{}
	mockTransactionRepository := &mocks.MockTransactionRepository{}
	mockAccountRepo := &mocks.MockAccountRepository{}

	mockBlockRepository.On("RetrieveLatest").Return(exampleBlock())
	mockAccountRepo.On("RetrieveBalanceAtBlock").Return(exampleAmount())

	// given:
	exampleCommons := NewCommons(mockBlockRepository, mockTransactionRepository)
	exampleAccountService := NewAccountAPIService(exampleCommons, mockAccountRepo)
	expectedResult := &rTypes.AccountBalanceResponse{
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
	var expectedNilError *rTypes.Error = nil

	// when:
	actualResult, error := exampleAccountService.AccountBalance(nil, exampleRequest())

	// then:
	assert.Equal(t, expectedResult, actualResult)
	assert.Equal(t, expectedNilError, error)
}
