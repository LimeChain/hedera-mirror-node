package hedera

import (
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-sdk-go"
)

func ToHederaAccountId(account *types.Account) hedera.AccountID {
	return hedera.AccountID{
		Shard:   uint64(account.Shard),
		Realm:   uint64(account.Realm),
		Account: uint64(account.Number),
	}
}

func ToHbarAmount(amount int64) hedera.Hbar {
	return hedera.HbarFromTinybar(amount)
}

func TransactionId(account hedera.AccountID) hedera.TransactionID {
	return hedera.NewTransactionID(account)
}
