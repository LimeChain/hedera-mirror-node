package maphelper

import (
	"fmt"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetsCorrectStringValuesFromMap(t *testing.T) {
	// given:
	testData := map[int]string{
		1: "abc",
	}

	// when:
	result := GetStringValuesFromIntStringMap(testData)

	// then:
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "abc", result[0])
}

func TestGetsCorrectErrorValuesFromMap(t *testing.T) {
	// given:
	error := newErrorDummy(32, true)

	testData := map[string]*types.Error{
		"error": error,
	}

	// when:
	result := GetErrorValuesFromStringErrorMap(testData)

	// then:
	assert.Equal(t, 1, len(result))
	assert.Equal(t, error, result[0])
}

func newErrorDummy(code int32, retryable bool) *types.Error {
	return errors.New(fmt.Sprintf("error_dummy_%d", code), code, retryable)
}
