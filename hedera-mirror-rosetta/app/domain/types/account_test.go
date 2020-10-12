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
	entityid "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/services/encoding"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exampleAccount() *Account {
	return &Account{
		entityid.EntityId{
			ShardNum:  0,
			RealmNum:  0,
			EntityNum: 0,
		},
	}
}

func exampleAccountWith(shard, realm, entity int64) *Account {
	return &Account{
		entityid.EntityId{
			ShardNum:  shard,
			RealmNum:  realm,
			EntityNum: entity,
		},
	}
}

func expectedAccount() *types.AccountIdentifier {
	return &types.AccountIdentifier{
		Address:    "0.0.0",
		SubAccount: nil,
		Metadata:   nil,
	}
}

func expectedAccountWith(shard int64, realm int64, number int64) *Account {
	return &Account{
		entityid.EntityId{
			ShardNum:  shard,
			RealmNum:  realm,
			EntityNum: number,
		},
	}
}

func TestToRosettaAccount(t *testing.T) {
	// when:
	rosettaAccount := exampleAccount().ToRosetta()

	// then:
	assert.Equal(t, expectedAccount(), rosettaAccount)
}

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
		assert.Equal(t, expectedAccountWith(tt.shard, tt.realm, tt.number), res)
	}
}

func TestComputeEncodedID(t *testing.T) {
	var testData = []struct {
		shard, realm, number int64
	}{
		{-1, 123, 246},
		{123, -123, 246},
		{123, 23, -246},
	}

	for _, tt := range testData {
		res, e := exampleAccountWith(tt.shard, tt.realm, tt.number).ComputeEncodedID()
		assert.Zero(t, res)
		assert.NotNil(t, e)
	}
}

func TestNewAccountFromEncodedIDThrows(t *testing.T) {
	// given:
	var testData = struct {
		input int64
	}{
		-123,
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
		assert.Equal(t, expectedAccountWith(tt.shard, tt.realm, tt.number), res)
	}
}

func TestAccountFromStringThrows(t *testing.T) {
	// given:
	var testData = []struct {
		input string
	}{
		{"a.0.0"},
		{"0.b.0"},
		{"0.0c"},
		{"0.0.c"},
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
