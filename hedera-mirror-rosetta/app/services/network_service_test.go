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
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/mocks"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/maphelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkList(t *testing.T) {
	// given:
	expectedResult := &rTypes.NetworkListResponse{
		NetworkIdentifiers: []*rTypes.NetworkIdentifier{
			{
				Blockchain: "SomeBlockchain",
				Network:    "SomeNetwork",
				SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
					Network:  "SomeSumNetwork",
					Metadata: nil,
				},
			},
		},
	}

	mocks.Setup()

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&rTypes.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "SomeSumNetwork",
				Metadata: nil,
			},
		},
		&rTypes.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)

	// when:
	res, e := networkService.NetworkList(nil, nil)

	// then:
	assert.Equal(t, expectedResult, res)
	assert.Nil(t, e)
}

func TestNetworkOptions(t *testing.T) {
	// given:
	expectedResult := &rTypes.NetworkOptionsResponse{
		Version: &rTypes.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
		Allow: &rTypes.Allow{
			OperationStatuses: []*rTypes.OperationStatus{
				{
					Status:     "Pending",
					Successful: true,
				},
			},
			OperationTypes:          []string{"Transfer"},
			Errors:                  maphelper.GetErrorValuesFromStringErrorMap(errors.Errors),
			HistoricalBalanceLookup: true,
		},
	}

	mocks.Setup()
	mocks.MTransactionRepository.On("Statuses").Return(map[int]string{1: "Pending"})
	mocks.MTransactionRepository.On("TypesAsArray").Return([]string{"Transfer"})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&rTypes.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "SomeSubNetwork",
				Metadata: nil,
			},
		},
		&rTypes.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)

	// when:
	res, e := networkService.NetworkOptions(nil, nil)

	// then:
	assert.Equal(t, expectedResult.Version, res.Version)
	assert.Equal(t, expectedResult.Allow.HistoricalBalanceLookup, res.Allow.HistoricalBalanceLookup)
	assert.Equal(t, expectedResult.Allow.OperationStatuses, res.Allow.OperationStatuses)
	assert.Equal(t, expectedResult.Allow.OperationTypes, res.Allow.OperationTypes)
	assert.Equal(t, len(expectedResult.Allow.Errors), len(res.Allow.Errors))
	assert.Nil(t, e)
}

func TestNetworkStatus(t *testing.T) {
	// given:
	exampleGenesisBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         0,
		ParentHash:          "",
	}

	exampleLatestBlock := &types.Block{
		Index:               2,
		Hash:                "0x1323jsjs",
		ConsensusStartNanos: 40000000,
		ConsensusEndNanos:   70000000,
		ParentIndex:         1,
		ParentHash:          "0x123jsjs",
	}

	exampleEntries := &types.AddressBookEntries{Entries: []*types.AddressBookEntry{}}

	expectedResult := &rTypes.NetworkStatusResponse{
		CurrentBlockIdentifier: &rTypes.BlockIdentifier{
			Index: 2,
			Hash:  "0x1323jsjs",
		},
		CurrentBlockTimestamp: 40,
		GenesisBlockIdentifier: &rTypes.BlockIdentifier{
			Index: 1,
			Hash:  "0x123jsjs",
		},
		Peers: []*rTypes.Peer{},
	}

	mocks.Setup()
	mocks.MBlockRepository.On("RetrieveGenesis").Return(exampleGenesisBlock, mocks.NilError)
	mocks.MBlockRepository.On("RetrieveLatest").Return(exampleLatestBlock, mocks.NilError)
	mocks.MAddressBookEntryRepository.On("Entries").Return(exampleEntries, mocks.NilError)

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&rTypes.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "SomeSumNetwork",
				Metadata: nil,
			},
		},
		&rTypes.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)

	// when:
	res, e := networkService.NetworkStatus(nil, nil)

	// then:
	assert.Equal(t, expectedResult, res)
	assert.Nil(t, e)
}

func TestNetworkStatusThrowsWhenRetrieveGenesisFails(t *testing.T) {
	// given:
	mocks.Setup()
	mocks.MBlockRepository.On("RetrieveGenesis").Return(mocks.NilBlock, &rTypes.Error{})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&rTypes.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "SomeSumNetwork",
				Metadata: nil,
			},
		},
		&rTypes.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)

	// when:
	res, e := networkService.NetworkStatus(nil, nil)

	// then
	assert.Nil(t, res)
	assert.NotNil(t, e)
}

func TestNetworkStatusThrowsWhenRetrieveLatestFails(t *testing.T) {
	// given:
	exampleGenesisBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         0,
		ParentHash:          "",
	}

	mocks.Setup()
	mocks.MBlockRepository.On("RetrieveGenesis").Return(exampleGenesisBlock, mocks.NilError)
	mocks.MBlockRepository.On("RetrieveLatest").Return(mocks.NilBlock, &rTypes.Error{})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&rTypes.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "SomeSumNetwork",
				Metadata: nil,
			},
		},
		&rTypes.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)

	// when:
	res, e := networkService.NetworkStatus(nil, nil)

	// then:
	assert.Nil(t, res)
	assert.NotNil(t, e)
}

func TestNetworkStatusThrowsWhenEntriesFail(t *testing.T) {
	// given:
	exampleGenesisBlock := &types.Block{
		Index:               1,
		Hash:                "0x123jsjs",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   20000000,
		ParentIndex:         0,
		ParentHash:          "",
	}

	exampleLatestBlock := &types.Block{
		Index:               2,
		Hash:                "0x1323jsjs",
		ConsensusStartNanos: 40000000,
		ConsensusEndNanos:   70000000,
		ParentIndex:         1,
		ParentHash:          "0x123jsjs",
	}

	mocks.Setup()
	mocks.MBlockRepository.On("RetrieveGenesis").Return(exampleGenesisBlock, mocks.NilError)
	mocks.MBlockRepository.On("RetrieveLatest").Return(exampleLatestBlock, mocks.NilError)
	mocks.MAddressBookEntryRepository.On("Entries").Return(mocks.NilEntries, &rTypes.Error{})

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&rTypes.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
				Network:  "SomeSumNetwork",
				Metadata: nil,
			},
		},
		&rTypes.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)

	// when:
	res, e := networkService.NetworkStatus(nil, nil)

	// then:
	assert.Nil(t, res)
	assert.NotNil(t, e)
}
