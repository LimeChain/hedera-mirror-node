package transaction

import (
	"encoding/hex"
	"fmt"
	dbTypes "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistence/postgres/common"
	hexUtils "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/hex"
	"log"
	"strconv"
	"strings"

	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/maphelper"

	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/domain/types"
	"github.com/jinzhu/gorm"
)

const (
	whereClauseBetweenConsensus         string = "consensus_ns >= ? AND consensus_ns <= ?"
	whereTimestampsInConsensusTimestamp string = "consensus_timestamp IN (%s)"
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

type transactionType struct {
	ProtoID int    `gorm:"type:integer;primary_key"`
	Name    string `gorm:"size:30"`
}

type transactionResult struct {
	ProtoID int    `gorm:"type:integer;primary_key"`
	Result  string `gorm:"size:100"`
}

// TableName - Set table name of the Transactions to be `record_file`
func (transaction) TableName() string {
	return "transaction"
}

// TableName - Set table name of the Transaction Types to be `t_transaction_types`
func (transactionType) TableName() string {
	return "t_transaction_types"
}

// TableName - Set table name of the Transaction Results to be `t_transaction_results`
func (transactionResult) TableName() string {
	return "t_transaction_results"
}

func (t *transaction) getHashString() string {
	return hexUtils.SafeAddHexPrefix(hex.EncodeToString(t.TransactionHash))
}

// TransactionRepository struct that has connection to the Database
type TransactionRepository struct {
	dbClient *gorm.DB
	statuses map[int]string
	types    map[int]string
}

// NewTransactionRepository creates an instance of a TransactionRepository struct
func NewTransactionRepository(dbClient *gorm.DB) *TransactionRepository {
	return &TransactionRepository{dbClient: dbClient}
}

// Types returns map of all Transaction Types
func (tr *TransactionRepository) Types() (map[int]string, *rTypes.Error) {
	if tr.types != nil {
		return tr.types, nil
	}

	typesArray := tr.retrieveTransactionTypes()
	if len(typesArray) == 0 {
		log.Println("No Transaction Types were found in the database.")
		return nil, errors.Errors[errors.OperationTypesNotFound]
	}

	tr.types = make(map[int]string)
	for _, t := range typesArray {
		tr.types[t.ProtoID] = t.Name
	}
	return tr.types, nil
}

// Statuses returns map of all Transaction Results
func (tr *TransactionRepository) Statuses() (map[int]string, *rTypes.Error) {
	if tr.statuses != nil {
		return tr.statuses, nil
	}

	rArray := tr.retrieveTransactionResults()
	if len(rArray) == 0 {
		log.Println("No Transaction Results were found in the database.")
		return nil, errors.Errors[errors.OperationStatusesNotFound]
	}

	tr.statuses = make(map[int]string)
	for _, s := range rArray {
		tr.statuses[s.ProtoID] = s.Result
	}
	return tr.statuses, nil
}

func (tr *TransactionRepository) TypesAsArray() ([]string, *rTypes.Error) {
	transactionTypes, err := tr.Types()
	if err != nil {
		return nil, err
	}
	return maphelper.GetStringValuesFromIntStringMap(transactionTypes), nil
}

// FindBetween retrieves all Transactions between the provided start and end timestamp
func (tr *TransactionRepository) FindBetween(start int64, end int64) ([]*types.Transaction, *rTypes.Error) {
	if start > end {
		return nil, errors.Errors[errors.StartMustNotBeAfterEnd]
	}
	var transactions []transaction
	tr.dbClient.Where(whereClauseBetweenConsensus, start, end).Find(&transactions)

	sameHashMap := make(map[string][]transaction)
	for _, t := range transactions {
		h := t.getHashString()
		sameHashMap[h] = append(sameHashMap[h], t)
	}
	res := make([]*types.Transaction, 0, len(sameHashMap))
	for _, sameHashTransactions := range sameHashMap {
		transaction, err := tr.constructTransaction(sameHashTransactions)
		if err != nil {
			return nil, err
		}
		res = append(res, transaction)
	}
	return res, nil
}

// FindByHashInBlock retrieves a transaction by Hash
func (tr *TransactionRepository) FindByHashInBlock(hashStr string, consensusStart int64, consensusEnd int64) (*types.Transaction, *rTypes.Error) {
	var transactions []transaction
	transactionHash, err := hex.DecodeString(hexUtils.SafeRemoveHexPrefix(hashStr))
	if err != nil {
		return nil, errors.Errors[errors.InvalidTransactionIdentifier]
	}
	tr.dbClient.Where(&transaction{TransactionHash: transactionHash}).Find(&transactions)
	transactions = filterTransactionsForRange(transactions, consensusStart, consensusEnd)

	if len(transactions) == 0 {
		return nil, errors.Errors[errors.TransactionNotFound]
	}

	transaction, err1 := tr.constructTransaction(transactions)
	if err1 != nil {
		return nil, err1
	}
	return transaction, nil
}

// filterTransactionsForRange - Filters the passed transactions. If the ConnsensusNS is not in the given [consensusStart; consensusEnd] range, the transaction is removed from the list
func filterTransactionsForRange(transactions []transaction, consensusStart int64, consensusEnd int64) []transaction {
	var length int
	for _, t := range transactions {
		if t.ConsensusNS < consensusStart || t.ConsensusNS > consensusEnd {
			continue
		}
		transactions[length] = t
		length++
	}
	return transactions[:length]
}

func (tr *TransactionRepository) findCryptoTransfers(timestamps []int64) []dbTypes.CryptoTransfer {
	var cryptoTransfers []dbTypes.CryptoTransfer
	timestampsStr := intsToString(timestamps)
	tr.dbClient.Where(fmt.Sprintf(whereTimestampsInConsensusTimestamp, timestampsStr)).Find(&cryptoTransfers)
	return cryptoTransfers
}

func (tr *TransactionRepository) retrieveTransactionTypes() []transactionType {
	var transactionTypes []transactionType
	tr.dbClient.Find(&transactionTypes)
	return transactionTypes
}

func (tr *TransactionRepository) retrieveTransactionResults() []transactionResult {
	var tResults []transactionResult
	tr.dbClient.Find(&tResults)
	return tResults
}

func (tr *TransactionRepository) constructTransaction(sameHashTransactions []transaction) (*types.Transaction, *rTypes.Error) {
	tResult := &types.Transaction{Hash: sameHashTransactions[0].getHashString()}

	transactionsMap := make(map[int64]transaction)
	timestamps := make([]int64, len(sameHashTransactions))
	for i, t := range sameHashTransactions {
		transactionsMap[t.ConsensusNS] = t
		timestamps[i] = t.ConsensusNS
	}
	cryptoTransfers := tr.findCryptoTransfers(timestamps)
	operations, err := tr.constructOperations(cryptoTransfers, transactionsMap)
	if err != nil {
		return nil, err
	}
	tResult.Operations = operations

	return tResult, nil
}

func (tr *TransactionRepository) constructOperations(cryptoTransfers []dbTypes.CryptoTransfer, transactionsMap map[int64]transaction) ([]*types.Operation, *rTypes.Error) {
	transactionTypes, err := tr.Types()
	if err != nil {
		return nil, err
	}
	transactionStatuses, err := tr.Statuses()
	if err != nil {
		return nil, err
	}

	operations := make([]*types.Operation, len(cryptoTransfers))
	for i, ct := range cryptoTransfers {
		a, err := constructAccount(ct.EntityID)
		if err != nil {
			return nil, err
		}
		operationType := transactionTypes[transactionsMap[ct.ConsensusTimestamp].Type]
		operationStatus := transactionStatuses[transactionsMap[ct.ConsensusTimestamp].Result]
		operations[i] = &types.Operation{Index: int64(i), Type: operationType, Status: operationStatus, Account: a, Amount: &types.Amount{Value: ct.Amount}}
	}
	return operations, nil
}

func constructAccount(encodedID int64) (*types.Account, *rTypes.Error) {
	acc, err := types.NewAccountFromEncodedID(encodedID)
	if err != nil {
		log.Printf(errors.CreateAccountDbIdFailed, encodedID)
		return nil, errors.Errors[errors.InternalServerError]
	}
	return acc, nil
}

func intsToString(ints []int64) string {
	r := make([]string, len(ints))
	for i, v := range ints {
		r[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(r, ",")
}
