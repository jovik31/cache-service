// Package model provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package model

// Id Integer representing the worker's Id
type Id = int

// Request defines model for Request.
type Request struct {
	// Ids Array of integers representing the worker's Id
	Ids []Id `json:"ids"`

	// Timestamp Integer representing the current timestamp.
	// Your Service clock must be updated with this value.
	// Lower Timestamps do not update the Service clock.
	Timestamp Timestamp `json:"timestamp"`
}

// Response defines model for Response.
type Response = []ResponseItem

// ResponseItem Object with the worker's Id and Status
type ResponseItem struct {
	// Id Integer representing the worker's Id
	Id     Id     `json:"id"`
	Status Status `json:"status"`
}

// Status defines model for Status.
type Status = int

// Timestamp Integer representing the current timestamp.
// Your Service clock must be updated with this value.
// Lower Timestamps do not update the Service clock.
type Timestamp = int
