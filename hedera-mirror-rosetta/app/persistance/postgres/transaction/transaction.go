package transaction

import "github.com/jinzhu/gorm"

type transaction struct {
	ID int64 `gorm:"type:bigint;primary_key"`
}

type TransactionRepository struct {
	dbClient *gorm.DB
}

func NewTransactionRepository(dbClient *gorm.DB) *TransactionRepository {
	return &TransactionRepository{dbClient: dbClient}
}
