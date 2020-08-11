package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistance/postgres/block"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistance/postgres/transaction"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/services"
	"github.com/jinzhu/gorm"
)

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(network *types.NetworkIdentifier, asserter *asserter.Asserter, version *types.Version, dbClient *gorm.DB) http.Handler {
	blockRepo := block.NewBlockRepository(dbClient)
	transactionRepo := transaction.NewTransactionRepository(dbClient)

	networkAPIService := services.NewNetworkAPIService(network, version, blockRepo, transactionRepo)
	networkAPIController := server.NewNetworkAPIController(networkAPIService, asserter)

	blockAPIService := services.NewBlockAPIService(network, blockRepo, transactionRepo)
	blockAPIController := server.NewBlockAPIController(blockAPIService, asserter)

	mempoolAPIService := services.NewMempoolAPIService()
	mempoolAPIController := server.NewMempoolAPIController(mempoolAPIService, asserter)

	return server.NewRouter(networkAPIController, blockAPIController, mempoolAPIController)
}

func main() {
	config := LoadConfig()

	network := &types.NetworkIdentifier{
		Blockchain: "Hedera",
		Network:    strings.ToLower(config.Hedera.Mirror.Rosetta.Network),
		SubNetworkIdentifier: &types.SubNetworkIdentifier{
			Network: fmt.Sprintf("shard %s realm %s", config.Hedera.Mirror.Rosetta.Shard, config.Hedera.Mirror.Rosetta.Realm),
		},
	}

	version := &types.Version{
		RosettaVersion:    config.Hedera.Mirror.Rosetta.ApiVersion,
		NodeVersion:       config.Hedera.Mirror.Rosetta.NodeVersion,
		MiddlewareVersion: &config.Hedera.Mirror.Rosetta.Version,
	}

	dbClient := connectToDb(config.Hedera.Mirror.Rosetta.Db)
	defer dbClient.Close()

	asserter, err := asserter.NewServer(
		[]string{"Transfer"},
		false,
		[]*types.NetworkIdentifier{network},
	)
	if err != nil {
		log.Fatal(err)
	}

	router := NewBlockchainRouter(network, asserter, version, dbClient)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	log.Printf("Listening on port %s\n", config.Hedera.Mirror.Rosetta.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Hedera.Mirror.Rosetta.Port), corsRouter))
}
