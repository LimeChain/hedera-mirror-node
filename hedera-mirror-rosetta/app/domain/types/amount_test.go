package types

import (
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/parse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToRosettaAmount(t *testing.T) {
	properAmount, _ := parse.ToInt64("100")
	exampleAmount := Amount{
		properAmount,
	}

	result := exampleAmount.ToRosettaAmount()

	assert.Equal(t, "100", result.Value)
	assert.Equal(t, config.CurrencyHbar, result.Currency)
}
