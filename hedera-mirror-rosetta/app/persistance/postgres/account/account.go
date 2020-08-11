package account

import (
	"fmt"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/jinzhu/gorm"
)

const (
	whereClauseBeforeConsensusEnd string = "consensus_timestamp = (SELECT MAX(consensus_timestamp) FROM %s WHERE consensus_timestamp <= %d)"
)

type accountBalance struct {
	ConsensusTimestamp int64 `gorm:"type:bigint;primary_key"`
	Balance            int64 `gorm:"type:bigint"`
	AccountRealmNum    int16 `gorm:"type:smallint;primary_key"`
	AccountNum         int32 `gorm:"type:integer;primary_key"`
}

// TableName - Set table name of the accountBalance to be `account_balance`
func (accountBalance) TableName() string {
	return "account_balance"
}

// AccountRepository struct that has connection to the Database
type AccountRepository struct {
	dbClient *gorm.DB
}

// NewAccountRepository creates an instance of a TransactionRepository struct. Populates the transaction types and statuses on init
func NewAccountRepository(dbClient *gorm.DB) *AccountRepository {
	return &AccountRepository{
		dbClient: dbClient,
	}
}

// RetrieveBalanceAtBlock returns the balance of the account at a given block (provided by consensusEnd timestamp)
func (ar *AccountRepository) RetrieveBalanceAtBlock(addressStr string, consensusEnd int64) (*types.Amount, *rTypes.Error) {
	acc, err := types.AccountFromString(addressStr)
	if err != nil {
		return nil, err
	}

	ab := &accountBalance{}
	if ar.dbClient.Where(&accountBalance{AccountRealmNum: int16(acc.Realm), AccountNum: int32(acc.Number)}).Where(fmt.Sprintf(whereClauseBeforeConsensusEnd, ab.TableName(), consensusEnd)).Find(&ab).RecordNotFound() {
		return nil, errors.Errors[errors.AccountNotFound]
	}

	return &types.Amount{
		Value: ab.Balance,
	}, nil
}
