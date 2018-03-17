// Code generated by go-swagger; DO NOT EDIT.

package workflow_invocation_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/fission/fission-workflows/cmd/wfcli/swagger-client/models"
)

// InvokeSync2Reader is a Reader for the InvokeSync2 structure.
type InvokeSync2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *InvokeSync2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewInvokeSync2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewInvokeSync2OK creates a InvokeSync2OK with default headers values
func NewInvokeSync2OK() *InvokeSync2OK {
	return &InvokeSync2OK{}
}

/*InvokeSync2OK handles this case with default header values.

InvokeSync2OK invoke sync2 o k
*/
type InvokeSync2OK struct {
	Payload *models.WorkflowInvocation
}

func (o *InvokeSync2OK) Error() string {
	return fmt.Sprintf("[GET /invocation/sync][%d] invokeSync2OK  %+v", 200, o.Payload)
}

func (o *InvokeSync2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.WorkflowInvocation)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}