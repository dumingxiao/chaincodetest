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
	} else if function == "verify" {
		if len(args) == 4 { //verify the request of the  hespital
			return t.P2Verify(stub, args)
		} else if len(args) == 5 { //verify the request of the  celeres
			return t.P3Verify(stub, args)
		}
	}

	return nil, errors.New("Received unknown function invocation")
}

//add an new Authorization code
func (t *SimpleChaincode) Add(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	account := args[0]
	code := args[1]
	//accountTest, err := stub.GetState(account)

	//test if the account has been existed
	// if accountTest != nil {
	// 	return nil, errors.New("the ccount is existed")
	// }
	// add the account
	err := stub.PutState(account, []byte(code))
	if err != nil {
		return nil, errors.New("Failed to add the account")
	}

	return nil, nil
}

// 医院请求获取病人的病例数据
//发送 医院编号|| 医院授权码 || 病人编号 ||病人授权码

func (t *SimpleChaincode) P2Verify(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	p2 := args[0] //hospital
	p2Code := args[1]
	p1 := args[2] //patient
	p1Code := args[3]
	p2Test, err := stub.GetState(p2)
	p1Test, errs := stub.GetState(p1)
	//test if the account has been existed
	if err != nil {
		return nil, errors.New("error in reading hospital's code")
	}
	if errs != nil {
		return nil, errors.New("error in reading patient's code")
	}
	if p2Code != string(p2Test) {
		return nil, errors.New("hospital's code is error")
	} else if p1Code != string(p1Test) {
		return nil, errors.New("patient's code is error")
	}

	//
	err = stub.PutState(p1, []byte("used code"))
	if err != nil {
		return nil, errors.New("Failed to store the verification record")
	}
	return []byte("ok"), nil

}

// 医疗分析机构请求获取病人的病例数据，
//发送 医疗分析机构编号 || 医院编号 || 医院授权码 || 病人编号 ||病人授权码
func (t *SimpleChaincode) P3Verify(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	p3 := args[0]
	p2 := args[1]
	p2Code := args[2]
	p1 := args[3]
	p1Code := args[4]
	p2Test, err := stub.GetState(p2)
	p1Test, errs := stub.GetState(p1)
	//test if the account has been existed
	if err != nil {
		return nil, errors.New("error in reading hospital's code")
	}
	if errs != nil {
		return nil, errors.New("error in reading patient's code")
	}

	if p2Code != string(p2Test) {
		return nil, errors.New("hospital's code is error")
	} else if p1Code != string(p1Test) {
		return nil, errors.New("patient's code is error")
	}

	// reset the account's password
	err = stub.PutState(p1, []byte("used code"))
	err = stub.PutState(p3, []byte(p1))
	if err != nil {
		return nil, errors.New("Failed to store the verification record")
	}
	return []byte("ok"), nil

}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//account := args[0]
	// Handle different functions
	//password, err := stub.GetState(account) //get the var from chaincode state

	if function == "verify" {
		if len(args) == 4 { //verify the request of the  hespital
			return t.P2VerifyQuery(stub, args)
		} else if len(args) == 5 { //verify the request of the  celeres
			return t.P3VerifyQuery(stub, args)
		}
	}
	if function == "test" { //add a new Administrator account
		return t.Test(stub, args)
	}

	return nil, errors.New("failed to query")

}

func (t *SimpleChaincode) P2VerifyQuery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	p2 := args[0] //hospital
	p2Code := args[1]
	p1 := args[2] //patient
	p1Code := args[3]
	p2Test, err := stub.GetState(p2)
	p1Test, errs := stub.GetState(p1)
	//test if the account has been existed
	if err != nil {
		return nil, errors.New("error in reading hospital's code")
	}
	if errs != nil {
		return nil, errors.New("error in reading patient's code")
	}
	if p2Code != string(p2Test) {
		return nil, errors.New("hospital's code is error")
	} else if p1Code != string(p1Test) {
		return nil, errors.New("patient's code is error")
	}

	return []byte("ok"), nil

}

// 医疗分析机构请求获取病人的病例数据，
//发送 医疗分析机构编号 || 医院编号 || 医院授权码 || 病人编号 ||病人授权码
func (t *SimpleChaincode) P3VerifyQuery(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//p3 := args[0]
	p2 := args[1]
	p2Code := args[2]
	p1 := args[3]
	p1Code := args[4]
	p2Test, err := stub.GetState(p2)
	p1Test, errs := stub.GetState(p1)
	//test if the account has been existed
	if err != nil {
		return nil, errors.New("error in reading hospital's code")
	}
	if errs != nil {
		return nil, errors.New("error in reading patient's code")
	}

	if p2Code != string(p2Test) {
		return nil, errors.New("hospital's code is error")
	} else if p1Code != string(p1Test) {
		return nil, errors.New("patient's code is error")
	}

	return []byte("ok"), nil

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
