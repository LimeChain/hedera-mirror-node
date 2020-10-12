package types

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exampleAddressBookEntry() *AddressBookEntry {
	return &AddressBookEntry{
		PeerId:   exampleAccount(),
		Metadata: nil,
	}
}

func exampleAddressBookEntries() *AddressBookEntries {
	return &AddressBookEntries{
		[]*AddressBookEntry{
			exampleAddressBookEntry(),
		},
	}
}

func expectedRosettaPeers() []*types.Peer {
	return []*types.Peer{
		{PeerID: exampleAccount().String(), Metadata: nil},
	}
}

func TestToRosettaPeers(t *testing.T) {
	// when:
	result := exampleAddressBookEntries().ToRosetta()

	// then:
	assert.Equal(t, expectedRosettaPeers(), result)
	assert.Len(t, result, 1)
}
