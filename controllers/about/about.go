package about

import (
	"github.com/RackHD/on-network/models"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
)

// Info is a struct for the http objects
type About struct {
	Request *http.Request
}

// MiddleWare handles the route call
func MiddleWare(r *http.Request) middleware.Responder {
	return &About{
		Request: r,
	}
}

// WriteResponse implements the CRUD logic behind the /credentials route
func (c *About) WriteResponse(rw http.ResponseWriter, rp runtime.Producer) {
	switch c.Request.Method {
	case http.MethodGet:
		c.getAbout(rw, rp)
	default:
		c.notSupported(rw, rp)
	}
}

//
func (c *About) notSupported(rw http.ResponseWriter, rp runtime.Producer) {
	rw.WriteHeader(http.StatusNotImplemented)
}
func (c *About) getAbout(rw http.ResponseWriter, rp runtime.Producer) {
	name := "on-network"
	info := models.About{
		Name: &name,
	}
	if err := rp.Produce(rw, info); err != nil {
		panic(err)
	}
}
