package repositories

import (
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
)

// BlockRepository Interface that all BlockRepository structs must implement
type BlockRepository interface {
	FindByIndex(index int64) (*types.Block, errors.Error)
	FindByHash(hash string) (*types.Block, errors.Error)
	FindByIndentifier(index int64, hash string) (*types.Block, errors.Error)
	RetrieveLatest() (*types.Block, errors.Error)
}
