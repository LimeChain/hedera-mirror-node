/*-
 * ‌
 * Hedera Mirror Node
 *
 * Copyright (C) 2019 - 2020 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ‍
 */

package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToRosettaBlock(t *testing.T) {
	// given:
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
	transactions := make([]*Transaction, 1)
	transactions[0] = exampleTransaction
	exampleBlock := &Block{
		Index:               2,
		Hash:                "somehash",
		ConsensusStartNanos: 10000000,
		ConsensusEndNanos:   12300000,
		ParentIndex:         1,
		ParentHash:          "someparenthash",
		Transactions:        transactions,
	}

	// when:
	rosettaBlockResult := exampleBlock.ToRosettaBlock()

	// then:
	assert.Equal(t, int64(10), rosettaBlockResult.Timestamp)
	assert.Equal(t, "0xsomehash", rosettaBlockResult.BlockIdentifier.Hash)
	assert.Equal(t, "0xsomeparenthash", rosettaBlockResult.ParentBlockIdentifier.Hash)
	assert.Len(t, rosettaBlockResult.Transactions, 1)
	assert.Equal(t, exampleTransaction.Hash, rosettaBlockResult.Transactions[0].TransactionIdentifier.Hash)
}

func TestGetTimestampMillis(t *testing.T) {
	// given:
	exampleBlock := &Block{
		ConsensusStartNanos: 10000000,
	}

	// when:
	resultMillis := exampleBlock.GetTimestampMillis()

	// then:
	assert.Equal(t, int64(10), resultMillis)
}
