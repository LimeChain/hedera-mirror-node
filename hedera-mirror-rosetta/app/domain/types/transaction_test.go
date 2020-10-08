package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToRosettaTransaction(t *testing.T) {
	exampleAccount, _ := AccountFromString("0.0.0")
	exampleAmount := &Amount{Value: int64(400)}
	exampleOperation := &Operation{
		Index:   1,
		Type:    "transfer",
		Status:  "pending",
		Account: exampleAccount,
		Amount:  exampleAmount,
	}
	operations := make([]*Operation, 1)
	operations[0] = exampleOperation
	exampleTransaction := &Transaction{
		Hash:       "somehash",
		Operations: operations,
	}

	rosettaTransaction := exampleTransaction.ToRosettaTransaction()

	assert.Len(t, rosettaTransaction.Operations, 1)
	assert.Equal(t, exampleOperation.ToRosettaOperation(), rosettaTransaction.Operations[0])
	assert.Equal(t, exampleTransaction.Hash, rosettaTransaction.TransactionIdentifier.Hash)
}
