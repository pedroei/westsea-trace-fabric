/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

type ProductLot struct {
	DocType           string        `json:"docType"`
	ID                string        `json:"ID"`
	ReferenceNumber   string        `json:"referenceNumber"`
	IsSerialNumber    bool          `json:"isSerialNumber"`
	Designation       string        `json:"designation"`
	ProductType       string        `json:"productType"`
	InitialQuantity   float32       `json:"initialQuantity"`
	AvailableQuantity float32       `json:"availableQuantity"`
	DocumentKeys      []DocumentKey `json:"documentKeys,omitempty" metadata:"documentKeys,optional"`
}

type Activity struct {
	DocType          string             `json:"docType"`
	ID               string             `json:"ID"`
	Designation      string             `json:"designation"`
	UserId           string             `json:"userId"`
	DateTime         string             `json:"dateTime"`
	InputProductLots map[string]float32 `json:"inputProductLots"` //input id and amount
	OutputProductLot ProductLot         `json:"outputProductLot"`
}

type DocumentKey struct {
	DocumentKey     string `json:"documentKey"`
	FileFingerPrint string `json:"fileFingerPrint"`
}
