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
	"encoding/hex"
	"github.com/DATA-DOG/go-sqlmock"
	entityid "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/services/encoding"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/test/mocks"
	"github.com/stretchr/testify/assert"
	"reflect"
	"regexp"
	"testing"
)

var (
	transactionTypeColumns   = mocks.GetFieldsNamesToSnakeCase(transactionType{})
	transactionResultColumns = mocks.GetFieldsNamesToSnakeCase(transactionResult{})
	dbTransactionTypes       = []transactionType{
		{
			ProtoID: 12,
			Name:    "CRYPTODELETE",
		},
		{
			ProtoID: 14,
			Name:    "CRYPTOTRANSFER",
		},
	}
	dbTransactionResults = []transactionResult{
		{ProtoID: 0, Result: "OK"},
		{ProtoID: 1, Result: "INVALID_TRANSACTION"},
	}
	tRepoTypes = map[int]string{
		12: "CRYPTODELETE",
		14: "CRYPTOTRANSFER",
	}
	tRepoStatuses = map[int]string{
		0: "OK",
		1: "INVALID_TRANSACTION",
	}
	tRepoTypesAsArray = []string{"CRYPTODELETE", "CRYPTOTRANSFER"}
)

func TestShouldSuccessReturnStatuses(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	rows = willReturnRows(transactionResultColumns, dbTransactionResults)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(rows)

	// when
	result, err := tr.Statuses()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, err)

	assert.Equal(t, tRepoStatuses, result)
}

func TestShouldFailReturnStatuses(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(sqlmock.NewRows(transactionResultColumns))

	// when
	result, err := tr.Statuses()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.NotNil(t, err)

	assert.Nil(t, result)
}

func TestShouldFailReturnTypes(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(sqlmock.NewRows(transactionResultColumns))

	// when
	result, err := tr.Types()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.NotNil(t, err)

	assert.Nil(t, result)
}

func TestShouldFailReturnTypesAsArray(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(sqlmock.NewRows(transactionResultColumns))

	// when
	result, err := tr.TypesAsArray()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.NotNil(t, err)

	assert.Nil(t, result)
}

func TestShouldSuccessReturnTypesAsArray(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	rows = willReturnRows(transactionResultColumns, dbTransactionResults)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(rows)

	// when
	result, err := tr.TypesAsArray()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, err)

	assert.Equal(t, tRepoTypesAsArray, result)
}

func TestShouldSuccessReturnTypes(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	rows = willReturnRows(transactionResultColumns, dbTransactionResults)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(rows)

	// when
	result, err := tr.Types()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, err)

	assert.Equal(t, tRepoTypes, result)
}

func TestShouldSuccessSaveTransactionTypesAndStatuses(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	rows = willReturnRows(transactionResultColumns, dbTransactionResults)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(rows)

	// when
	result := tr.retrieveTransactionTypesAndStatuses()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Nil(t, result)

	assert.Equal(t, tRepoStatuses, tr.statuses)
	assert.Equal(t, tRepoTypes, tr.types)
}

func TestShouldFailReturnTransactionTypesAndStatusesDueToNoResults(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(sqlmock.NewRows(transactionResultColumns))

	// when
	result := tr.retrieveTransactionTypesAndStatuses()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, errors.Errors[errors.OperationStatusesNotFound], result)
}

func TestShouldFailReturnTransactionTypesAndStatusesDueToNoTypes(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(sqlmock.NewRows(transactionTypeColumns))
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(sqlmock.NewRows(transactionResultColumns))

	// when
	result := tr.retrieveTransactionTypesAndStatuses()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, errors.Errors[errors.OperationTypesNotFound], result)
}

func TestShouldSuccessReturnTransactionResults(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionResultColumns, dbTransactionResults)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionResults)).
		WillReturnRows(rows)

	// when
	result := tr.retrieveTransactionResults()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, dbTransactionResults, result)
}

func TestShouldSuccessReturnTransactionTypes(t *testing.T) {
	// given
	tr, mock := setupRepository(t)
	defer tr.dbClient.DB().Close()

	rows := willReturnRows(transactionTypeColumns, dbTransactionTypes)
	mock.ExpectQuery(regexp.QuoteMeta(selectTransactionTypes)).
		WillReturnRows(rows)

	// when
	result := tr.retrieveTransactionTypes()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, dbTransactionTypes, result)
}

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

func TestShouldSuccessGetHashString(t *testing.T) {
	// given
	txStr := "967f26876ad492cc27b4c384dc962f443bcc9be33cbb7add3844bc864de047340e7a78c0fbaf40ab10948dc570bbc25edb505f112d0926dffb65c93199e6d507"
	bytesTx, _ := hex.DecodeString(txStr)
	givenTx := transaction{
		TransactionHash: bytesTx,
	}

	// when
	res := givenTx.getHashString()

	// then
	assert.Equal(t, "0x"+txStr, res)
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

func setupRepository(t *testing.T) (*TransactionRepository, sqlmock.Sqlmock) {
	gormDbClient, mock := mocks.DatabaseMock(t)

	aber := NewTransactionRepository(gormDbClient)
	return aber, mock
}

func willReturnRows(columns []string, data ...interface{}) *sqlmock.Rows {
	converter := sqlmock.NewRows(columns)

	for _, value := range data {
		s := reflect.ValueOf(value)

		for i := 0; i < s.Len(); i++ {
			row := mocks.GetFieldsValuesAsDriverValue(s.Index(i).Interface())
			converter.AddRow(row...)
		}
	}

	return converter
}
