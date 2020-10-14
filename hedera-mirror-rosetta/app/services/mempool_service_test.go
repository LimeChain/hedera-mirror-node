package services

import (
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMempool(t *testing.T) {
	res, e := NewMempoolAPIService().Mempool(nil, nil)

	assert.Equal(
		t,
		&rTypes.MempoolResponse{
			TransactionIdentifiers: []*rTypes.TransactionIdentifier{},
		},
		res,
	)

	assert.Nil(t, e)
}

func TestMempoolTransaction(t *testing.T) {
	res, e := NewMempoolAPIService().MempoolTransaction(nil, nil)

	assert.Equal(t, errors.Errors[errors.TransactionNotFound], e)

	assert.Nil(t, res)
}
