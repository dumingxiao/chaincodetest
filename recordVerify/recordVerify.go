// ============================================================================================================================
// This smart contract is used for storing and verifying the original medical record
// It includes creating the Administrator account ,adding the original medical record
// and verifying the aoriginal medical record
// ============================================================================================================================

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

type RecordDetail struct {
	DbIP       string `json:"dbIP"`       //user who created the open trade order
	RecordHash string `json:"recordHash"` //utc timestamp of creation
}

// ============================================================================================================================
// Init function is used for creating the Administrator account
// ============================================================================================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var ID string    // Administrator's ID
	var IDval string // Password of the Administrator ID
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. ")
	}

	// Initialize the chaincode
	ID = args[0]
	IDval = args[1]

	// Write the state to the ledger
	err = stub.PutState(ID, []byte(IDval))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Invoke function is the entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" { //add a new Administrator account
		return t.Init(stub, "init", args)

	} else if function == "add" { //add a new account
		return t.Add(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

//add a new record
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	recordID := args[0]
	var recordDetail RecordDetail
	recordDetail.DbIP = args[1]
	recordDetail.RecordHash = args[2]
	JsonRecordDetail, _ := json.Marshal(recordDetail)
	recordTest, err := stub.GetState(recordID)

	//test if the account has been existed
	if recordTest != nil {
		return nil, errors.New("the record is existed")
	}
	// add the hash
	err = stub.PutState(recordID, JsonRecordDetail)
	if err != nil {
		return nil, errors.New("Failed to add the record")
	}

	return nil, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "verify" { //deletes an account from its state
		if len(args) == 1 {
			return t.Get(stub, args)
		} else if len(args) == 2 {
			return t.VerifyRecordHash(stub, args)

		}
	}

	return nil, errors.New("failed to query")

}

func (t *SimpleChaincode) VerifyRecordHash(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	recordID := args[0]
	recordHash := args[1]
	recordTest, err := stub.GetState(recordID)
	var JsonRecordTest RecordDetail
	json.Unmarshal(recordTest, &JsonRecordTest)
	ver := []byte("ok")
	jsonResp := "{\"Error\":\"Failed to get state for " + recordID + "\"}"
	//test if the account has been existed
	if err != nil {
		return nil, errors.New(jsonResp)
	}
	if recordTest == nil {
		return nil, errors.New(jsonResp)
	}

	// verify
	if recordHash == string(JsonRecordTest.RecordHash) {
		return ver, nil

	} else {
		return nil, errors.New("Failed to verify the record")
	}

}

func (t *SimpleChaincode) Get(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	recordID := args[0]
	// Handle different functions
	recordDetail, err := stub.GetState(recordID) //get the var from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + recordID + "\"}"
		return nil, errors.New(jsonResp)
	}

	return recordDetail, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
