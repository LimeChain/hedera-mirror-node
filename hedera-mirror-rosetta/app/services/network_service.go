package services

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/repositories"
)

type NetworkService struct {
	network         *types.NetworkIdentifier
	blockRepo       repositories.BlockRepository
	transactionRepo repositories.TransactionRepository
	version         *types.Version
}

func (n NetworkService) NetworkList(ctx context.Context, request *types.MetadataRequest) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			n.network,
		},
	}, nil
}

func (n NetworkService) NetworkOptions(ctx context.Context, request *types.NetworkRequest) (*types.NetworkOptionsResponse, *types.Error) {
	return &types.NetworkOptionsResponse{
		Version: n.version,
		Allow: &types.Allow{
			OperationStatuses:       nil,
			OperationTypes:          n.transactionRepo.GetTypesAsArray(),
			Errors:                  nil,
			HistoricalBalanceLookup: false,
		},
	}, nil
}

func (n NetworkService) NetworkStatus(ctx context.Context, request *types.NetworkRequest) (*types.NetworkStatusResponse, *types.Error) {
	genesisBlock, err := n.blockRepo.RetrieveGenesis()
	if err != nil {
		return nil, &types.Error{}
	}

	latestBlock, err := n.blockRepo.RetrieveLatest()
	if err != nil {
		return nil, &types.Error{}
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: latestBlock.ConsensusStart,
			Hash:  latestBlock.Hash,
		},
		CurrentBlockTimestamp: latestBlock.ConsensusStart,
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: genesisBlock.ConsensusStart,
			Hash:  genesisBlock.Hash,
		},
		Peers: nil,
	}, nil
}

func NewNetworkAPIService(network *types.NetworkIdentifier, version *types.Version, blockRepo repositories.BlockRepository, transactionRepo repositories.TransactionRepository) server.NetworkAPIServicer {
	return &NetworkService{
		network:         network,
		version:         version,
		blockRepo:       blockRepo,
		transactionRepo: transactionRepo,
	}
}
