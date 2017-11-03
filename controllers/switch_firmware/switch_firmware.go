package switch_firmware

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

// SwitchConfig is a struct for the http objects
type SwitchFirmware struct {
	Request *http.Request
	Client  switch_operations.ISwitch
}

// MiddleWare handles the route call
func MiddleWare(r *http.Request, body *models.Switch) middleware.Responder {
	fmt.Println("postSwitchFirmware")

	var client switch_operations.ISwitch

	if *body.Endpoint.SwitchType == "cisco" || *body.Endpoint.SwitchType == "nexus" {
		client = &cisco.Switch{
			Runner: &nexus.Runner{
				IP:       *body.Endpoint.Ipaddress,
				Username: *body.Endpoint.Username,
				Password: *body.Endpoint.Password,
			},
		}
	}

	return &SwitchFirmware{
		Request: r,
		Client:  client,
	}
}

// WriteResponse implements the CRUD logic behind the /credentials route
func (c *SwitchFirmware) WriteResponse(rw http.ResponseWriter, rp runtime.Producer) {
	switch c.Request.Method {
	case http.MethodPost:
		c.postSwitchFirmware(rw, rp)
	default:
		c.notSupported(rw, rp)
	}
}

func (c *SwitchFirmware) notSupported(rw http.ResponseWriter, rp runtime.Producer) {
	rw.WriteHeader(http.StatusNotImplemented)
}

func (c *SwitchFirmware) postSwitchFirmware(rw http.ResponseWriter, rp runtime.Producer) {
	version, err := c.Client.GetFirmware()
	if err != nil {
		rw.WriteHeader(400)
		rp.Produce(rw, fmt.Sprintf("failed to fetch firmware version: %+v", err))
		return
	}

	versionResponse := models.SwitchVersionResponse{
		Version: version,
	}

	if err := rp.Produce(rw, versionResponse); err != nil {
		panic(err)
	}
}
