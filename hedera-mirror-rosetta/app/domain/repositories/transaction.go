package repositories

import (
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
)

// TransactionRepository Interface that all TransactionRepository structs must implement
type TransactionRepository interface {
	FindById(id string) *types.Transaction
	FindBetween(start int64, end int64) []*types.Transaction
}
