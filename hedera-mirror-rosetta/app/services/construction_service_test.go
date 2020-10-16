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

package services

import (
	"encoding/hex"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newDummyConstructionCombineRequest(unsignedTxHash, signingPayloadBytes, publicKeyBytes, signatureBytes string) *types.ConstructionCombineRequest {
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
	mocks.Setup()
	constructionService := NewConstructionAPIService()
	assert.IsType(t, &ConstructionAPIService{}, constructionService)
}

func TestConstructionCombine(t *testing.T) {
	// given:
	exampleConstructionCombineRequest := newDummyConstructionCombineRequest(
		"0x1a00223c0a130a0b08c7af94fa0510f7d9fc76120418d8c307120218041880c2d72f2202087872180a160a090a0418d8c30710cf0f0a090a0418fec40710d00f",
		"967f26876ad492cc27b4c384dc962f443bcc9be33cbb7add3844bc864de047340e7a78c0fbaf40ab10948dc570bbc25edb505f112d0926dffb65c93199e6d507",
		"d25025bad248dbd4c6ca704eefba7ab4f3e3f48089fa5f20e4e1d10303f97ade",
		"0a130a0b08c7af94fa0510f7d9fc76120418d8c307120218041880c2d72f2202087872180a160a090a0418d8c30710cf0f0a090a0418fec40710d00f",
	)
	expectedConstructionCombineResponse := &types.ConstructionCombineResponse{
		SignedTransaction: "0x1a660a640a20d25025bad248dbd4c6ca704eefba7ab4f3e3f48089fa5f20e4e1d10303f97ade1a40967f26876ad492cc27b4c384dc962f443bcc9be33cbb7add3844bc864de047340e7a78c0fbaf40ab10948dc570bbc25edb505f112d0926dffb65c93199e6d507223c0a130a0b08c7af94fa0510f7d9fc76120418d8c307120218041880c2d72f2202087872180a160a090a0418d8c30710cf0f0a090a0418fec40710d00f",
	}

	mocks.Setup()

	// when:
	res, e := NewConstructionAPIService().ConstructionCombine(nil, exampleConstructionCombineRequest)

	// then:
	assert.Equal(t, expectedConstructionCombineResponse, res)
	assert.Nil(t, e)
}
