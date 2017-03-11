/*
This Chaincode stores the accounts and passwords

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type RequestState struct {
	CeleresID  string `json:"celeresID"`  //user who created the open trade order
	HospitalID string `json:"hospitalID"` //utc timestamp of creation
	Hcode      string `json:"hcode"`      //description of desired marble
	PatientID  string `json:"patientID"`
	Pcode      string `json:"pcode"` //array of marbles willing to trade away
}

// the Init function is used for deploying the chaincode and setting the Administrator account

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

	// Write the password to the ledger
	err = stub.PutState(ID, []byte(IDval))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" { //add a new Administrator account
		return t.Init(stub, "init", args)
	} else if function == "add" { //add a new account
		return t.Add(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

//add an new Authorization code
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	requestID := args[0]
	var requestState RequestState
	requestState.CeleresID = args[1]
	requestState.HospitalID = args[2]
	requestState.Hcode = args[3]
	requestState.PatientID = args[4]
	requestState.Pcode = args[5]
	JsonRequestState, _ := json.Marshal(requestState)

	// var requestState []string
	// requestState[0] = args[1]
	// requestState[1] = args[2]
	// requestState[2] = args[3]
	// requestState[3] = args[4]
	// requestState[4] = args[5]
	// JsonRequestState, _ := json.Marshal(requestState)
	//accountTest, err := stub.GetState(account)

	//test if the account has been existed
	// if accountTest != nil {
	// 	return nil, errors.New("the ccount is existed")
	// }
	// add the account
	err := stub.PutState(requestID, JsonRequestState)
	if err != nil {
		return nil, errors.New("Failed to add the account")
	}

	return nil, nil
}

// 医院请求获取病人的病例数据
//发送 医院编号|| 医院授权码 || 病人编号 ||病人授权码

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//account := args[0]
	// Handle different functions
	//password, err := stub.GetState(account) //get the var from chaincode state

	if function == "verify" {
		return t.VerifyQuery(stub, args)
	} else if function == "test" { //add a new Administrator account
		return t.Test(stub, args)
	}

	return nil, errors.New("failed to query")

}

func (t *SimpleChaincode) VerifyQuery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	celeres := args[0]
	result, err := stub.GetState(celeres)
	//test if the account has been existed
	if err != nil {
		return nil, errors.New("error in reading request record")
	}

	return result, nil

}

func (t *SimpleChaincode) Test(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//fmt.Println("query is running " + function)
	account := args[0]
	// Handle different functions
	password, err := stub.GetState(account) //get the var from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + account + "\"}"
		return nil, errors.New(jsonResp)
	}

	return password, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
