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
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestNewAccountAPIService(t *testing.T) {
	mocks.Setup()
	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	accountService := NewAccountAPIService(commons, mocks.MAccountRepository)

	assert.IsType(t, &AccountAPIService{}, accountService)
}

func TestAccountBalance(t *testing.T) {
	// given:
	exampleBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}

	expectedAccountBalanceResponse := &rTypes.AccountBalanceResponse{
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

	exampleAmount := &types.Amount{
		Value: int64(1000),
	}

	mocks.Setup()
	mocks.MBlockRepository.On("RetrieveLatest").Return(exampleBlock, mocks.NilError)
	mocks.MAccountRepository.On("RetrieveBalanceAtBlock").Return(exampleAmount, mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	accountService := NewAccountAPIService(commons, mocks.MAccountRepository)

	// when:
	actualResult, e := accountService.AccountBalance(nil, exampleRequest(false))

	// then:
	assert.Equal(t, expectedAccountBalanceResponse, actualResult)
	assert.Nil(t, e)
}

func TestAccountBalanceFullData(t *testing.T) {
	// given:
	exampleBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         2,
		ParentHash:          "0xparenthash",
	}

	expectedAccountBalanceResponse := &rTypes.AccountBalanceResponse{
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

	exampleAmount := &types.Amount{
		Value: int64(1000),
	}

	mocks.Setup()
	mocks.MBlockRepository.On("FindByIdentifier").Return(exampleBlock, mocks.NilError)
	mocks.MAccountRepository.On("RetrieveBalanceAtBlock").Return(exampleAmount, mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	accountService := NewAccountAPIService(commons, mocks.MAccountRepository)

	// when:
	actualResult, e := accountService.AccountBalance(nil, exampleRequest(true))

	//then:
	assert.Equal(t, expectedAccountBalanceResponse, actualResult)
	assert.Nil(t, e)
	mocks.MBlockRepository.AssertNotCalled(t, "RetrieveLatest")
}

func TestAccountBalanceThrows(t *testing.T) {
	// given:
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

func TestAccountBalanceThrowsWhenRetrieveBalanceAtBlockFails(t *testing.T) {
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
	mocks.MAccountRepository.On("RetrieveBalanceAtBlock").Return(mocks.NilAmount, &rTypes.Error{})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	accountService := NewAccountAPIService(commons, mocks.MAccountRepository)

	// when:
	actualResult, e := accountService.AccountBalance(nil, exampleRequest(true))

	//then:
	assert.Nil(t, actualResult)
	assert.NotNil(t, e)
}
