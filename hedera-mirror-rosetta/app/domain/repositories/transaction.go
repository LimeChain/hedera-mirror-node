package repositories

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
)

// TransactionRepository Interface that all TransactionRepository structs must implement
type TransactionRepository interface {
	FindByHashInBlock(identifier string, consensusStart int64, consensusEnd int64) (*types.Transaction, *rTypes.Error)
	FindBetween(start int64, end int64) ([]*types.Transaction, *rTypes.Error)
	Types() (map[int]string, *rTypes.Error)
	TypesAsArray() ([]string, *rTypes.Error)
	Statuses() (map[int]string, *rTypes.Error)
}
