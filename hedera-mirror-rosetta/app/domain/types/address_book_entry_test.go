package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToRosettaPeers(t *testing.T) {
	exampleAccount, _ := AccountFromString("0.0.0")
	var testData = []*AddressBookEntry{
		{exampleAccount, nil},
	}
	exampleEntries := &AddressBookEntries{
		testData,
	}
	result := exampleEntries.ToRosettaPeers()

	assert.Len(t, result, 1)
	assert.Equal(t, testData[0].PeerId.String(), result[0].PeerID)
	assert.Equal(t, testData[0].Metadata, result[0].Metadata)
}
