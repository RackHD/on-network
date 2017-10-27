// Code generated by go-swagger; DO NOT EDIT.

package update_switch

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/RackHD/on-network/models"
)

// UpdateSwitchCreatedCode is the HTTP code returned for type UpdateSwitchCreated
const UpdateSwitchCreatedCode int = 201

/*UpdateSwitchCreated Successfully issued update switch firmware

swagger:response updateSwitchCreated
*/
type UpdateSwitchCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Status `json:"body,omitempty"`
}

// NewUpdateSwitchCreated creates UpdateSwitchCreated with default headers values
func NewUpdateSwitchCreated() *UpdateSwitchCreated {
	return &UpdateSwitchCreated{}
}

// WithPayload adds the payload to the update switch created response
func (o *UpdateSwitchCreated) WithPayload(payload *models.Status) *UpdateSwitchCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update switch created response
func (o *UpdateSwitchCreated) SetPayload(payload *models.Status) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateSwitchCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*UpdateSwitchDefault Error

swagger:response updateSwitchDefault
*/
type UpdateSwitchDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewUpdateSwitchDefault creates UpdateSwitchDefault with default headers values
func NewUpdateSwitchDefault(code int) *UpdateSwitchDefault {
	if code <= 0 {
		code = 500
	}

	return &UpdateSwitchDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the update switch default response
func (o *UpdateSwitchDefault) WithStatusCode(code int) *UpdateSwitchDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the update switch default response
func (o *UpdateSwitchDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the update switch default response
func (o *UpdateSwitchDefault) WithPayload(payload *models.ErrorResponse) *UpdateSwitchDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update switch default response
func (o *UpdateSwitchDefault) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateSwitchDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
