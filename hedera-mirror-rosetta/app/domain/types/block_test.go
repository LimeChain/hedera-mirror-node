package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToRosettaBlock(t *testing.T) {
	exampleBlock := &Block{
		Index:               2,
		Hash:                "somehash",
		ConsensusStartNanos: 0,
		ConsensusEndNanos:   123,
		ParentIndex:         1,
		ParentHash:          "someparenthash",
		Transactions:        nil,
	}

	rosettaBlockResult := exampleBlock.ToRosettaBlock()

	assert.Equal(t, int64(0), rosettaBlockResult.Timestamp)
	assert.Equal(t, "0xsomehash", rosettaBlockResult.BlockIdentifier.Hash)
	assert.Equal(t, "0xsomeparenthash", rosettaBlockResult.ParentBlockIdentifier.Hash)
	assert.Len(t, rosettaBlockResult.Transactions, 0)
}

func TestGetTimestampMillis(t *testing.T) {
	exampleBlock := &Block{
		Index:               2,
		Hash:                "somehash",
		ConsensusStartNanos: 1000000,
		ConsensusEndNanos:   1123000,
		ParentIndex:         1,
		ParentHash:          "someparenthash",
		Transactions:        nil,
	}

	resultMillis := exampleBlock.GetTimestampMillis()

	assert.Equal(t, int64(1), resultMillis)
}
