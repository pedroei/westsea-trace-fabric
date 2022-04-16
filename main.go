/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	westseaTraceShipContract := new(WestseaTraceShipContract)
	westseaTraceShipContract.Info.Version = "0.0.1"
	westseaTraceShipContract.Info.Description = "My Smart Contract"
	westseaTraceShipContract.Info.License = new(metadata.LicenseMetadata)
	westseaTraceShipContract.Info.License.Name = "Apache-2.0"
	westseaTraceShipContract.Info.Contact = new(metadata.ContactMetadata)
	westseaTraceShipContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(westseaTraceShipContract)
	chaincode.Info.Title = "westsea-trace-fabric chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from WestseaTraceShipContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
