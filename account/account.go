/*
This Chaincode stores the accounts and passwords

*/

package main

import (
	"errors"
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
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

	// Write the state to the ledger
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
	} else if function == "delete" { //deletes an account from its state
		return t.Delete(stub, args)
	} else if function == "add" { //add a new account
		return t.Add(stub, args)
	} else if function == "edit" { //change the password of the account
		return t.Edit(stub, args)
	} else if function == "reset" { // reset the password of the account
		return t.Reset(stub, args)
	} else if function == "verify" { //verify the account and password
		return t.Verify(stub, args)
	}

	return nil, errors.New("Received unknown function invocation")
}

// delete an account
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	account := args[0]
	err := stub.DelState(account) //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete account")
	}

	return nil, nil
}

//add an new account
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	account := args[0]
	password := args[1]
	accountTest, err := stub.GetState(account)

	//test if the account has been existed
	if accountTest != nil {
		return nil, errors.New("the ccount is existed")
	}
	// add the account
	err = stub.PutState(account, []byte(password))
	if err != nil {
		return nil, errors.New("Failed to add the account")
	}

	return nil, nil
}

//change the account's password
func (t *SimpleChaincode) Edit(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	account := args[0]
	password := args[2]
	accountTest, err := stub.GetState(account)

	//test if the account has been existed
	if err != nil {
		return nil, errors.New("account not found")
	}
	if accountTest == nil {
		return nil, errors.New("account not found")
	}

	// edit the account
	err = stub.PutState(account, []byte(password))
	if err != nil {
		return nil, errors.New("Failed to edit the account")
	}

	return nil, nil
}

func (t *SimpleChaincode) Reset(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	account := args[0]
	password := args[0]
	accountTest, err := stub.GetState(account)

	//test if the account has been existed
	if err != nil {
		return nil, errors.New("account is not here")
	}
	if accountTest == nil {
		return nil, errors.New("account not found")
	}

	// reset the account's password
	err = stub.PutState(account, []byte(password))
	if err != nil {
		return nil, errors.New("Failed to edit the account")
	}

	return nil, nil
}

func (t *SimpleChaincode) Verify(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	account := args[0]
	password := args[1]
	accountTest, err := stub.GetState(account)
	ver := []byte("ok")
	//test if the account has been existed
	if err != nil {
		return nil, errors.New("account not found")
	}
	if accountTest == nil {
		return nil, errors.New("account not found")
	}

	// verify
	if password == string(accountTest) {
		return ver, nil

	} else {
		return nil, errors.New("Failed to verify the account")
	}

}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
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
