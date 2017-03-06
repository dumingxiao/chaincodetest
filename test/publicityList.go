/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//两公示一公告名单

import (
	"errors"
	"fmt"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var tableTime string  // time of the publicity list
	var operatorID string // ID of the operator
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. ")
	}

	// Initialize the chaincode
	tableTime = args[0]
	operatorID = args[1]

	// Write the state to the ledger
	err = stub.PutState(tableTime, []byte(operatorID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. ")
	}
	applicantID := args[0]
	personalInfo := args[1]
	var err error
	// Perform the execution

	// Write the state back to the ledger
	if function == "add" {
		// Perform the execution

		// Write the state back to the ledger
		err = stub.PutState(applicantID, []byte(personalInfo))
		if err != nil {
			return nil, err
		}
	}
	if function == "update" {
		ListIDvalTemp, errs := stub.GetState(applicantID)

		if errs != nil {
			return nil, errors.New("list is not here")
		}
		if ListIDvalTemp == nil {
			return nil, errors.New("Entity not found")
		}
		err = stub.PutState(applicantID, []byte(personalInfo))
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var applicantID string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	applicantID = args[0]

	// Get the state from the ledger
	personalInfo, err := stub.GetState(applicantID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + applicantID + "\"}"
		return nil, errors.New(jsonResp)
	}

	if personalInfo == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + applicantID + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + applicantID + "\",\"Amount\":\"" + string(personalInfo) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return personalInfo, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
