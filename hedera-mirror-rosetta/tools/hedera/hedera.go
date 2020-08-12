package hedera

import (
	"github.com/hashgraph/hedera-sdk-go"
)

func ToHbarAmount(amount int64) hedera.Hbar {
	return hedera.HbarFromTinybar(amount)
}

func TransactionId(account hedera.AccountID) hedera.TransactionID {
	return hedera.NewTransactionID(account)
}
