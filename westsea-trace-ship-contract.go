/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// WestseaTraceShipContract contract for managing CRUD for WestseaTraceShip
type WestseaTraceShipContract struct {
	contractapi.Contract
}

/*
 ****************************************
 ****************************************
 * ProductLot TRANSCATION METHDOS *
 ****************************************
 ****************************************
 */

// GetAllProductLot queries for all productLots.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (docType).
// Only available on state databases that support rich query (e.g. CouchDB)
// Example: Parameterized rich query
func (c *WestseaTraceShipContract) GetAllProductLot(ctx contractapi.TransactionContextInterface) ([]*ProductLot, error) {
	queryString := buildQueryString("docType", "productLot")
	productLots, _, err := getQueryResultForQueryString(ctx, queryString, IterationType(0))
	return productLots, err
}

// GetAllActivities queries for all activities.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (docType).
// Only available on state databases that support rich query (e.g. CouchDB)
// Example: Parameterized rich query
func (c *WestseaTraceShipContract) GetAllActivities(ctx contractapi.TransactionContextInterface) ([]*Activity, error) {
	queryString := buildQueryString("docType", "activity")
	_, activities, err := getQueryResultForQueryString(ctx, queryString, IterationType(1))
	return activities, err
}

//FIXME: order by date?
func (c *WestseaTraceShipContract) GetTraceability(ctx contractapi.TransactionContextInterface, referenceNum string) ([]*Activity, error) {
	//get product
	product, err := c.ReadProductLotByReferenceNum(ctx, referenceNum)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	}

	//get all activities
	allActivities, err := c.GetAllActivities(ctx)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	}

	var tracedActivities []*Activity

	for _, activity := range allActivities {
		//filter for the activities that have the outputLot == product
		if activity.OutputProductLot.ID == product.ID {
			tracedActivities = append(tracedActivities, activity)
		}

		//filter for the activities that have the product.ID == inputLots.key
		for inputID := range activity.InputProductLots {
			if inputID == product.ID {
				tracedActivities = append(tracedActivities, activity)
			}
		}
	}

	return tracedActivities, nil
}

// ProductLotExists returns true when asset with given ID exists in world state
func (c *WestseaTraceShipContract) ProductLotExists(ctx contractapi.TransactionContextInterface, productLotID string) (bool, error) {
	data, err := ctx.GetStub().GetState(productLotID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateProductLot creates a new instance of WestseaTraceShip
func (c *WestseaTraceShipContract) CreateProductLot(ctx contractapi.TransactionContextInterface,
	productLotID string,
	referenceNumber string,
	isSerialNumber bool,
	designation string,
	productType string,
	initialAmount float32,
	documentKeys []string,
) (string, error) {

	exists, err := c.ProductLotExists(ctx, productLotID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return "", fmt.Errorf("the productLot %s already exists", productLotID)
	}

	_, err = c.ReadProductLotByReferenceNum(ctx, referenceNumber)

	if err == nil {
		return "", fmt.Errorf("the productLot with the reference number %s already exists", referenceNumber)
	}

	//the referenceNumber of a productLot can be the serial number, if isSerialNumber; or the lotNumber if !isSerialNumber)
	if isSerialNumber {
		initialAmount = 1
	}

	productLot := &ProductLot{
		DocType:           "productLot",
		ID:                productLotID,
		ReferenceNumber:   referenceNumber,
		IsSerialNumber:    isSerialNumber,
		Designation:       designation,
		ProductType:       productType,
		InitialQuantity:   initialAmount,
		AvailableQuantity: initialAmount,
		DocumentKeys:      documentKeys,
	}

	bytes, err := json.Marshal(productLot)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(productLot.ID, bytes)
	if err != nil {
		return "", fmt.Errorf("failed to put to world state: %v", err)
	}

	return fmt.Sprintf("%s created successfully", productLotID), nil
}

// ReadProductLot retrieves an instance of ProductLot from the world state
func (c *WestseaTraceShipContract) ReadProductLot(ctx contractapi.TransactionContextInterface, productLotID string) (*ProductLot, error) {
	exists, err := c.ProductLotExists(ctx, productLotID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", productLotID)
	}

	bytes, _ := ctx.GetStub().GetState(productLotID)

	productLot := new(ProductLot)

	err = json.Unmarshal(bytes, productLot)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type WestseaTraceShip")
	}

	return productLot, nil
}

// ReadProductLot retrieves an instance of ProductLot from the world state
func (c *WestseaTraceShipContract) ReadProductLotByReferenceNum(ctx contractapi.TransactionContextInterface, referenceNum string) (*ProductLot, error) {
	queryString := buildQueryString("referenceNumber", referenceNum)
	productLots, _, err := getQueryResultForQueryString(ctx, queryString, IterationType(0))

	if err != nil {
		return nil, fmt.Errorf("Could not read productLot with reference number: %s", referenceNum)
	}

	if len(productLots) <= 0 {
		return nil, fmt.Errorf("ProductLot with reference number: %s was not found", referenceNum)
	}

	productLotFound := productLots[0]

	exists, err := c.ProductLotExists(ctx, productLotFound.ID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", productLotFound.ID)
	}

	bytes, _ := ctx.GetStub().GetState(productLotFound.ID)

	productLot := new(ProductLot)

	err = json.Unmarshal(bytes, productLot)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type WestseaTraceShip")
	}

	return productLot, nil
}

