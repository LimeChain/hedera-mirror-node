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

func exampleOperation() *Operation {
	return &Operation{
		Index:   1,
		Type:    "transfer",
		Status:  "pending",
		Account: exampleAccount(),
		Amount:  exampleAmount(),
	}
}

func exampleOperations() []*Operation {
	return []*Operation{
		exampleOperation(),
	}
}

func expectedOperationIdentifier() *types.OperationIdentifier {
	return &types.OperationIdentifier{
		Index:        1,
		NetworkIndex: nil,
	}
}

func expectedRelatedOperations() []*types.OperationIdentifier {
	return []*types.OperationIdentifier{}
}

func expectedOperation() *types.Operation {
	return &types.Operation{
		OperationIdentifier: expectedOperationIdentifier(),
		RelatedOperations:   expectedRelatedOperations(),
		Type:                "transfer",
		Status:              "pending",
		Account:             expectedAccount(),
		Amount:              expectedAmount(),
	}
}

func TestToRosettaOperation(t *testing.T) {
	// when:
	rosettaOperation := exampleOperation().ToRosetta()

	// then:
	assert.Equal(t, expectedOperation(), rosettaOperation)
}
