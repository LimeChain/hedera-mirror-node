package repositories

import (
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
)

// TransactionRepository Interface that all TransactionRepository structs must implement
type TransactionRepository interface {
	FindByTimestamp(timestamp int64) *types.Transaction
	FindBetween(start int64, end int64) ([]*types.Transaction, errors.Error)
	GetTypes() map[int]string
	GetStatuses() map[int]string
}
