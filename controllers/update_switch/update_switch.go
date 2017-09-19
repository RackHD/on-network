package update_switch

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

// UpdateSwitch is a struct for the http objects
type UpdateSwitch struct {
	Request    *http.Request
	Client     switch_operations.Switch
	UpdateType string
	ImageURL   string
}

// MiddleWare handles the route call
func MiddleWare(r *http.Request, body *models.UpdateSwitch) middleware.Responder {
	var client switch_operations.Switch

	if *body.SwitchType == "cisco" {
		client = &cisco.Switch{
			Runner: &nexus.Runner{
				IP:       *body.IP,
				Username: *body.Username,
				Password: *body.Password,
			},
		}
	}

	return &UpdateSwitch{
		Request:    r,
		Client:     client,
		UpdateType: *body.UpdateType,
		ImageURL:   *body.ImageURL,
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
	err := c.Client.Update(c.UpdateType, c.ImageURL)
	if err != nil {
		rp.Produce(rw, fmt.Sprintf("failed to update switch: %+v", err))
		return
	}

	if err := rp.Produce(rw, fmt.Sprintf("succeeded to update switch!!!")); err != nil {
		panic(err)
	}
}