// UpdateProductLotDocumentKeys retrieves an instance of ProductLot from the world state and updates its document keys
func (c *WestseaTraceShipContract) UpdateProductLotDocumentKeys(ctx contractapi.TransactionContextInterface, productLotID string, newDocumentKeys []string) (string, error) {
	exists, err := c.ProductLotExists(ctx, productLotID)
	if err != nil {
		return "", fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return "", fmt.Errorf("The asset %s does not exist", productLotID)
	}

	outdatedProductLotBytes, _ := ctx.GetStub().GetState(productLotID) // Gets "old" ProductLot bytes from productLotID

	outdatedProductLot := new(ProductLot) // Initialize outdated/"old" ProductLot object

	// Parses the JSON-encoded data in bytes (outdatedProductLotBytes) to the "old" ProductLot object (outdatedProductLot)
	err = json.Unmarshal(outdatedProductLotBytes, outdatedProductLot)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal world state data to type ProductLot")
	}

	productLot := &ProductLot{
		DocType:           "productLot",
		ID:                productLotID,
		ReferenceNumber:   outdatedProductLot.ReferenceNumber,
		IsSerialNumber:    outdatedProductLot.IsSerialNumber,
		Designation:       outdatedProductLot.Designation,
		ProductType:       outdatedProductLot.ProductType,
		InitialQuantity:   outdatedProductLot.InitialQuantity,
		AvailableQuantity: outdatedProductLot.AvailableQuantity,
		DocumentKeys:      newDocumentKeys,
	}

	bytes, _ := json.Marshal(productLot)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(productLotID, bytes)
	if err != nil {
		return "", fmt.Errorf("failed to put to world state: %v", err)
	}

	return fmt.Sprintf("%s document keys updated successfully", productLotID), nil
}

