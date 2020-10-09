package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToRosettaPeers(t *testing.T) {
	result := exampleAddressBookEntries().ToRosettaPeers()

	assert.Equal(t, expectedRosettaPeers(), result)
	assert.Len(t, result, 1)
}
