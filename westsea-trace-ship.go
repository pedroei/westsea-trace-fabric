/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

//FIXME: a product can only be on ONE activity?
type ProductSerialNum struct {
	DocType 	 	string   `json:"docType"`
	ID      	 	string   `json:"ID"`
	SerialNumber 	string   `json:"serialNumber"`
	Designation  	string   `json:"designation"`
	ProductType  	string   `json:"productType"`
	DocumentKeys 	[]string `json:"documentKeys,omitempty" metadata:"documentKeys,optional"`
	// BuilderActivity string 	 `json:"builderActivity"` //activity that originated the product
}

//FIXME: the outputs can be more than one ??
type Activity struct {
	DocType 	 			string   			`json:"docType"`
	ID      	 			string   			`json:"ID"`
	Designation  			string   			`json:"designation"`
	UserId  				string   			`json:"userId"` //user that performed the activity
	DateTime 				string 	 			`json:"dateTime"`
	InputProductSerialNums 	[]string 			`json:"inputProductSerialNums"`
	OutputProductSerialNum 	ProductSerialNum   	`json:"outputProductSerialNum"` //FIXME: should this be just the ID?
}
