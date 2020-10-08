package types

import (
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccountFromEncodedID(t *testing.T) {
	// given:
	var testData = []struct {
		input, shard, realm, number int64
	}{
		{0, 0, 0, 0},
		{10, 0, 0, 10},
		{4294967295, 0, 0, 4294967295},
		{2814792716779530, 10, 10, 10},
		{9223372036854775807, 32767, 65535, 4294967295},
		{9223090561878065152, 32767, 0, 0},
	}

	for _, tt := range testData {
		// when:
		res, _ := NewAccountFromEncodedID(tt.input)

		// then:
		assert.Equal(t, tt.shard, res.Shard)
		assert.Equal(t, tt.realm, res.Realm)
		assert.Equal(t, tt.number, res.Number)
	}
}

func TestNewAccountFromEncodedIDThrows(t *testing.T) {
	// given:
	var testData = struct {
		input, shard, realm, number int64
	}{
		-123, 0, 0, 0,
	}

	// when:
	res, err := NewAccountFromEncodedID(testData.input)

	// then:
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestAccountFromString(t *testing.T) {
	// given:
	var testData = []struct {
		input                string
		shard, realm, number int64
	}{
		{"0.0.0", 0, 0, 0},
		{"0.0.10", 0, 0, 10},
		{"0.0.4294967295", 0, 0, 4294967295},
		{"10.10.10", 10, 10, 10},
		{"32767.65535.4294967295", 32767, 65535, 4294967295},
		{"32767.0.0", 32767, 0, 0},
	}

	for _, tt := range testData {
		// when:
		res, _ := AccountFromString(tt.input)

		// then:
		assert.Equal(t, tt.shard, res.Shard)
		assert.Equal(t, tt.realm, res.Realm)
		assert.Equal(t, tt.number, res.Number)
	}
}

func TestAccountFromStringThrows(t *testing.T) {
	// given:
	var testData = []struct {
		input                string
		shard, realm, number int64
	}{
		{"a.0.0", 0, 0, 0},
		{"0.b.0", 0, 0, 10},
		{"0.0.c", 0, 0, 4294967295},
	}

	var expectedNil *Account = nil

	for _, tt := range testData {
		// when:
		res, err := AccountFromString(tt.input)

		// then:
		assert.Equal(t, expectedNil, res)
		assert.Equal(t, errors.Errors[errors.InvalidAccount], err)
	}
}
