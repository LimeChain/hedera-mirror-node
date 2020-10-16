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

package transaction

import (
	entityid "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/services/encoding"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldSuccessReturnTransactionTableName(t *testing.T) {
	assert.Equal(t, tableNameTransaction, transaction{}.TableName())
}

func TestShouldSuccessReturnTransactionTypesTableName(t *testing.T) {
	assert.Equal(t, tableNameTransactionTypes, transactionType{}.TableName())
}

func TestShouldSuccessReturnTransactionResultsTableName(t *testing.T) {
	assert.Equal(t, tableNameTransactionResults, transactionResult{}.TableName())
}

func TestShouldSuccessReturnRepository(t *testing.T) {
	// given
	gormDbClient, _ := mocks.DatabaseMock(t)

	// when
	result := NewTransactionRepository(gormDbClient)

	// then
	assert.IsType(t, &TransactionRepository{}, result)
	assert.Equal(t, result.dbClient, gormDbClient)
}

func TestShouldSuccessIntsToString(t *testing.T) {
	data := []int64{1, 2, 2394238471841, 2394143718391293}
	expected := "1,2,2394238471841,2394143718391293"

	result := intsToString(data)

	assert.Equal(t, expected, result)
}

func TestShouldFailConstructAccount(t *testing.T) {
	data := int64(-1)
	expected := errors.Errors[errors.InternalServerError]

	result, err := constructAccount(data)

	assert.Nil(t, result)
	assert.Equal(t, expected, err)
}

func TestShouldSuccessConstructAccount(t *testing.T) {
	// given
	data := int64(5)
	expected := &types.Account{EntityId: entityid.EntityId{
		ShardNum:  0,
		RealmNum:  0,
		EntityNum: 5,
	}}

	// when
	result, err := constructAccount(data)

	// then
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
