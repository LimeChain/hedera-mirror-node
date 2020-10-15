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
	"github.com/stretchr/testify/assert"
	"testing"
)

func exampleAddressBookEntries() *AddressBookEntries {
	return &AddressBookEntries{
		[]*AddressBookEntry{
			{
				PeerId: &Account{
					entityid.EntityId{
						ShardNum:  0,
						RealmNum:  0,
						EntityNum: 0,
					},
				},
				Metadata: map[string]interface{}{
					"ip":   "123",
					"port": "20514",
				},
			},
		},
	}
}

func expectedRosettaPeers() []*types.Peer {
	return []*types.Peer{
		{
			PeerID: (&Account{
				entityid.EntityId{
					ShardNum:  0,
					RealmNum:  0,
					EntityNum: 0,
				},
			}).String(),
			Metadata: map[string]interface{}{
				"ip":   "123",
				"port": "20514",
			},
		},
	}
}

func TestToRosettaPeers(t *testing.T) {
	// when:
	result := exampleAddressBookEntries().ToRosetta()

	// then:
	assert.Equal(t, expectedRosettaPeers(), result)
	assert.Len(t, result, 1)
}
