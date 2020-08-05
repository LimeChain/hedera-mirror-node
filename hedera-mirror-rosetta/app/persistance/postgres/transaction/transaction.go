package transaction

import (
	"log"

	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/jinzhu/gorm"
)

type transaction struct {
	ConsensusNS          int64  `gorm:"type:bigint;primary_key"`
	Type                 int    `gorm:"type:smallint"`
	Result               int    `gorm:"type:smallint"`
	PayerAccountID       int64  `gorm:"type:bigint"`
	ValidStartNS         int64  `gorm:"type:bigint"`
	ValidDurationSeconds int64  `gorm:"type:bigint"`
	NodeAccountID        int64  `gorm:"type:bigint"`
	EntityID             int64  `gorm:"type:bigint"`
	InitialBalance       int64  `gorm:"type:bigint"`
	MaxFee               int64  `gorm:"type:bigint"`
	ChargedTxFee         int64  `gorm:"type:bigint"`
	Memo                 []byte `gorm:"type:bytea"`
	TransactionHash      []byte `gorm:"type:bytea"`
	TransactionBytes     []byte `gorm:"type:bytea"`
}

type cryptoTransfer struct {
}

// TableName - Set table name to be `record_file`
func (transaction) TableName() string {
	return "transaction"
}

// TableName - Set table name to be `record_file`
func (cryptoTransfer) TableName() string {
	return "crypto_transfer"
}

// TransactionRepository struct that has connection to the Database
type TransactionRepository struct {
	dbClient *gorm.DB
}

// NewTransactionRepository creates an instance of a TransactionRepository struct
func NewTransactionRepository(dbClient *gorm.DB) *TransactionRepository {
	return &TransactionRepository{dbClient: dbClient}
}

// FindByTimestamp retrieves Transaction by given timestmap
func (tr *TransactionRepository) FindByTimestamp(timestamp int64) *types.Transaction {
	t := &transaction{}
	tr.dbClient.Find(t, timestamp)

	log.Println(t)
	return nil
}

// FindBetween retrieves all Transactions between the provided start and end timestamp
func (tr *TransactionRepository) FindBetween(start int64, end int64) []*types.Transaction {
	if start > end {
		// TODO throw error
	}
	tArray := &[]transaction{}
	tr.dbClient.Where("consensus_ns >= ? AND consensus_ns <= ?", start, end).Find(tArray)

	log.Println(tArray)
	return nil
}
