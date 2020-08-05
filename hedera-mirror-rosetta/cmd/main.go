package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/persistance/postgres/block"
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/app/services"
	"github.com/jinzhu/gorm"
)

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(network *types.NetworkIdentifier, asserter *asserter.Asserter, dbClient *gorm.DB) http.Handler {

	blockRepo := block.NewBlockRepository(dbClient)
	blockAPIService := services.NewBlockAPIService(network, blockRepo)
	blockAPIController := server.NewBlockAPIController(blockAPIService, asserter)

	return server.NewRouter(blockAPIController)
}

func main() {
	config := LoadConfig()

	network := &types.NetworkIdentifier{
		Blockchain: "Hedera",
		Network:    config.Hedera.Mirror.Rosetta.Network,
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

	router := NewBlockchainRouter(network, asserter, dbClient)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	log.Printf("Listening on port %s\n", config.Hedera.Mirror.Rosetta.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Hedera.Mirror.Rosetta.Port), corsRouter))
}
