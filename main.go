package main

import (
	"fmt"
	"os"

	"github.com/parvmor/docshare/blockchain"
)

const totalUsers int = 100

var fSetupMap map[string]blockchain.FabricSetup

func main() {
	// Definition of the Fabric SDK properties
	for i := 1; i <= totalUsers; i++ {
		userName := "User" + string(i)
		fSetupMap[userName] = blockchain.FabricSetup{
			// Network parameters
			OrdererID: "orderer.hf.docshare.io",

			// Channel parameters
			ChannelID:     "docshare",
			ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/parvmor/docshare/network/artifacts/docshare.channel.tx",

			// Chaincode parameters
			ChainCodeID:     "docshare",
			ChaincodeGoPath: os.Getenv("GOPATH"),
			ChaincodePath:   "github.com/parvmor/docshare/chaincode/",
			OrgAdmin:        "Admin",
			OrgName:         "org1",
			ConfigFile:      "config.yaml",

			// User parameters
			UserName: userName,
		}

		// Initialize the SDK
		err := fSetupMap[userName].Initialize()
		if err != nil {
			fmt.Printf("Unable to initialize the Fabric SDK: %v for %s\n", err, userName)
			return
		}
		// Close SDK
		defer fSetupMap[userName].CloseSDK()

		// Install and instantiate the chaincode
		err = fSetup.InstallAndInstantiateCC()
		if err != nil {
			fmt.Printf("Unable to install and instantiate the chaincode: %v for %s\n", err, userName)
			return
		}
	}
}
