package mocks

import (
	"fmt"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/stretchr/testify/mock"
)

type MockBlockRepository struct {
	mock.Mock
}

func (m MockBlockRepository) FindByIndex(index int64) (*types.Block, *rTypes.Error) {
	panic("implement me")
}

func (m MockBlockRepository) FindByHash(hash string) (*types.Block, *rTypes.Error) {
	fmt.Printf("Finding by hash [%s]", hash)
	args := m.Called()
	return args.Get(0).(*types.Block), args.Get(1).(*rTypes.Error)
}

func (m MockBlockRepository) FindByIdentifier(index int64, hash string) (*types.Block, *rTypes.Error) {
	fmt.Println("Finding by identifier")
	args := m.Called()
	return args.Get(0).(*types.Block), args.Get(1).(*rTypes.Error)
}

func (m MockBlockRepository) RetrieveGenesis() (*types.Block, *rTypes.Error) {
	panic("implement me")
}

func (m MockBlockRepository) RetrieveLatest() (*types.Block, *rTypes.Error) {
	fmt.Println("Retrieving latest")
	args := m.Called()
	return args.Get(0).(*types.Block), args.Get(1).(*rTypes.Error)
}
