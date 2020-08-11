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
	network         *rTypes.NetworkIdentifier
	blockRepo       repositories.BlockRepository
	transactionRepo repositories.TransactionRepository
}

// NewBlockAPIService creates a new instance of a BlockAPIService.
func NewBlockAPIService(network *rTypes.NetworkIdentifier, blockRepo repositories.BlockRepository, transactionRepo repositories.TransactionRepository) server.BlockAPIServicer {
	return &BlockAPIService{
		network:         network,
		blockRepo:       blockRepo,
		transactionRepo: transactionRepo,
	}
}

// Block implements the /block endpoint.
func (s *BlockAPIService) Block(ctx context.Context, request *rTypes.BlockRequest) (*rTypes.BlockResponse, *rTypes.Error) {
	var block = &types.Block{}
	var err *rTypes.Error
	if request.BlockIdentifier.Hash != nil && request.BlockIdentifier.Index != nil {
		block, err = s.blockRepo.FindByIdentifier(*request.BlockIdentifier.Index, *request.BlockIdentifier.Hash)
	} else if request.BlockIdentifier.Hash == nil {
		block, err = s.blockRepo.FindByIndex(*request.BlockIdentifier.Index)
	} else if request.BlockIdentifier.Index == nil {
		block, err = s.blockRepo.FindByHash(*request.BlockIdentifier.Hash)
	} else {
		block, err = s.blockRepo.RetrieveLatest()
	}

	if err != nil {
		return nil, err
	}

	tArray, err := s.transactionRepo.FindBetween(block.ConsensusStart, block.ConsensusEnd)
	if err != nil {
		return nil, err
	}

	block.Transactions = tArray
	rBlock := block.ToRosettaBlock()
	return &rTypes.BlockResponse{
		Block: rBlock,
	}, nil
}

// BlockTransaction implements the /block/transaction endpoint.
func (s *BlockAPIService) BlockTransaction(
	ctx context.Context,
	request *rTypes.BlockTransactionRequest,
) (*rTypes.BlockTransactionResponse, *rTypes.Error) {
	block, err := s.blockRepo.FindByIdentifier(request.BlockIdentifier.Index, request.BlockIdentifier.Hash)
	if err != nil {
		return nil, err
	}

	transaction, err := s.transactionRepo.FindByIdentifierInBlock(request.TransactionIdentifier.Hash, block.ConsensusStart, block.ConsensusEnd)
	if err != nil {
		return nil, err
	}

	rTransaction := transaction.ToRosettaTransaction()
	return &rTypes.BlockTransactionResponse{
		Transaction: rTransaction,
	}, nil
}
