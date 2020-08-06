package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/server"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/repositories"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
)

// BlockAPIService implements the server.BlockAPIServicer interface.
type BlockAPIService struct {
	network   *rTypes.NetworkIdentifier
	blockRepo repositories.BlockRepository
}

// NewBlockAPIService creates a new instance of a BlockAPIService.
func NewBlockAPIService(network *rTypes.NetworkIdentifier, blockRepo repositories.BlockRepository) server.BlockAPIServicer {
	return &BlockAPIService{
		network:   network,
		blockRepo: blockRepo,
	}
}

// Block implements the /block endpoint.
func (s *BlockAPIService) Block(ctx context.Context, request *rTypes.BlockRequest) (*rTypes.BlockResponse, *rTypes.Error) {
	var block = &types.Block{}
	if request.BlockIdentifier.Hash != nil && request.BlockIdentifier.Index != nil {
		block = s.blockRepo.FindByIndentifier(*request.BlockIdentifier.Index, *request.BlockIdentifier.Hash)
	} else if request.BlockIdentifier.Hash == nil {
		block = s.blockRepo.FindByIndex(*request.BlockIdentifier.Index)
	} else if request.BlockIdentifier.Index == nil {
		block = s.blockRepo.FindByHash(*request.BlockIdentifier.Hash)
	} else {
		block = s.blockRepo.RetrieveLatest()
	}
	// TODO what is it does not exist?

	rBlock := block.ToRosettaBlock()

	return &rTypes.BlockResponse{
		Block: &rTypes.Block{},
	}, nil
}

// BlockTransaction implements the /block/transaction endpoint.
func (s *BlockAPIService) BlockTransaction(
	ctx context.Context,
	request *rTypes.BlockTransactionRequest,
) (*rTypes.BlockTransactionResponse, *rTypes.Error) {

	// TODO
	return &rTypes.BlockTransactionResponse{}, nil
}
