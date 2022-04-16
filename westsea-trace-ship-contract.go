/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// WestseaTraceShipContract contract for managing CRUD for WestseaTraceShip
type WestseaTraceShipContract struct {
	contractapi.Contract
}

/*
 ****************************************
 ****************************************
 * ProductSerialNum TRANSCATION METHDOS *
 ****************************************
 ****************************************
 */

// ProductSerialNumExists returns true when asset with given ID exists in world state
func (c *WestseaTraceShipContract) ProductSerialNumExists(ctx contractapi.TransactionContextInterface, productSerialNumID string) (bool, error) {
	data, err := ctx.GetStub().GetState(productSerialNumID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateProductSerialNum creates a new instance of WestseaTraceShip
func (c *WestseaTraceShipContract) CreateProductSerialNum(ctx contractapi.TransactionContextInterface, 
	productSerialNumID string, 
	serialNumber string,
	designation string,
	productType string,
	documentKeys []string,
	) (string, error) {

	exists, err := c.ProductSerialNumExists(ctx, productSerialNumID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return "", fmt.Errorf("the lot %s already exists", productSerialNumID)
	}

	productSerialNum := &ProductSerialNum{
		DocType: "productSerialNum",
		ID:      productSerialNumID,
		SerialNumber: serialNumber,
		Designation: designation,
		ProductType: productType,
		DocumentKeys: documentKeys,
	}

	bytes, err := json.Marshal(productSerialNum)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(productSerialNum.ID, bytes)
	if err != nil {
		return "", fmt.Errorf("failed to put to world state: %v", err)
	}

	return fmt.Sprintf("%s created successfully", productSerialNumID), nil
}

// ReadProductSerialNum retrieves an instance of ProductSerialNum from the world state
func (c *WestseaTraceShipContract) ReadProductSerialNum(ctx contractapi.TransactionContextInterface, productSerialNumID string) (*ProductSerialNum, error) {
	exists, err := c.ProductSerialNumExists(ctx, productSerialNumID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", productSerialNumID)
	}

	bytes, _ := ctx.GetStub().GetState(productSerialNumID)

	productSerialNum := new(ProductSerialNum)

	err = json.Unmarshal(bytes, productSerialNum)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type WestseaTraceShip")
	}

	return productSerialNum, nil
}

// UpdateProductSerialNumDocumentKeys retrieves an instance of ProductSerialNum from the world state and updates its document keys
func (c *WestseaTraceShipContract) UpdateProductSerialNumDocumentKeys(ctx contractapi.TransactionContextInterface, productSerialNumID string, newDocumentKeys []string) (string, error) {
	exists, err := c.ProductSerialNumExists(ctx, productSerialNumID)
	if err != nil {
		return "", fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return "", fmt.Errorf("The asset %s does not exist", productSerialNumID)
	}

	outdatedProductSerialNumBytes, _ := ctx.GetStub().GetState(productSerialNumID) // Gets "old" ProductSerialNum bytes from productSerialNumID

	outdatedProductSerialNum := new(ProductSerialNum) // Initialize outdated/"old" ProductSerialNum object

	// Parses the JSON-encoded data in bytes (outdatedProductSerialNumBytes) to the "old" ProductSerialNum object (outdatedProductSerialNum)
	err = json.Unmarshal(outdatedProductSerialNumBytes, outdatedProductSerialNum)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal world state data to type ProductSerialNum")
	}

	productSerialNum := &ProductSerialNum{
		DocType: "productSerialNum",
		ID:      productSerialNumID,
		SerialNumber: outdatedProductSerialNum.SerialNumber,
		Designation: outdatedProductSerialNum.Designation,
		ProductType: outdatedProductSerialNum.ProductType,
		DocumentKeys: newDocumentKeys,
	}

	bytes, _ := json.Marshal(productSerialNum)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(productSerialNumID, bytes)
	if err != nil {
		return "", fmt.Errorf("failed to put to world state: %v", err)
	}

	return fmt.Sprintf("%s document keys updated successfully", productSerialNumID), nil
}

// DeleteProductSerialNum deletes an instance of ProductSerialNum from the world state
func (c *WestseaTraceShipContract) DeleteProductSerialNum(ctx contractapi.TransactionContextInterface, productSerialNumID string) error {
	exists, err := c.ProductSerialNumExists(ctx, productSerialNumID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", productSerialNumID)
	}

	return ctx.GetStub().DelState(productSerialNumID)
}



 /*
  ****************************************
  ****************************************
  * Activity TRANSCATION METHDOS *
  ****************************************
  ****************************************
  */
 
 // ActivityExists returns true when asset with given ID exists in world state
 func (c *WestseaTraceShipContract) ActivityExists(ctx contractapi.TransactionContextInterface, activityID string) (bool, error) {
	data, err := ctx.GetStub().GetState(activityID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateActivity creates a new instance of WestseaTraceShip
func (c *WestseaTraceShipContract) CreateActivity(ctx contractapi.TransactionContextInterface, 
	activityID string, 
	designation string,
	userId string,
	dateTime string,
	inputProductSerialNums []string,
	outputProductSerialNum ProductSerialNum,
	) (string, error) {

	exists, err := c.ActivityExists(ctx, activityID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return "", fmt.Errorf("the asset %s already exists", activityID)
	}

	//the input products must exist
	for _, inputID := range inputProductSerialNums {

	   exists, err = c.ProductSerialNumExists(ctx, inputID)
	   if err != nil {
		   return "", fmt.Errorf("Could not read from world state. %s", err)
	   } else if !exists {
		   return "", fmt.Errorf("The input product [%s] does not exists", inputID)
	   }
   }

	//FIXME: O produto output deve ser criado previamente? Ou durante a criação de uma atividade?
	exists, err = c.ActivityExists(ctx, outputProductSerialNum.ID)
	if err != nil {
		return "", fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return "", fmt.Errorf("The output product [%s] already exists", outputProductSerialNum.ID)
	}

	_, err = c.CreateProductSerialNum(ctx, 
		outputProductSerialNum.ID,
		outputProductSerialNum.SerialNumber,
		outputProductSerialNum.Designation,
		outputProductSerialNum.ProductType,
		outputProductSerialNum.DocumentKeys,
	)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	}

	activity := &Activity{
		DocType: "activity",
		ID:      activityID,
		Designation: designation,
		UserId: userId,
		DateTime: dateTime,
		InputProductSerialNums: inputProductSerialNums,
		OutputProductSerialNum: outputProductSerialNum,
	}

	bytes, err := json.Marshal(activity)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(activity.ID, bytes)
	if err != nil {
		return "", fmt.Errorf("failed to put to world state: %v", err)
	}

	return fmt.Sprintf("%s created successfully", activityID), nil
}

// ReadActivity retrieves an instance of Activity from the world state
func (c *WestseaTraceShipContract) ReadActivity(ctx contractapi.TransactionContextInterface, activityID string) (*Activity, error) {
	exists, err := c.ActivityExists(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", activityID)
	}

	bytes, _ := ctx.GetStub().GetState(activityID)

	activity := new(Activity)

	err = json.Unmarshal(bytes, activity)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type WestseaTraceShip")
	}

	return activity, nil
}

//TODO: should we be able to update any property of the activity?

// DeleteActivity deletes an instance of Activity from the world state
func (c *WestseaTraceShipContract) DeleteActivity(ctx contractapi.TransactionContextInterface, activityID string) error {
	exists, err := c.ActivityExists(ctx, activityID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", activityID)
	}

	return ctx.GetStub().DelState(activityID)
}
