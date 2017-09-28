// Code generated by go-swagger; DO NOT EDIT.

package switch_config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/RackHD/on-network/models"
)

// SwitchConfigOKCode is the HTTP code returned for type SwitchConfigOK
const SwitchConfigOKCode int = 200

/*SwitchConfigOK Successfully returned switch running config

swagger:response switchConfigOK
*/
type SwitchConfigOK struct {
}

// NewSwitchConfigOK creates SwitchConfigOK with default headers values
func NewSwitchConfigOK() *SwitchConfigOK {
	return &SwitchConfigOK{}
}

// WriteResponse to the client
func (o *SwitchConfigOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
}

/*SwitchConfigDefault Error

swagger:response switchConfigDefault
*/
type SwitchConfigDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload models.ErrorResponse `json:"body,omitempty"`
}

// NewSwitchConfigDefault creates SwitchConfigDefault with default headers values
func NewSwitchConfigDefault(code int) *SwitchConfigDefault {
	if code <= 0 {
		code = 500
	}

	return &SwitchConfigDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the switch config default response
func (o *SwitchConfigDefault) WithStatusCode(code int) *SwitchConfigDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the switch config default response
func (o *SwitchConfigDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the switch config default response
func (o *SwitchConfigDefault) WithPayload(payload models.ErrorResponse) *SwitchConfigDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the switch config default response
func (o *SwitchConfigDefault) SetPayload(payload models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SwitchConfigDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}