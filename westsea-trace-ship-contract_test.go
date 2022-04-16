/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

// func configureStub() (*MockContext, *MockStub) {
// 	var nilBytes []byte

// 	testWestseaTraceShip := new(WestseaTraceShip)
// 	testWestseaTraceShip.Value = "set value"
// 	westseaTraceShipBytes, _ := json.Marshal(testWestseaTraceShip)

// 	ms := new(MockStub)
// 	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
// 	ms.On("GetState", "missingkey").Return(nilBytes, nil)
// 	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
// 	ms.On("GetState", "westseaTraceShipkey").Return(westseaTraceShipBytes, nil)
// 	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
// 	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

// 	mc := new(MockContext)
// 	mc.On("GetStub").Return(ms)

// 	return mc, ms
// }

// func TestWestseaTraceShipExists(t *testing.T) {
// 	var exists bool
// 	var err error

// 	ctx, _ := configureStub()
// 	c := new(WestseaTraceShipContract)

// 	exists, err = c.WestseaTraceShipExists(ctx, "statebad")
// 	assert.EqualError(t, err, getStateError)
// 	assert.False(t, exists, "should return false on error")

// 	exists, err = c.WestseaTraceShipExists(ctx, "missingkey")
// 	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
// 	assert.False(t, exists, "should return false when no value for key in world state")

// 	exists, err = c.WestseaTraceShipExists(ctx, "existingkey")
// 	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
// 	assert.True(t, exists, "should return true when value for key in world state")
// }

// func TestCreateWestseaTraceShip(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(WestseaTraceShipContract)

// 	err = c.CreateWestseaTraceShip(ctx, "statebad", "some value")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

// 	err = c.CreateWestseaTraceShip(ctx, "existingkey", "some value")
// 	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

// 	err = c.CreateWestseaTraceShip(ctx, "missingkey", "some value")
// 	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
// }

// func TestReadWestseaTraceShip(t *testing.T) {
// 	var westseaTraceShip *WestseaTraceShip
// 	var err error

// 	ctx, _ := configureStub()
// 	c := new(WestseaTraceShipContract)

// 	westseaTraceShip, err = c.ReadWestseaTraceShip(ctx, "statebad")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
// 	assert.Nil(t, westseaTraceShip, "should not return WestseaTraceShip when exists errors when reading")

// 	westseaTraceShip, err = c.ReadWestseaTraceShip(ctx, "missingkey")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
// 	assert.Nil(t, westseaTraceShip, "should not return WestseaTraceShip when key does not exist in world state when reading")

// 	westseaTraceShip, err = c.ReadWestseaTraceShip(ctx, "existingkey")
// 	assert.EqualError(t, err, "Could not unmarshal world state data to type WestseaTraceShip", "should error when data in key is not WestseaTraceShip")
// 	assert.Nil(t, westseaTraceShip, "should not return WestseaTraceShip when data in key is not of type WestseaTraceShip")

// 	westseaTraceShip, err = c.ReadWestseaTraceShip(ctx, "westseaTraceShipkey")
// 	expectedWestseaTraceShip := new(WestseaTraceShip)
// 	expectedWestseaTraceShip.Value = "set value"
// 	assert.Nil(t, err, "should not return error when WestseaTraceShip exists in world state when reading")
// 	assert.Equal(t, expectedWestseaTraceShip, westseaTraceShip, "should return deserialized WestseaTraceShip from world state")
// }

// func TestUpdateWestseaTraceShip(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(WestseaTraceShipContract)

// 	err = c.UpdateWestseaTraceShip(ctx, "statebad", "new value")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

// 	err = c.UpdateWestseaTraceShip(ctx, "missingkey", "new value")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

// 	err = c.UpdateWestseaTraceShip(ctx, "westseaTraceShipkey", "new value")
// 	expectedWestseaTraceShip := new(WestseaTraceShip)
// 	expectedWestseaTraceShip.Value = "new value"
// 	expectedWestseaTraceShipBytes, _ := json.Marshal(expectedWestseaTraceShip)
// 	assert.Nil(t, err, "should not return error when WestseaTraceShip exists in world state when updating")
// 	stub.AssertCalled(t, "PutState", "westseaTraceShipkey", expectedWestseaTraceShipBytes)
// }

// func TestDeleteWestseaTraceShip(t *testing.T) {
// 	var err error

// 	ctx, stub := configureStub()
// 	c := new(WestseaTraceShipContract)

// 	err = c.DeleteWestseaTraceShip(ctx, "statebad")
// 	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

// 	err = c.DeleteWestseaTraceShip(ctx, "missingkey")
// 	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

// 	err = c.DeleteWestseaTraceShip(ctx, "westseaTraceShipkey")
// 	assert.Nil(t, err, "should not return error when WestseaTraceShip exists in world state when deleting")
// 	stub.AssertCalled(t, "DelState", "westseaTraceShipkey")
// }
