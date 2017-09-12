package update_switch

import (
	"net/http"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)
// Info is a struct for the http objects
type UpdateSwitch struct {
	Request *http.Request
}
// MiddleWare handles the route call
func MiddleWare(r *http.Request) middleware.Responder {
	return &UpdateSwitch{
		Request: r,
	}
}
// WriteResponse implements the CRUD logic behind the /credentials route
func (c *UpdateSwitch) WriteResponse(rw http.ResponseWriter, rp runtime.Producer) {
	switch c.Request.Method {
	case http.MethodPost:
		c.postUpdateSwitch(rw, rp)
	default:
		c.notSupported(rw, rp)
	}
}
func (c *UpdateSwitch) notSupported(rw http.ResponseWriter, rp runtime.Producer) {
	rw.WriteHeader(http.StatusNotImplemented)
}
func (c *UpdateSwitch) postUpdateSwitch(rw http.ResponseWriter, rp runtime.Producer) {
	if err := rp.Produce(rw, "Hello Masarah!"); err != nil {
		panic(err)
	}
}