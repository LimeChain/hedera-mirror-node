package validator

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateOperationsSum(t *testing.T) {
	// given:
	operationDummy := newOperationDummy("100")
	operationDummy2 := newOperationDummy("-100")

	testData := []*types.Operation{
		operationDummy,
		operationDummy2,
	}

	expectedError := errors.Errors[errors.InvalidOperationsTotalAmount]

	// when:
	result := ValidateOperationsSum(testData)

	// then:
	assert.Equal(t, nil, result)

	// and:
	testData = append(testData, operationDummy2)

	// then:
	result = ValidateOperationsSum(testData)
	assert.Equal(t, expectedError, result)
}

func newOperationDummy(amount string) *types.Operation {
	return &types.Operation{
		Amount: &types.Amount{
			Value: amount,
		},
	}
}
