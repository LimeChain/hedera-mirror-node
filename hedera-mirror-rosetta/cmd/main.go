package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistance/postgres/account"
	addressBookEntry "github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistance/postgres/addressbook/entry"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistance/postgres/block"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistance/postgres/transaction"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/services"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/config"
	"github.com/jinzhu/gorm"
)

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(network *types.NetworkIdentifier, asserter *asserter.Asserter, version *types.Version, dbClient *gorm.DB) http.Handler {
	blockRepo := block.NewBlockRepository(dbClient)
	transactionRepo := transaction.NewTransactionRepository(dbClient)
	accountRepo := account.NewAccountRepository(dbClient)
	addressBookEntryRepo := addressBookEntry.NewAddressBookEntryRepository(dbClient)

	commons := services.NewCommons(blockRepo, transactionRepo)

	networkAPIService := services.NewNetworkAPIService(commons, addressBookEntryRepo, network, version)
	networkAPIController := server.NewNetworkAPIController(networkAPIService, asserter)

	blockAPIService := services.NewBlockAPIService(commons)
	blockAPIController := server.NewBlockAPIController(blockAPIService, asserter)

	mempoolAPIService := services.NewMempoolAPIService()
	mempoolAPIController := server.NewMempoolAPIController(mempoolAPIService, asserter)

	constructionAPIService := services.NewConstructionAPIService()
	constructionAPIController := server.NewConstructionAPIController(constructionAPIService, asserter)

	accountAPIService := services.NewAccountAPIService(commons, accountRepo)
	accountAPIController := server.NewAccountAPIController(accountAPIService, asserter)

	return server.NewRouter(networkAPIController, blockAPIController, mempoolAPIController, constructionAPIController, accountAPIController)
}

func main() {
	configuration := LoadConfig()

	network := &types.NetworkIdentifier{
		Blockchain: config.Blockchain,
		Network:    strings.ToLower(configuration.Hedera.Mirror.Rosetta.Network),
		SubNetworkIdentifier: &types.SubNetworkIdentifier{
			Network: fmt.Sprintf("shard %s realm %s", configuration.Hedera.Mirror.Rosetta.Shard, configuration.Hedera.Mirror.Rosetta.Realm),
		},
	}

	version := &types.Version{
		RosettaVersion:    configuration.Hedera.Mirror.Rosetta.ApiVersion,
		NodeVersion:       configuration.Hedera.Mirror.Rosetta.NodeVersion,
		MiddlewareVersion: &configuration.Hedera.Mirror.Rosetta.Version,
	}

	dbClient := connectToDb(configuration.Hedera.Mirror.Rosetta.Db)
	defer dbClient.Close()

	asserter, err := asserter.NewServer(
		[]string{config.OperationTypeCryptoTransfer},
		true,
		[]*types.NetworkIdentifier{network},
	)
	if err != nil {
		log.Fatal(err)
	}

	router := NewBlockchainRouter(network, asserter, version, dbClient)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	log.Printf("Listening on port %s\n", configuration.Hedera.Mirror.Rosetta.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", configuration.Hedera.Mirror.Rosetta.Port), corsRouter))
}
