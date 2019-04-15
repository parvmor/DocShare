package main

import (
	"fmt"

	"github.com/parvmor/docshare/blockchain"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: "orderer.hf.docshare.io",

		// Channel parameters
		ChannelID:     "docshare",
		ChannelConfig: "/home/centos/go/src/github.com/parvmor/docshare/network/artifacts/docshare.channel.tx",

		// Chaincode parameters
		ChainCodeID:     "docshare",
		ChaincodeGoPath: "/home/centos/go",
		ChaincodePath:   "github.com/parvmor/docshare/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: "User1",
	}

	// Initialize the SDK
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}
}
