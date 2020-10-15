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

package entry

import (
	"github.com/DATA-DOG/go-sqlmock"
	entityid "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/services/encoding"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	entityId = &entityid.EntityId{
		ShardNum:  0,
		RealmNum:  0,
		EntityNum: 5,
	}
	peerId                   = &types.Account{EntityId: *entityId}
	expectedAddressBookEntry = &types.AddressBookEntry{
		PeerId: peerId,
		Metadata: map[string]interface{}{
			"ip":   "127.0.0.1",
			"port": int32(0),
		},
	}
	expectedResult *types.AddressBookEntries = &types.AddressBookEntries{
		Entries: []*types.AddressBookEntry{expectedAddressBookEntry, expectedAddressBookEntry},
	}
)

func TestShouldSuccessReturnAddressBookEntryTableName(t *testing.T) {
	assert.Equal(t, tableNameAddressBookEntry, addressBookEntry{}.TableName())
}

func TestShouldSuccessReturnRepository(t *testing.T) {
	// given
	gormDbClient, _ := mocks.DatabaseMock(t)

	// when
	result := NewAddressBookEntryRepository(gormDbClient)

	// then
	assert.IsType(t, &AddressBookEntryRepository{}, result)
	assert.Equal(t, result.dbClient, gormDbClient)
}

func TestShouldSuccessReturnAddressBookEntries(t *testing.T) {
	// given
	aber, columns, mock := setupRepository(t)
	defer aber.dbClient.DB().Close()

	mockedRows := sqlmock.NewRows(columns).
		AddRow(1, 1, "127.0.0.1", 0, "0.0.5", nil, nil, nil, nil).
		AddRow(1, 1, "127.0.0.1", 0, "0.0.5", nil, nil, nil, nil)
	mock.ExpectQuery(latestAddressBookEntries).
		WillReturnRows(mockedRows)

	// when
	result, err := aber.Entries()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())

	assert.Len(t, expectedResult.Entries, 2)
	for i, resultEntry := range result.Entries {
		assert.Equal(t, expectedResult.Entries[i].Metadata, resultEntry.Metadata)
		assert.Equal(t, expectedResult.Entries[i].PeerId, resultEntry.PeerId)
	}
	assert.Nil(t, err)
}

func TestShouldFailReturnEntriesDueToInvalidDbData(t *testing.T) {
	// given
	aber, columns, mock := setupRepository(t)
	defer aber.dbClient.DB().Close()

	mock.ExpectQuery(latestAddressBookEntries).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow(1, 1, "127.0.0.1", 0, "0.0.a", nil, nil, nil, nil))

	// when
	result, err := aber.Entries()

	// then
	assert.Nil(t, mock.ExpectationsWereMet())

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.Errors[errors.InternalServerError], err)
}

func TestShouldSuccessReturnPeerId(t *testing.T) {
	// given
	abe := addressBookEntry{
		Memo: peerId.String(),
	}

	// when
	result, err := abe.getPeerId()

	// then
	assert.Equal(t, peerId, result)
	assert.Nil(t, err)
}

func TestShouldFailReturnPeerId(t *testing.T) {
	// given
	abe := addressBookEntry{
		Memo: "0.0.a",
	}

	// when
	result, err := abe.getPeerId()

	// then
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, errors.Errors[errors.InternalServerError], err)
}

func setupRepository(t *testing.T) (*AddressBookEntryRepository, []string, sqlmock.Sqlmock) {
	gormDbClient, mock := mocks.DatabaseMock(t)

	columns := mocks.GetFieldsToSnakeCase(addressBookEntry{})

	aber := NewAddressBookEntryRepository(gormDbClient)
	return aber, columns, mock
}
