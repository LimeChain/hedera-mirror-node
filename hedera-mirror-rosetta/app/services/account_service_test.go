package services

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/repositories/mocks"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var mockBlockRepository *mocks.MockBlockRepository
var mockTransactionRepository *mocks.MockTransactionRepository
var mockAccountRepo *mocks.MockAccountRepository
var exampleAccountService *AccountAPIService
var exampleCommons Commons
var nilBlock *types.Block = nil
var nilError *rTypes.Error = nil

func setup() {
	mockBlockRepository = &mocks.MockBlockRepository{}
	mockTransactionRepository = &mocks.MockTransactionRepository{}
	mockAccountRepo = &mocks.MockAccountRepository{}
	exampleCommons = NewCommons(mockBlockRepository, mockTransactionRepository)
	exampleAccountService = NewAccountAPIService(exampleCommons, mockAccountRepo)
}

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
	setup()
	mockBlockRepository.On("RetrieveLatest").Return(exampleBlock(), nilError)
	mockAccountRepo.On("RetrieveBalanceAtBlock").Return(exampleAmount(), nilError)

	// when:
	actualResult, e := exampleAccountService.AccountBalance(nil, exampleRequest(false))

	// then:
	assert.Equal(t, expectedAccountBalanceResponse(), actualResult)
	assert.Nil(t, e)
}

func TestAccountBalanceFullData(t *testing.T) {
	// mocks:
	setup()
	mockBlockRepository.On("FindByIdentifier").Return(exampleBlock(), nilError)
	mockAccountRepo.On("RetrieveBalanceAtBlock").Return(exampleAmount(), nilError)

	// when:
	actualResult, e := exampleAccountService.AccountBalance(nil, exampleRequest(true))

	//then:
	assert.Equal(t, expectedAccountBalanceResponse(), actualResult)
	assert.Nil(t, e)
	mockBlockRepository.AssertNotCalled(t, "RetrieveLatest")
}

func TestAccountBalanceThrows(t *testing.T) {
	// mocks:
	setup()

	mockBlockRepository.On("RetrieveLatest").Return(nilBlock, &rTypes.Error{})

	// when:
	actualResult, e := exampleAccountService.AccountBalance(nil, exampleRequest(false))

	// then:
	assert.Nil(t, actualResult)
	assert.NotNil(t, e)
}
