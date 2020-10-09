package types

import (
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
)

func expectedOperationIdentifier() *types.OperationIdentifier {
	return &types.OperationIdentifier{
		Index:        1,
		NetworkIndex: nil,
	}
}

func expectedRelatedOperations() []*types.OperationIdentifier {
	return []*types.OperationIdentifier{}
}

func expectedAccount() *types.AccountIdentifier {
	return &types.AccountIdentifier{
		Address:    "0.0.0",
		SubAccount: nil,
		Metadata:   nil,
	}
}

func expectedAmount() *types.Amount {
	return &types.Amount{
		Value:    "400",
		Currency: config.CurrencyHbar,
		Metadata: nil,
	}
}

func expectedTransaction() *types.Transaction {
	return &types.Transaction{
		TransactionIdentifier: &types.TransactionIdentifier{Hash: "somehash"},
		Operations:            []*types.Operation{expectedOperation()},
		Metadata:              nil,
	}
}

func expectedBlockIdentifier() *types.BlockIdentifier {
	return &types.BlockIdentifier{
		Index: 2,
		Hash:  "0xsomehash",
	}
}

func expectedParentBlockIdentifier() *types.BlockIdentifier {
	return &types.BlockIdentifier{
		Index: 1,
		Hash:  "0xsomeparenthash",
	}
}

func expectedBlock() *types.Block {
	return &types.Block{
		BlockIdentifier:       expectedBlockIdentifier(),
		ParentBlockIdentifier: expectedParentBlockIdentifier(),
		Timestamp:             int64(10),
		Transactions:          []*types.Transaction{expectedTransaction()},
	}
}

func expectedOperation() *types.Operation {
	return &types.Operation{
		OperationIdentifier: expectedOperationIdentifier(),
		RelatedOperations:   expectedRelatedOperations(),
		Type:                "transfer",
		Status:              "pending",
		Account:             expectedAccount(),
		Amount:              expectedAmount(),
	}
}

func expectedRosettaPeers() []*types.Peer {
	return []*types.Peer{
		{PeerID: exampleAccount().String(), Metadata: nil},
	}
}

func expectedAccountWith(shard int64, realm int64, number int64) *Account {
	return &Account{
		Shard:  shard,
		Realm:  realm,
		Number: number,
	}
}
