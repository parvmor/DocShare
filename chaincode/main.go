package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// DocShareChainCode implementation of Chaincode
type DocShareChainCode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *DocShareChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### DocShareChainCode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke - All future requests named invoke will arrive here.
func (t *DocShareChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### DocShareChainCode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *DocShareChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### DocShareChainCode query ###########")

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
	// We also have only one possible argument: getFile
	if args[1] == "getFile" && len(args) == 4 {
		// Key will be username_filename
		key := args[3] + "_" + args[2]

		// Get the state of the value matching the key hello in the ledger
		state, err := stub.GetState(key)
		if err != nil {
			return shim.Error("Failed to get state of " + key)
		}

		// Return this value in response
		return shim.Success(state)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *DocShareChainCode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### DocShareChainCode invoke ###########")

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Check if the ledger key is "hello" and process if it is the case. Otherwise it returns an error.
	if args[1] == "putFile" && len(args) == 5 {
		// Key will be username_filename
		key := args[4] + "_" + args[3]

		// Write the new value in the ledger
		err := stub.PutState(key, []byte(args[2]))
		if err != nil {
			return shim.Error("Failed to update state of " + key)
		}

		// Notify listeners that an event "eventInvokePutFile" have been executed
		err = stub.SetEvent("eventInvokePutFile", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return this value in response
		return shim.Success(nil)
	} else if args[1] == "shareFile" && len(args) == 6 {
		// Key will be username_filename
		key := args[5] + "_" + args[4] + "_" + args[3]

		// Write the new value in the ledger
		err := stub.PutState(key, []byte(args[2]))
		if err != nil {
			return shim.Error("Failed to update state of " + key)
		}

		// Notify listeners that an event "eventInvokePutFile" have been executed
		err = stub.SetEvent("eventInvokeShareFile", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return this value in response
		return shim.Success(nil)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown invoke action, check the second argument.")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(DocShareChainCode))
	if err != nil {
		fmt.Printf("Error starting DocShare Service chaincode: %s", err)
	}
}
