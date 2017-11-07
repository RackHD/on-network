// Copyright 2017, Dell EMC, Inc.

package check_vlan

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

// CheckVlan is a struct for the http objects
type CheckVlan struct {
	Request *http.Request
	Client  switch_operations.ISwitch
	VlanID  int64
}

// MiddleWare handles the route call
func MiddleWare(r *http.Request, body *models.CheckVlan) middleware.Responder {
	fmt.Println("postCheckVlan")

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

	return &CheckVlan{
		Request: r,
		Client:  client,
		VlanID: *body.VlanID,

	}
}

// WriteResponse implements the CRUD logic behind the /credentials route
func (c *CheckVlan) WriteResponse(rw http.ResponseWriter, rp runtime.Producer) {
	switch c.Request.Method {
	case http.MethodPost:
		c.postCheckVlan(rw, rp)
	default:
		c.notSupported(rw, rp)
	}
}

func (c *CheckVlan) notSupported(rw http.ResponseWriter, rp runtime.Producer) {
	rw.WriteHeader(http.StatusNotImplemented)
}

func (c *CheckVlan) postCheckVlan(rw http.ResponseWriter, rp runtime.Producer) {
	isExist, err := c.Client.CheckVlan(c.VlanID)
	if err != nil {
		rw.WriteHeader(400)
		rp.Produce(rw, fmt.Sprintf("failed to fetch checkVlan: %+v", err))
		return
	}

	if err := rp.Produce(rw, isExist); err != nil {
		panic(err)
	}
}
