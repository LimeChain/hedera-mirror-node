/*-
 * ‌
 * Hedera Mirror Node
 *
 * Copyright (C) 2019 - 2020 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ‍
 */

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
			Timestamp: 1,
			Transactions: []*rTypes.Transaction{
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
			},
			Metadata: nil,
		},
		OtherTransactions: nil,
	}
}

func TestBlock(t *testing.T) {
	// given:
	exampleBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}

	exampleTransactions := []*types.Transaction{
		{
			Hash:       "123",
			Operations: nil,
		},
		{
			Hash:       "246",
			Operations: nil,
		},
	}

	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(exampleBlock, mocks.NilError)
	mocks.MTransactionRepository.On("FindBetween").Return(exampleTransactions, mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	// when:
	res, e := blockService.Block(nil, exampleBlockRequest())

	// then:
	assert.Nil(t, e)
	assert.Equal(t, exampleBlockResponse(), res)
}

func TestBlockThrowsWhenFindByIdentifierFails(t *testing.T) {
	// given:
	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(
		mocks.NilBlock,
		&rTypes.Error{},
	)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	// when:
	res, e := blockService.Block(nil, exampleBlockRequest())

	// then:
	assert.NotNil(t, e)
	assert.Nil(t, res)
}

func TestBlockThrowsWhenFindBetweenFails(t *testing.T) {
	// given:
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

	// when:
	res, e := blockService.Block(nil, exampleBlockRequest())

	// then:
	assert.NotNil(t, e)
	assert.Nil(t, res)
}

func TestBlockTransaction(t *testing.T) {
	// given:
	exampleBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}

	exampleTransaction := &types.Transaction{
		Hash:       "somehash",
		Operations: nil,
	}

	exampleTransactionRequest := &rTypes.BlockTransactionRequest{
		NetworkIdentifier: &rTypes.NetworkIdentifier{
			Blockchain: "someblockchain",
			Network:    "somenetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "somesubnetwork",
				Metadata: nil,
			},
		},
		BlockIdentifier: &rTypes.BlockIdentifier{
			Index: 1,
			Hash:  "someblockhash",
		},
		TransactionIdentifier: &rTypes.TransactionIdentifier{Hash: "somehash"},
	}

	expectedResult := &rTypes.BlockTransactionResponse{Transaction: &rTypes.Transaction{
		TransactionIdentifier: &rTypes.TransactionIdentifier{Hash: "somehash"},
		Operations:            []*rTypes.Operation{},
		Metadata:              nil,
	}}

	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(exampleBlock, mocks.NilError)
	mocks.MTransactionRepository.On("FindByHashInBlock").Return(exampleTransaction, mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	// when:
	res, e := blockService.BlockTransaction(nil, exampleTransactionRequest)

	// then:
	assert.Equal(t, expectedResult, res)
	assert.Nil(t, e)
}

func TestBlockTransactionThrowsWhenFindByIdentifierFails(t *testing.T) {
	// given:
	exampleTransactionRequest := &rTypes.BlockTransactionRequest{
		NetworkIdentifier: &rTypes.NetworkIdentifier{
			Blockchain: "someblockchain",
			Network:    "somenetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "somesubnetwork",
				Metadata: nil,
			},
		},
		BlockIdentifier: &rTypes.BlockIdentifier{
			Index: 1,
			Hash:  "someblockhash",
		},
		TransactionIdentifier: &rTypes.TransactionIdentifier{Hash: "somehash"},
	}

	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(mocks.NilBlock, &rTypes.Error{})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	// when:
	res, e := blockService.BlockTransaction(nil, exampleTransactionRequest)

	// then:
	assert.Nil(t, res)
	assert.NotNil(t, e)
}

func TestBlockTransactionThrowsWhenFindByHashInBlockFails(t *testing.T) {
	// given:
	exampleBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}

	exampleTransactionRequest := &rTypes.BlockTransactionRequest{
		NetworkIdentifier: &rTypes.NetworkIdentifier{
			Blockchain: "someblockchain",
			Network:    "somenetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "somesubnetwork",
				Metadata: nil,
			},
		},
		BlockIdentifier: &rTypes.BlockIdentifier{
			Index: 1,
			Hash:  "someblockhash",
		},
		TransactionIdentifier: &rTypes.TransactionIdentifier{Hash: "somehash"},
	}

	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(exampleBlock, mocks.NilError)
	mocks.MTransactionRepository.On("FindByHashInBlock").Return(mocks.NilTransaction, &rTypes.Error{})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	blockService := NewBlockAPIService(commons)

	// when:
	res, e := blockService.BlockTransaction(nil, exampleTransactionRequest)

	// then:
	assert.Nil(t, res)
	assert.NotNil(t, e)
}
