package errors

import (
	"github.com/coinbase/rosetta-sdk-go/types"
)

// Errors - map of all Errors that this API can return
var Errors = map[string]*types.Error{
	BlockNotFound:                  New(BlockNotFound, 1, true),
	AccountNotFound:                New(AccountNotFound, 2, true),
	TransactionBuildFailed:         New(TransactionBuildFailed, 3, false),
	TransactionDecodeFailed:        New(TransactionDecodeFailed, 4, false),
	TransactionMarshallingFailed:   New(TransactionMarshallingFailed, 5, false),
	TransactionUnmarshallingFailed: New(TransactionUnmarshallingFailed, 6, false),
	TransactionNotFound:            New(TransactionNotFound, 7, true),
	MultipleOperationTypesPresent:  New(MultipleOperationTypesPresent, 8, false),
	StartMustBeBeforeEnd:           New(StartMustBeBeforeEnd, 9, false),
	InvalidAccount:                 New(InvalidAccount, 10, false),
	InvalidAmount:                  New(InvalidAmount, 11, false),
	InvalidOperationsAmount:        New(InvalidOperationsAmount, 12, false),
	InvalidOperationsTotalAmount:   New(InvalidOperationsTotalAmount, 13, false),
	InvalidTransactionIdentifier:   New(InvalidTransactionIdentifier, 14, false),
	NotImplemented:                 New(NotImplemented, 15, false),
}

const (
	BlockNotFound                  string = "Block not found"
	AccountNotFound                string = "Account not found"
	TransactionBuildFailed         string = "Transaction build failed"
	TransactionDecodeFailed        string = "Transaction Decode failed"
	TransactionMarshallingFailed   string = "Transaction marshalling failed"
	TransactionUnmarshallingFailed string = "Transaction unmarshalling failed"
	TransactionNotFound            string = "Transaction not found"
	MultipleOperationTypesPresent  string = "Only one Operation Type must be present"
	StartMustBeBeforeEnd           string = "Start must be before end"
	InvalidAccount                 string = "Invalid Account provided"
	InvalidAmount                  string = "Invalid Amount provided"
	InvalidOperationsAmount        string = "Invalid Operations amount provided"
	InvalidOperationsTotalAmount   string = "Operations total amount must be 0"
	InvalidTransactionIdentifier   string = "Invalid Transaction Identifier provided"
	NotImplemented                 string = "Not implemented"
)

func New(message string, statusCode int32, retryable bool) *types.Error {
	return &types.Error{
		Message:   message,
		Code:      statusCode,
		Retriable: retryable,
		Details:   nil,
	}
}
