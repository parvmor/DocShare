package blockchain

import (
	"fmt"
	"io/ioutil"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// QueryGetFile query the chaincode to fetch the given filename
func (setup *FabricSetup) QueryGetFile(filename string) ([]byte, error) {
	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "getFile")
	args = append(args, filename)
	args = append(args, setup.UserName)

	response, err := setup.client.Query(channel.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn:         args[0],
		Args: [][]byte{
			[]byte(args[1]),
			[]byte(args[2]),
			[]byte(args[3]),
			[]byte(args[4]),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	// Payload is the IPFS cid
	cid := string(response.Payload)
	reader, err := setup.sh.Cat(cid)
	if err != nil {
		return nil, err
	}
	value, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return value, nil
}
