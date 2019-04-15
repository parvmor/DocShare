package blockchain

import (
	"bytes"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// InvokePutFile will call the chaincode to put the file
func (setup *FabricSetup) InvokePutFile(value []byte, filename, username string) (string, error) {
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
	args = append(args, username)

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

// InvokeShareFile will call the chaincode to share the file
// Assume public key cryptography to be done at controller
// Following is the mechanism:
// - file is encrpted with a randomly chosen key
// - encrypted file is stored in IPFS at cid_file
// - key || cid_file is encrypted using public key of receiver
// - this is value. this is to be stored at sharer_receiver_filename
func (setup *FabricSetup) InvokeShareFile(value []byte, filename, sharer, receiver string) (string, error) {
	// Add value into ipfs
	cid, err := setup.sh.Add(bytes.NewReader(value))
	if err != nil {
		return "", err
	}

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "invoke")
	args = append(args, "shareFile")
	args = append(args, cid)
	args = append(args, filename)
	args = append(args, sharer)
	args = append(args, receiver)

	eventID := "eventInvokeShareFile"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in share file invoke")

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
			[]byte(args[6]),
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
