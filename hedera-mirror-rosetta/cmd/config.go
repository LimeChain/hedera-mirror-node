package main

import (
	"github.com/hashgraph/hedera-mirror-node/hedera-mirror-rosetta/types"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func GetConfig() *types.Config {
	filename, _ := filepath.Abs("config/application.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var config types.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