// UpdateProductAvailableAmount retrieves an instance of ProductLot from the world state and updates its available quantity
func (c *WestseaTraceShipContract) UpdateProductAvailableQuantity(ctx contractapi.TransactionContextInterface, productLotID string, newAvailableQuantity float32) (string, error) {
	exists, err := c.ProductLotExists(ctx, productLotID)
	if err != nil {
		return "", fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return "", fmt.Errorf("The asset %s does not exist", productLotID)
	}

	outdatedProductLotBytes, _ := ctx.GetStub().GetState(productLotID) // Gets "old" ProductLot bytes from productLotID

	outdatedProductLot := new(ProductLot) // Initialize outdated/"old" ProductLot object

	// Parses the JSON-encoded data in bytes (outdatedProductLotBytes) to the "old" ProductLot object (outdatedProductLot)
	err = json.Unmarshal(outdatedProductLotBytes, outdatedProductLot)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal world state data to type ProductLot")
	}

	productLot := &ProductLot{
		DocType:           "productLot",
		ID:                productLotID,
		ReferenceNumber:   outdatedProductLot.ReferenceNumber,
		IsSerialNumber:    outdatedProductLot.IsSerialNumber,
		Designation:       outdatedProductLot.Designation,
		ProductType:       outdatedProductLot.ProductType,
		InitialQuantity:   outdatedProductLot.InitialQuantity,
		AvailableQuantity: newAvailableQuantity,
		DocumentKeys:      outdatedProductLot.DocumentKeys,
	}

	bytes, _ := json.Marshal(productLot)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(productLotID, bytes)
	if err != nil {
		return "", fmt.Errorf("failed to put to world state: %v", err)
	}

	return fmt.Sprintf("%s available quantity updated successfully to %.2f", productLotID, newAvailableQuantity), nil
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
	inputProductLots map[string]float32,
	outputProductLot ProductLot,
) (string, error) {

	exists, err := c.ActivityExists(ctx, activityID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return "", fmt.Errorf("the asset %s already exists", activityID)
	}

	for inputID, usedAmount := range inputProductLots {

		//the input products must exist
		exists, err = c.ProductLotExists(ctx, inputID)
		if err != nil {
			return "", fmt.Errorf("Could not read from world state. %s", err)
		} else if !exists {
			return "", fmt.Errorf("The input product [%s] does not exists", inputID)
		}

		//get product
		product, err := c.ReadProductLot(ctx, inputID)
		if err != nil {
			return "", fmt.Errorf("Could not read from world state. %s", err)
		}

		//make sure the inputProduct quantity is valid quantities

		if usedAmount <= 0 {
			return "", fmt.Errorf("inputProduct amount must be greater than 0 (inputProduct amount for inputProduct [%s] is %.2f)", inputID, usedAmount)
		} else if usedAmount > product.AvailableQuantity {
			return "", fmt.Errorf("inputProduct amount must not exceed the inputProduct availableQuantity (inputProduct [%s] maximum quantity is %.2f)", inputID, product.AvailableQuantity)
		}

		//the amounts on inputLots should be reduced on the availableAmount of a productLot
		_, err = c.UpdateProductAvailableQuantity(ctx, inputID, product.AvailableQuantity-usedAmount)
		if err != nil {
			return "", fmt.Errorf("Could not write to world state. %s", err)
		}
	}

	exists, err = c.ProductLotExists(ctx, outputProductLot.ID)
	if err != nil {
		return "", fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return "", fmt.Errorf("The output product [%s] already exists", outputProductLot.ID)
	}

	_, err = c.CreateProductLot(ctx,
		outputProductLot.ID,
		outputProductLot.ReferenceNumber,
		outputProductLot.IsSerialNumber,
		outputProductLot.Designation,
		outputProductLot.ProductType,
		outputProductLot.InitialQuantity,
		outputProductLot.DocumentKeys,
	)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	}

	activity := &Activity{
		DocType:          "activity",
		ID:               activityID,
		Designation:      designation,
		UserId:           userId,
		DateTime:         time.Now().Format(time.RFC3339),
		InputProductLots: inputProductLots,
		OutputProductLot: outputProductLot,
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

/*
 ****************************************
 ****************************************
 * COMMON METHDOS *
 ****************************************
 ****************************************
 */

type IterationType int

const (
	PRODUCT_LOT IterationType = iota
	ACTIVITY
)

// constructQueryResponseFromIterator constructs a slice of lots from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface, t IterationType) ([]*ProductLot, []*Activity, error) {
	var productLots []*ProductLot
	var activities []*Activity

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, nil, err
		}

		if t == IterationType(0) {
			var prod ProductLot
			err = json.Unmarshal(queryResult.Value, &prod)
			if err != nil {
				return nil, nil, err
			}
			productLots = append(productLots, &prod)
		}

		if t == IterationType(1) {
			var activity Activity
			err = json.Unmarshal(queryResult.Value, &activity)
			if err != nil {
				return nil, nil, err
			}
			activities = append(activities, &activity)
		}

	}

	return productLots, activities, nil
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string, t IterationType) ([]*ProductLot, []*Activity, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator, t)
}

func buildQueryString(key string, value string) string {
	return fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}", key, value)
}
