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
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exampleBlock() *Block {
	return &Block{
		Index:               2,
		Hash:                "somehash",
		ConsensusStartNanos: 10000000,
		ConsensusEndNanos:   12300000,
		ParentIndex:         1,
		ParentHash:          "someparenthash",
		Transactions:        exampleTransactions(),
	}
}

func expectedBlockIdentifier() *types.BlockIdentifier {
	return &types.BlockIdentifier{
		Index: 2,
		Hash:  "0xsomehash",
	}
}

func expectedParentBlockIdentifier() *types.BlockIdentifier {
	return &types.BlockIdentifier{
		Index: 1,
		Hash:  "0xsomeparenthash",
	}
}

func expectedBlock() *types.Block {
	return &types.Block{
		BlockIdentifier:       expectedBlockIdentifier(),
		ParentBlockIdentifier: expectedParentBlockIdentifier(),
		Timestamp:             int64(10),
		Transactions:          []*types.Transaction{expectedTransaction()},
	}
}

func TestToRosettaBlock(t *testing.T) {
	// when:
	rosettaBlockResult := exampleBlock().ToRosetta()

	// then:
	assert.Equal(t, expectedBlock(), rosettaBlockResult)
}

func TestGetTimestampMillis(t *testing.T) {
	// given:
	exampleBlock := exampleBlock()
	exampleBlock.ConsensusStartNanos = 10000000

	// when:
	resultMillis := exampleBlock.GetTimestampMillis()

	// then:
	assert.Equal(t, int64(10), resultMillis)
}
