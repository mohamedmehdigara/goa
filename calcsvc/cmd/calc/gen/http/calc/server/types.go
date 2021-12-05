// Code generated by goa v3.5.2, DO NOT EDIT.
//
// calc HTTP server types
//
// Command:
// $ goa gen calcsvc/design

package server

import (
	calc "calcsvc/cmd/calc/gen/calc"
)

// NewAddPayload builds a calc service add endpoint payload.
func NewAddPayload(a int, b int) *calc.AddPayload {
	v := &calc.AddPayload{}
	v.A = a
	v.B = b

	return v
}