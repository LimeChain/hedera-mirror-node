package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/services"
)

// Account is domain level struct used to represent Rosetta Account
type Account struct {
	Shard  int64
	Realm  int64
	Number int64
}

// NewAccountFromEncodedID - creates new instance of Account struct
func NewAccountFromEncodedID(encodedID int64) (*Account, error) {
	d, err := services.Decode(encodedID)
	if err != nil {
		return nil, err
	}

	return &Account{
		Shard:  d.ShardNum,
		Realm:  d.RealmNum,
		Number: d.EntityNum,
	}, err
}

// FormatToString - returns the string representation of the account
func (a *Account) FormatToString() string {
	return fmt.Sprintf("%d.%d.%d", a.Shard, a.Realm, a.Number)
}

// FromRosettaAccount populates domain type Account from Rosetta type Account
func FromRosettaAccount(rAccount *rTypes.AccountIdentifier) (*Account, error) {
	inputs := strings.Split(rAccount.Address, ".")
	if len(inputs) != 3 {
		return nil, errors.New("Invalid account provided")
	}

	shard, err := strconv.Atoi(inputs[0])
	realm, err1 := strconv.Atoi(inputs[1])
	number, err2 := strconv.Atoi(inputs[2])
	if err != nil || err1 != nil || err2 != nil {
		return nil, errors.New("Invalid account provided")
	}

	return &Account{
		Shard:  int64(shard),
		Realm:  int64(realm),
		Number: int64(number),
	}, nil
}

// ToRosettaAccount returns Rosetta type Account from the current domain type Account
func (a *Account) ToRosettaAccount() *rTypes.AccountIdentifier {
	return &rTypes.AccountIdentifier{
		Address: a.FormatToString(),
	}
}
