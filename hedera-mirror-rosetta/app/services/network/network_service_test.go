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

package network

import (
	"github.com/coinbase/rosetta-sdk-go/server"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/services/base"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tests/mocks/repository"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/maphelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNetworkAPIService(t *testing.T) {
	repository.Setup()
	baseService := base.NewBaseService(repository.MBlockRepository, repository.MTransactionRepository)
	networkService := networkAPIService(baseService)

	assert.IsType(t, &NetworkAPIService{}, networkService)
}

func networkAPIService(base base.BaseService) server.NetworkAPIServicer {
	return NewNetworkAPIService(
		base,
		repository.MAddressBookEntryRepository,
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
}

func TestNetworkList(t *testing.T) {
	// given:
	expectedResult := &rTypes.NetworkListResponse{
		NetworkIdentifiers: []*rTypes.NetworkIdentifier{
			{
				Blockchain: "SomeBlockchain",
				Network:    "SomeNetwork",
				SubNetworkIdentifier: &rTypes.SubNetworkIdentifier{
					Network:  "SomeSubNetwork",
					Metadata: nil,
				},
			},
		},
	}

	repository.Setup()

	commons := base.NewBaseService(repository.MBlockRepository, repository.MTransactionRepository)
	networkService := networkAPIService(commons)

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

	repository.Setup()
	repository.MTransactionRepository.On("Statuses").Return(map[int]string{1: "Pending"}, repository.NilError)
	repository.MTransactionRepository.On("TypesAsArray").Return([]string{"Transfer"}, repository.NilError)

	commons := base.NewBaseService(repository.MBlockRepository, repository.MTransactionRepository)
	networkService := networkAPIService(commons)

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

	repository.Setup()
	repository.MBlockRepository.On("RetrieveGenesis").Return(exampleGenesisBlock, repository.NilError)
	repository.MBlockRepository.On("RetrieveLatest").Return(exampleLatestBlock, repository.NilError)
	repository.MAddressBookEntryRepository.On("Entries").Return(exampleEntries, repository.NilError)

	commons := base.NewBaseService(repository.MBlockRepository, repository.MTransactionRepository)
	networkService := networkAPIService(commons)

	// when:
	res, e := networkService.NetworkStatus(nil, nil)

	// then:
	assert.Equal(t, expectedResult, res)
	assert.Nil(t, e)
}

func TestNetworkStatusThrowsWhenRetrieveGenesisFails(t *testing.T) {
	// given:
	repository.Setup()
	repository.MBlockRepository.On("RetrieveGenesis").Return(repository.NilBlock, &rTypes.Error{})

	commons := base.NewBaseService(repository.MBlockRepository, repository.MTransactionRepository)
	networkService := networkAPIService(commons)

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

	repository.Setup()
	repository.MBlockRepository.On("RetrieveGenesis").Return(exampleGenesisBlock, repository.NilError)
	repository.MBlockRepository.On("RetrieveLatest").Return(repository.NilBlock, &rTypes.Error{})

	commons := base.NewBaseService(repository.MBlockRepository, repository.MTransactionRepository)
	networkService := networkAPIService(commons)

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

	repository.Setup()
	repository.MBlockRepository.On("RetrieveGenesis").Return(exampleGenesisBlock, repository.NilError)
	repository.MBlockRepository.On("RetrieveLatest").Return(exampleLatestBlock, repository.NilError)
	repository.MAddressBookEntryRepository.On("Entries").Return(repository.NilEntries, &rTypes.Error{})

	commons := base.NewBaseService(repository.MBlockRepository, repository.MTransactionRepository)
	networkService := networkAPIService(commons)

	// when:
	res, e := networkService.NetworkStatus(nil, nil)

	// then:
	assert.Nil(t, res)
	assert.NotNil(t, e)
}
