package services

import (
	"context"
	"encoding/hex"
	"github.com/coinbase/rosetta-sdk-go/server"
	rTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/errors"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
	hederatools "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/hedera"
	hexutils "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/hex"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/tools/validator"
	"github.com/hashgraph/hedera-sdk-go"
	"strconv"
)

type ConstructionService struct {
	hederaClient *hedera.Client
}

func (c *ConstructionService) ConstructionCombine(ctx context.Context, request *rTypes.ConstructionCombineRequest) (*rTypes.ConstructionCombineResponse, *rTypes.Error) {
	panic("implement me")
}

func (c *ConstructionService) ConstructionDerive(ctx context.Context, request *rTypes.ConstructionDeriveRequest) (*rTypes.ConstructionDeriveResponse, *rTypes.Error) {
	return nil, errors.Errors[errors.NotImplemented]
}

func (c *ConstructionService) ConstructionHash(ctx context.Context, request *rTypes.ConstructionHashRequest) (*rTypes.TransactionIdentifierResponse, *rTypes.Error) {
	panic("implement me")
}

func (c *ConstructionService) ConstructionMetadata(ctx context.Context, request *rTypes.ConstructionMetadataRequest) (*rTypes.ConstructionMetadataResponse, *rTypes.Error) {
	return &rTypes.ConstructionMetadataResponse{}, nil
}

func (c *ConstructionService) ConstructionParse(ctx context.Context, request *rTypes.ConstructionParseRequest) (*rTypes.ConstructionParseResponse, *rTypes.Error) {
	panic("implement me")
}

func (c *ConstructionService) ConstructionPayloads(ctx context.Context, request *rTypes.ConstructionPayloadsRequest) (*rTypes.ConstructionPayloadsResponse, *rTypes.Error) {
	operationType, err := validator.ValidateOperations(request.Operations)
	if err != nil {
		return nil, err
	}

	switch *operationType {
	case config.OperationTypeCryptoTransfer:
		return c.handleCryptoTransferPayload(request.Operations)
	default:
		return c.handleCryptoCreateAccountPayload(request.Operations)
	}
}
func (c *ConstructionService) ConstructionPreprocess(ctx context.Context, request *rTypes.ConstructionPreprocessRequest) (*rTypes.ConstructionPreprocessResponse, *rTypes.Error) {
	return &rTypes.ConstructionPreprocessResponse{}, nil
}

func (c *ConstructionService) ConstructionSubmit(ctx context.Context, request *rTypes.ConstructionSubmitRequest) (*rTypes.TransactionIdentifierResponse, *rTypes.Error) {
	panic("implement me")
}

func (c *ConstructionService) handleCryptoCreateAccountPayload(operations []*rTypes.Operation) (*rTypes.ConstructionPayloadsResponse, *rTypes.Error) {
	operationsLength := len(operations)
	if operationsLength != 1 {
		return nil, errors.Errors[errors.InvalidOperationsAmount]
	}

	operation := operations[0]
	sender, err := hedera.AccountIDFromString(operation.Account.Address)
	if err != nil {
		return nil, errors.Errors[errors.InvalidAccount]
	}

	amount, err := strconv.Atoi(operation.Amount.Value)
	if err != nil {
		return nil, errors.Errors[errors.InvalidAmount]
	}

	transaction, err := hedera.
		NewAccountCreateTransaction().
		SetInitialBalance(hederatools.ToHbarAmount(int64(amount))).
		SetTransactionID(hederatools.TransactionId(sender)).
		Build(c.hederaClient)

	if err != nil {
		return nil, errors.Errors[errors.TransactionBuildFailed]
	}

	bytesTransaction, err := transaction.MarshalBinary()
	if err != nil {
		return nil, errors.Errors[errors.TransactionMarshallingFailed]
	}

	return &rTypes.ConstructionPayloadsResponse{
		UnsignedTransaction: hexutils.SafeAddHexPrefix(hex.EncodeToString(bytesTransaction)),
		Payloads: []*rTypes.SigningPayload{{
			Address: sender.String(),
			Bytes:   bytesTransaction,
		}},
	}, nil
}

func (c *ConstructionService) handleCryptoTransferPayload(operations []*rTypes.Operation) (*rTypes.ConstructionPayloadsResponse, *rTypes.Error) {
	builderTransaction := hedera.NewCryptoTransferTransaction()
	var sender hedera.AccountID

	for _, operation := range operations {
		account, err := hedera.AccountIDFromString(operation.Account.Address)
		if err != nil {
			return nil, errors.Errors[errors.InvalidAccount]
		}

		amount, err := strconv.Atoi(operation.Amount.Value)
		if err != nil {
			return nil, errors.Errors[errors.InvalidAmount]
		}

		if amount < 0 {
			sender = account
			builderTransaction.AddSender(
				sender,
				hedera.HbarFromTinybar(int64(amount)))
		} else {
			builderTransaction.AddRecipient(sender,
				hederatools.ToHbarAmount(int64(amount)))
		}
	}

	transaction, err := builderTransaction.SetTransactionID(hedera.NewTransactionID(sender)).Build(c.hederaClient)
	if err != nil {
		return nil, errors.Errors[errors.TransactionBuildFailed]
	}

	bytesTransaction, err := transaction.MarshalBinary()
	if err != nil {
		return nil, errors.Errors[errors.TransactionMarshallingFailed]
	}

	return &rTypes.ConstructionPayloadsResponse{
		UnsignedTransaction: hexutils.SafeAddHexPrefix(hex.EncodeToString(bytesTransaction)),
		Payloads: []*rTypes.SigningPayload{{
			Address: sender.String(),
			Bytes:   bytesTransaction,
		}},
	}, nil
}

func NewConstructionAPIService() server.ConstructionAPIServicer {
	return &ConstructionService{
		hederaClient: hedera.ClientForTestnet(),
	}
}
