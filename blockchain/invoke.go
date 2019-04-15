package blockchain

import (
	"bytes"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// InvokePutFile will call the chaincode to put the file
func (setup *FabricSetup) InvokePutFile(value []byte, filename string) (string, error) {
	// Add value into ipfs
	cid, err := setup.sh.Add(bytes.NewReader(value))
	if err != nil {
		return "", err
	}

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "invoke")
	args = append(args, "putFile")
	args = append(args, cid)
	args = append(args, filename)
	args = append(args, setup.UserName)

	eventID := "eventInvokePutFile"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in put file invoke")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{
		ChaincodeID:  setup.ChainCodeID,
		Fcn:          args[0],
		TransientMap: transientDataMap,
		Args: [][]byte{
			[]byte(args[1]),
			[]byte(args[2]),
			[]byte(args[3]),
			[]byte(args[4]),
			[]byte(args[5]),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	return string(response.TransactionID), nil
}
