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
				Metadata: nil,
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
			Metadata: nil,
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
