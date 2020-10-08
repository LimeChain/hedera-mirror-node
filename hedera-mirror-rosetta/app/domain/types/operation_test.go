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
	"strconv"
	"testing"
)

func TestToRosettaOperation(t *testing.T) {
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

	// when:
	rosettaOperation := exampleOperation.ToRosettaOperation()

	// then:
	assert.Equal(t, strconv.FormatInt(exampleOperation.Amount.Value, 10), rosettaOperation.Amount.Value)
	assert.Equal(t, exampleAccount.ToRosettaAccount(), rosettaOperation.Account)
	assert.Equal(t, exampleOperation.Status, rosettaOperation.Status)
	assert.Equal(t, exampleOperation.Type, rosettaOperation.Type)
	assert.Equal(t, exampleOperation.Index, rosettaOperation.OperationIdentifier.Index)
}
