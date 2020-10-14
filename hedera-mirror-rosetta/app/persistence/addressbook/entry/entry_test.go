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
		Entries: []*types.AddressBookEntry{expectedAddressBookEntry},
	}
)

func TestShouldReturnCorrectAddressBookEntries(t *testing.T) {
	// given
	aber, columns, mock := setupRepository(t)
	defer aber.dbClient.DB().Close()

	mock.ExpectQuery(latestAddressBookEntries).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow(1, 1, "127.0.0.1", 0, "0.0.5", nil, nil, nil, nil))

	// when
	result, err := aber.Entries()

	// assert
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err.Message)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Len(t, expectedResult.Entries, 1)
	assert.Equal(t, expectedResult.Entries[0].Metadata, result.Entries[0].Metadata)
	assert.Equal(t, expectedResult.Entries[0].PeerId, result.Entries[0].PeerId)
}

func TestShouldFailReturnEntries(t *testing.T) {
	// given
	aber, columns, mock := setupRepository(t)
	defer aber.dbClient.DB().Close()

	mock.ExpectQuery(latestAddressBookEntries).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow(1, 1, "127.0.0.1", 0, "0.0.a", nil, nil, nil, nil))

	// when
	result, err := aber.Entries()

	// assert
	if result != nil {
		t.Error("Expected no result.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	assert.Equal(t, err.Message, errors.InternalServerError)
}

func TestShouldReturnProperPeerId(t *testing.T) {
	// given
	abe := addressBookEntry{
		Memo: peerId.String(),
	}

	// when
	result, err := abe.getPeerId()

	// then
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err.Message)
	}

	assert.Equal(t, peerId, result)
}

func TestShouldFailReturningProperPeerId(t *testing.T) {
	// given
	abe := addressBookEntry{
		Memo: "0.0.a",
	}

	// when
	result, err := abe.getPeerId()

	// then
	if result != nil {
		t.Errorf("Expected no result, got '%s'", result)
	}

	assert.Equal(t, err.Message, errors.InternalServerError)
}

func setupRepository(t *testing.T) (*AddressBookEntryRepository, []string, sqlmock.Sqlmock) {
	gormDbClient, mock := mocks.DatabaseMock(t)

	columns := mocks.GetFieldsToSnakeCase(addressBookEntry{})

	aber := NewAddressBookEntryRepository(gormDbClient)
	return aber, columns, mock
}
