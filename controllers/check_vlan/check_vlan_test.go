// Copyright 2017, Dell EMC, Inc.

package check_vlan_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/RackHD/on-network/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/RackHD/on-network/controllers/check_vlan"
)

type TestProducer struct{}

// Produce is ...
func (t TestProducer) Produce(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

var _ = Describe("CheckVlan", func() {
	var prod TestProducer
	var buff *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up receiver to mock out where response would go
		prod = TestProducer{}
		buff = httptest.NewRecorder()

		os.Setenv("SWITCH_MODELS_FILE_PATH", "../../switch_operations/cisco/fake/switchModels.yml")
	})

	Context("when a message is routed to the /checkVlan handler", func() {
		It("info API should return whether the vlan exists", func() {
			// Create on-network api about
			serverURL := "http://localhost:8080"

			jsonBody := []byte(`{
				"endpoint": {
					"ipaddress": "test",
					"username": "test",
					"password": "test",
					"switchType": "cisco"
				},
				"vlanID": 1
			}`)

			req, err := http.NewRequest("POST", serverURL+"/checkVlan", bytes.NewBuffer(jsonBody))
			Expect(err).ToNot(HaveOccurred())

			checkVlan := &models.CheckVlan{}
			err = json.Unmarshal(jsonBody, checkVlan)

			Expect(err).ToNot(HaveOccurred())

			responder := check_vlan.MiddleWare(req, checkVlan)
			responder.WriteResponse(buff, prod)
			Expect(buff.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
