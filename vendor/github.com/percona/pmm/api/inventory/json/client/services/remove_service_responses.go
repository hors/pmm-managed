// Code generated by go-swagger; DO NOT EDIT.

package services

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// RemoveServiceReader is a Reader for the RemoveService structure.
type RemoveServiceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RemoveServiceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewRemoveServiceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewRemoveServiceOK creates a RemoveServiceOK with default headers values
func NewRemoveServiceOK() *RemoveServiceOK {
	return &RemoveServiceOK{}
}

/*RemoveServiceOK handles this case with default header values.

(empty)
*/
type RemoveServiceOK struct {
	Payload interface{}
}

func (o *RemoveServiceOK) Error() string {
	return fmt.Sprintf("[POST /v1/inventory/Services/Remove][%d] removeServiceOK  %+v", 200, o.Payload)
}

func (o *RemoveServiceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*RemoveServiceBody remove service body
swagger:model RemoveServiceBody
*/
type RemoveServiceBody struct {

	// Unique randomly generated instance identifier.
	ServiceID string `json:"service_id,omitempty"`
}

// Validate validates this remove service body
func (o *RemoveServiceBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *RemoveServiceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RemoveServiceBody) UnmarshalBinary(b []byte) error {
	var res RemoveServiceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
