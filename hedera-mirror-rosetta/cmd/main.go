package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

func NewBlockchainRouter(
	network *types.NetworkIdentifier,
	asserter *asserter.Asserter,
) http.Handler {
	return server.NewRouter()
}

func main() {
	config := LoadConfig()

	network := &types.NetworkIdentifier{
		Blockchain: "Hedera",
		Network:    config.Hedera.Mirror.Rosetta.Network,
	}

	asserter, err := asserter.NewServer(
		[]string{"Transfer"},
		false,
		[]*types.NetworkIdentifier{network},
	)
	if err != nil {
		log.Fatal(err)
	}

	router := NewBlockchainRouter(network, asserter)
	loggedRouter := server.LoggerMiddleware(router)
	corsRouter := server.CorsMiddleware(loggedRouter)
	log.Printf("Listening on port %s\n", config.Hedera.Mirror.Rosetta.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Hedera.Mirror.Rosetta.Port), corsRouter))
}
