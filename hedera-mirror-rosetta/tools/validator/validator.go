package validator

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
)

func ValidateOperations(operations []*types.Operation) (*string, *types.Error) {
	typeOperation := operations[0].Type
	length := len(operations)

	for i := 1; i < length; i++ {
		if operations[i].Type != typeOperation {
			return nil, errors.Errors[errors.OnlyOneOperationTypePresent]
		}
	}

	return &typeOperation, nil
}
