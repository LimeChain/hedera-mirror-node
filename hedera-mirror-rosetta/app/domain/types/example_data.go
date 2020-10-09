package types

func exampleBlock() *Block {
	return &Block{
		Index:               2,
		Hash:                "somehash",
		ConsensusStartNanos: 10000000,
		ConsensusEndNanos:   12300000,
		ParentIndex:         1,
		ParentHash:          "someparenthash",
		Transactions:        exampleTransactions(),
	}
}

func exampleTransactions() []*Transaction {
	return []*Transaction{
		exampleTransaction(),
	}
}

func exampleTransaction() *Transaction {
	return &Transaction{
		Hash:       "somehash",
		Operations: exampleOperations(),
	}
}

func exampleOperations() []*Operation {
	return []*Operation{
		exampleOperation(),
	}
}

func exampleOperation() *Operation {
	return &Operation{
		Index:   1,
		Type:    "transfer",
		Status:  "pending",
		Account: exampleAccount(),
		Amount:  exampleAmount(),
	}
}

func exampleAccount() *Account {
	return &Account{
		Shard:  0,
		Realm:  0,
		Number: 0,
	}
}

func exampleAmount() *Amount {
	return &Amount{Value: int64(400)}
}

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
