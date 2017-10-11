package switch_version

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/RackHD/on-network/models"
	"github.com/RackHD/on-network/switch_operations"
	"github.com/RackHD/on-network/switch_operations/cisco"
	"github.com/RackHD/on-network/switch_operations/cisco/nexus"
)

// SwitchVersion is a struct for the http objects
type SwitchVersion struct {
	Request *http.Request
	Client  switch_operations.ISwitch
}

// MiddleWare handles the route call
func MiddleWare(r *http.Request, body *models.Switch) middleware.Responder {
	fmt.Println("postSwitchFirmware")

	var client switch_operations.ISwitch

	if *body.Endpoint.SwitchType == "cisco" {
		client = &cisco.Switch{
			Runner: &nexus.Runner{
				IP:       *body.Endpoint.IP,
				Username: *body.Endpoint.Username,
				Password: *body.Endpoint.Password,
			},
		}
	}

	return &SwitchVersion{
		Request: r,
		Client:  client,
	}
}

// WriteResponse implements the CRUD logic behind the /credentials route
func (c *SwitchVersion) WriteResponse(rw http.ResponseWriter, rp runtime.Producer) {
	switch c.Request.Method {
	case http.MethodPost:
		c.postSwitchVersion(rw, rp)
	default:
		c.notSupported(rw, rp)
	}
}

func (c *SwitchVersion) notSupported(rw http.ResponseWriter, rp runtime.Producer) {
	rw.WriteHeader(http.StatusNotImplemented)
}

func (c *SwitchVersion) postSwitchVersion(rw http.ResponseWriter, rp runtime.Producer) {
	version, err := c.Client.GetFullVersion()
	if err != nil {
		rp.Produce(rw, fmt.Sprintf("failed to fetch firmware version: %+v", err))
		return
	}

	if err := rp.Produce(rw, version); err != nil {
		panic(err)
	}
}
