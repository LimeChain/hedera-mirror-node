package services

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exampleBlockRequest() *rTypes.BlockRequest {
	index := int64(100)
	hash := "somehashh"

	return &rTypes.BlockRequest{
		NetworkIdentifier: &rTypes.NetworkIdentifier{
			Blockchain:           "Hedera",
			Network:              "hhh",
			SubNetworkIdentifier: nil,
		},
		BlockIdentifier: &rTypes.PartialBlockIdentifier{
			Index: &index,
			Hash:  &hash,
		},
	}
}

func exampleBlockResponse() *rTypes.BlockResponse {
	return &rTypes.BlockResponse{
		Block: &rTypes.Block{
			BlockIdentifier: &rTypes.BlockIdentifier{
				Index: 1,
				Hash:  "0x123jsjs",
			},
			ParentBlockIdentifier: &rTypes.BlockIdentifier{
				Index: 2,
				Hash:  "0xparenthash",
			},
			Timestamp:    1,
			Transactions: exampleRosettaTransactions(),
			Metadata:     nil,
		},
		OtherTransactions: nil,
	}
}

func TestBlock(t *testing.T) {
	exampleBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}

	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(exampleBlock, mocks.NilError)
	mocks.MTransactionRepository.On("FindBetween").Return(exampleTransactions(), mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	res, e := blockService.Block(nil, exampleBlockRequest())
	assert.Nil(t, e)
	assert.Equal(t, exampleBlockResponse(), res)
}

func TestBlockThrowsWhenFindByIdentifierFails(t *testing.T) {
	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(
		mocks.NilBlock,
		&rTypes.Error{},
	)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	res, e := blockService.Block(nil, exampleBlockRequest())
	assert.NotNil(t, e)
	assert.Nil(t, res)
}

func TestBlockThrowsWhenFindBetweenFails(t *testing.T) {
	exampleBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}

	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(exampleBlock, mocks.NilError)
	mocks.MTransactionRepository.On("FindBetween").Return(
		[]*types.Transaction{},
		&rTypes.Error{},
	)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	res, e := blockService.Block(nil, exampleBlockRequest())
	assert.NotNil(t, e)
	assert.Nil(t, res)
}

func exampleTransactions() []*types.Transaction {
	return []*types.Transaction{
		{
			Hash:       "123",
			Operations: nil,
		},
		{
			Hash:       "246",
			Operations: nil,
		},
	}
}

func exampleRosettaTransactions() []*rTypes.Transaction {
	return []*rTypes.Transaction{
		{
			TransactionIdentifier: &rTypes.TransactionIdentifier{Hash: "123"},
			Operations:            []*rTypes.Operation{},
			Metadata:              nil,
		},
		{
			TransactionIdentifier: &rTypes.TransactionIdentifier{Hash: "246"},
			Operations:            []*rTypes.Operation{},
			Metadata:              nil,
		},
	}
}
