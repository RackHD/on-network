package switch_config

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
type SwitchConfig struct {
	Request *http.Request
	Client  switch_operations.ISwitch
}

// MiddleWare handles the route call
func MiddleWare(r *http.Request, body *models.Switch) middleware.Responder {
	fmt.Println("postSwitchConfig")

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

	return &SwitchConfig{
		Request: r,
		Client:  client,
	}
}

// WriteResponse implements the CRUD logic behind the /credentials route
func (c *SwitchConfig) WriteResponse(rw http.ResponseWriter, rp runtime.Producer) {
	switch c.Request.Method {
	case http.MethodPost:
		c.postSwitchConfig(rw, rp)
	default:
		c.notSupported(rw, rp)
	}
}

func (c *SwitchConfig) notSupported(rw http.ResponseWriter, rp runtime.Producer) {
	rw.WriteHeader(http.StatusNotImplemented)
}

func (c *SwitchConfig) postSwitchConfig(rw http.ResponseWriter, rp runtime.Producer) {
	config, err := c.Client.GetConfig()
	if err != nil {
		rw.WriteHeader(404)
		rp.Produce(rw, fmt.Sprintf("failed to update switch: %+v", err))
		return
	}

	configResponse := models.SwitchConfigResponse{
		Config: config,
	}

	if err := rp.Produce(rw, configResponse); err != nil {
		panic(err)
	}
}
