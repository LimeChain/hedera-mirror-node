package services

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/mocks"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/maphelper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkList(t *testing.T) {
	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&types.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &types.SubNetworkIdentifier{
				Network:  "SomeSumNetwork",
				Metadata: nil,
			},
		},
		&types.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)

	expectedResult := &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			{
				Blockchain: "SomeBlockchain",
				Network:    "SomeNetwork",
				SubNetworkIdentifier: &types.SubNetworkIdentifier{
					Network:  "SomeSumNetwork",
					Metadata: nil,
				},
			},
		},
	}

	res, e := networkService.NetworkList(nil, nil)

	assert.Equal(t, expectedResult, res)
	assert.Nil(t, e)
}

func TestNetworkOptions(t *testing.T) {
	mocks.Setup()
	mocks.MTransactionRepository.On("Statuses").Return(map[int]string{1: "Pending"})
	mocks.MTransactionRepository.On("TypesAsArray").Return([]string{"Transfer"})

	expectedResult := &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				&types.OperationStatus{
					Status:     "Pending",
					Successful: true,
				},
			},
			OperationTypes:          []string{"Transfer"},
			Errors:                  maphelper.GetErrorValuesFromStringErrorMap(errors.Errors),
			HistoricalBalanceLookup: true,
		},
	}

	commons := NewCommons(mocks.MBlockRepository, mocks.MTransactionRepository)
	networkService := NewNetworkAPIService(
		commons,
		mocks.MAddressBookEntryRepository,
		&types.NetworkIdentifier{
			Blockchain: "SomeBlockchain",
			Network:    "SomeNetwork",
			SubNetworkIdentifier: &types.SubNetworkIdentifier{
				Network:  "SomeSubNetwork",
				Metadata: nil,
			},
		},
		&types.Version{
			RosettaVersion:    "1",
			NodeVersion:       "1",
			MiddlewareVersion: nil,
			Metadata:          nil,
		},
	)
	res, e := networkService.NetworkOptions(nil, nil)

	assert.Equal(t, expectedResult.Version, res.Version)
	assert.Equal(t, expectedResult.Allow.HistoricalBalanceLookup, res.Allow.HistoricalBalanceLookup)
	assert.Equal(t, expectedResult.Allow.OperationStatuses, res.Allow.OperationStatuses)
	assert.Equal(t, expectedResult.Allow.OperationTypes, res.Allow.OperationTypes)
	assert.Equal(t, len(expectedResult.Allow.Errors), len(res.Allow.Errors))
	assert.Nil(t, e)
}
