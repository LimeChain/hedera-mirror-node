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

package account

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	accountString      = "0.0.5"
	consensusTimestamp = int64(2)
	dbAccountBalance   = &accountBalance{
		ConsensusTimestamp: 1,
		Balance:            10,
		AccountId:          int64(5),
	}
	expectedAmount = &types.Amount{Value: 10}
)

func TestShouldReturnValidAccountBalanceTableName(t *testing.T) {
	assert.Equal(t, tableNameAccountBalance, accountBalance{}.TableName())
}

func TestShouldSuccessReturnValidRepository(t *testing.T) {
	// given
	gormDbClient, _ := mocks.DatabaseMock(t)

	// when
	result := NewAccountRepository(gormDbClient)

	// then
	assert.IsType(t, &AccountRepository{}, result)
	assert.Equal(t, result.dbClient, gormDbClient)
}

func TestShouldFailRetrieveBalanceAtBlockDueToInvalidAddress(t *testing.T) {
	// given
	abr, _, mock := setupRepository(t)
	defer abr.dbClient.DB().Close()

	invalidAddressString := "0.0.a"

	// when
	result, err := abr.RetrieveBalanceAtBlock(invalidAddressString, consensusTimestamp)

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestShouldSuccessRetrieveBalanceAtBlock(t *testing.T) {
	abr, columns, mock := setupRepository(t)
	defer abr.dbClient.DB().Close()

	mock.ExpectQuery(latestBalanceBeforeConsensus).
		WithArgs("0.0.5", 10).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow(1, 10, 5))

	mock.ExpectQuery(balanceChangeBetween).WithArgs(1, 10, 5).WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow(1, 10, 5))

	// when
	result, err := abr.RetrieveBalanceAtBlock("0.0.5", consensusTimestamp)

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func setupRepository(t *testing.T) (*AccountRepository, []string, sqlmock.Sqlmock) {
	gormDbClient, mock := mocks.DatabaseMock(t)

	columns := mocks.GetFieldsNamesToSnakeCase(accountBalance{})

	aber := NewAccountRepository(gormDbClient)
	return aber, columns, mock
}
