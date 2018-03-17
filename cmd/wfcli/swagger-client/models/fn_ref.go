// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// FnRef FnRef is an immutable, unique reference to a function on a specific function runtime environment.
//
// The string representation (via String or Format): runtime://runtimeId
// swagger:model FnRef
type FnRef struct {

	// Runtime is the Function Runtime environment (fnenv) that was used to resolve the function.
	Runtime string `json:"runtime,omitempty"`

	// RuntimeId is the runtime-specific identifier of the function.
	RuntimeID string `json:"runtimeId,omitempty"`
}

// Validate validates this fn ref
func (m *FnRef) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *FnRef) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FnRef) UnmarshalBinary(b []byte) error {
	var res FnRef
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}