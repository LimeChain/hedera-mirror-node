package services

import (
	"context"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/repositories"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
)

// AccountAPIService implements the server.AccountAPIServicer interface.
type AccountAPIService struct {
	Commons
	accountRepo repositories.AccountRepository
}

// NewAccountAPIService creates a new instance of a AccountAPIService.
func NewAccountAPIService(commons Commons, accountRepo repositories.AccountRepository) *AccountAPIService {
	return &AccountAPIService{
		Commons:     commons,
		accountRepo: accountRepo,
	}
}

// AccountBalance implements the /account/balance endpoint.
func (a *AccountAPIService) AccountBalance(ctx context.Context, request *rTypes.AccountBalanceRequest) (*rTypes.AccountBalanceResponse, *rTypes.Error) {
	bIdentifier := &rTypes.BlockIdentifier{}
	balance := &types.Amount{}

	if request.BlockIdentifier != nil {
		block, err := a.RetrieveBlock(request.BlockIdentifier)
		if err != nil {
			return nil, err
		}
		balance, err = a.accountRepo.RetrieveBalanceAtBlock(request.AccountIdentifier.Address, block.ConsensusEnd)
		if err != nil {
			return nil, err
		}

		bIdentifier.Index = *request.BlockIdentifier.Index
		bIdentifier.Hash = *request.BlockIdentifier.Hash
	} else {
		// TODO populate latest block index and hash
	}

	return &rTypes.AccountBalanceResponse{
		BlockIdentifier: bIdentifier,
		Balances:        []*rTypes.Amount{balance.ToRosettaAmount()},
	}, nil
}
