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

package construction

import (
	"encoding/hex"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	nilConstructionCombineResponse *types.ConstructionCombineResponse = nil
)

func newDummyConstructionCombineRequestWith(unsignedTxHash, signingPayloadBytes, publicKeyBytes, signatureBytes string) *types.ConstructionCombineRequest {
	decodedSigningPayloadBytes, e1 := hex.DecodeString(signingPayloadBytes)
	decodedPublicKeyBytes, e2 := hex.DecodeString(publicKeyBytes)
	decodedSignatureBytes, e3 := hex.DecodeString(signatureBytes)

	if e1 != nil || e2 != nil || e3 != nil {
		return nil
	}

	return &types.ConstructionCombineRequest{
		NetworkIdentifier: &types.NetworkIdentifier{
			Blockchain: "someblockchain",
			Network:    "somenetwork",
			SubNetworkIdentifier: &types.SubNetworkIdentifier{
				Network:  "somesubnetwork",
				Metadata: nil,
			},
		},
		UnsignedTransaction: unsignedTxHash,
		Signatures: []*types.Signature{
			{
				SigningPayload: &types.SigningPayload{
					AccountIdentifier: &types.AccountIdentifier{
						Address:  "0.0.123352",
						Metadata: nil,
					},
					Bytes:         decodedSigningPayloadBytes,
					SignatureType: "ed25519",
				},
				PublicKey: &types.PublicKey{
					Bytes:     decodedPublicKeyBytes,
					CurveType: "edwards25519",
				},
				SignatureType: "ed25519",
				Bytes:         decodedSignatureBytes,
			},
		},
	}
}

func TestNewConstructionAPIService(t *testing.T) {
	repository.Setup()
	constructionService := NewConstructionAPIService()
	assert.IsType(t, &ConstructionAPIService{}, constructionService)
}

func TestConstructionCombine(t *testing.T) {
	// given:
	expectedConstructionCombineResponse := &types.ConstructionCombineResponse{
		SignedTransaction: "0x1a660a640a20d25025bad248dbd4c6ca704eefba7ab4f3e3f48089fa5f20e4e1d10303f97ade1a40967f26876ad492cc27b4c384dc962f443bcc9be33cbb7add3844bc864de047340e7a78c0fbaf40ab10948dc570bbc25edb505f112d0926dffb65c93199e6d507223c0a130a0b08c7af94fa0510f7d9fc76120418d8c307120218041880c2d72f2202087872180a160a090a0418d8c30710cf0f0a090a0418fec40710d00f",
	}

	repository.Setup()

	// when:
	res, e := NewConstructionAPIService().ConstructionCombine(nil, newDummyConstructionCombineRequest())

	// then:
	assert.Equal(t, expectedConstructionCombineResponse, res)
	assert.Nil(t, e)
}

func newDummyConstructionCombineRequest() *types.ConstructionCombineRequest {
	unsignedTxHash := "0x1a00223c0a130a0b08c7af94fa0510f7d9fc76120418d8c307120218041880c2d72f2202087872180a160a090a0418d8c30710cf0f0a090a0418fec40710d00f"
	signingPayloadBytes := "967f26876ad492cc27b4c384dc962f443bcc9be33cbb7add3844bc864de047340e7a78c0fbaf40ab10948dc570bbc25edb505f112d0926dffb65c93199e6d507"
	publicKeyBytes := "d25025bad248dbd4c6ca704eefba7ab4f3e3f48089fa5f20e4e1d10303f97ade"
	signatureBytes := "0a130a0b08c7af94fa0510f7d9fc76120418d8c307120218041880c2d72f2202087872180a160a090a0418d8c30710cf0f0a090a0418fec40710d00f"

	return newDummyConstructionCombineRequestWith(
		unsignedTxHash,
		signingPayloadBytes,
		publicKeyBytes,
		signatureBytes,
	)
}

func TestConstructionCombineThrowsWhenNoSingleSignature(t *testing.T) {
	// given:
	exampleConstructionCombineRequest := newDummyConstructionCombineRequest()
	exampleConstructionCombineRequest.Signatures = []*types.Signature{}

	// when:
	res, e := NewConstructionAPIService().ConstructionCombine(nil, exampleConstructionCombineRequest)

	// then:
	assert.Equal(t, nilConstructionCombineResponse, res)
	assert.Equal(t, errors.Errors[errors.MultipleSignaturesPresent], e)
}

func TestConstructionCombineThrowsWhenDecodeStringFails(t *testing.T) {
	// given:
	exampleInvalidTxHashConstructionCombineRequest := newDummyConstructionCombineRequest()
	exampleInvalidTxHashConstructionCombineRequest.UnsignedTransaction = "invalidTxHash"

	// when:
	res, e := NewConstructionAPIService().ConstructionCombine(nil, exampleInvalidTxHashConstructionCombineRequest)

	// then:
	assert.Equal(t, nilConstructionCombineResponse, res)
	assert.Equal(t, errors.Errors[errors.TransactionDecodeFailed], e)
}

func TestConstructionCombineThrowsWithInvalidPublicKey(t *testing.T) {
	// given:
	exampleInvalidPublicKeyConstructionCombineRequest := newDummyConstructionCombineRequest()
	exampleInvalidPublicKeyConstructionCombineRequest.Signatures[0].PublicKey = &types.PublicKey{}

	// when:
	res, e := NewConstructionAPIService().ConstructionCombine(nil, exampleInvalidPublicKeyConstructionCombineRequest)

	// then:
	assert.Equal(t, nilConstructionCombineResponse, res)
	assert.Equal(t, errors.Errors[errors.InvalidPublicKey], e)
}

func TestConstructionCombineThrowsWhenSignatureIsNotVerified(t *testing.T) {
	// given:
	exampleInvalidSigningPayloadConstructionCombineRequest := newDummyConstructionCombineRequest()
	exampleInvalidSigningPayloadConstructionCombineRequest.Signatures[0].SigningPayload = &types.SigningPayload{}

	// when:
	res, e := NewConstructionAPIService().ConstructionCombine(nil, exampleInvalidSigningPayloadConstructionCombineRequest)

	// then:
	assert.Equal(t, nilConstructionCombineResponse, res)
	assert.Equal(t, errors.Errors[errors.InvalidSignatureVerification], e)
}
