package types

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestToRosettaOperation(t *testing.T) {
	exampleAccount, _ := AccountFromString("0.0.0")
	exampleAmount := &Amount{Value: int64(400)}
	exampleOperation := &Operation{
		Index:   1,
		Type:    "transfer",
		Status:  "pending",
		Account: exampleAccount,
		Amount:  exampleAmount,
	}

	rosettaOperation := exampleOperation.ToRosettaOperation()

	assert.Equal(t, strconv.FormatInt(exampleOperation.Amount.Value, 10), rosettaOperation.Amount.Value)
	assert.Equal(t, exampleAccount.ToRosettaAccount(), rosettaOperation.Account)
	assert.Equal(t, exampleOperation.Status, rosettaOperation.Status)
	assert.Equal(t, exampleOperation.Type, rosettaOperation.Type)
	assert.Equal(t, exampleOperation.Index, rosettaOperation.OperationIdentifier.Index)
}
