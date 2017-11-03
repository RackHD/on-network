// Copyright 2017, Dell EMC, Inc.

package switch_config_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/RackHD/on-network/controllers/switch_config"
	"github.com/RackHD/on-network/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestProducer struct{}

// Produce is ...
func (t TestProducer) Produce(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

var _ = Describe("SwitchConfig", func() {
	var prod TestProducer
	var buff *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up receiver to mock out where response would go
		prod = TestProducer{}
		buff = httptest.NewRecorder()

		os.Setenv("SWITCH_MODELS_FILE_PATH", "../../switch_operations/cisco/fake/switchModels.yml")
	})

	Context("when a message is routed to the /switchConfig handler", func() {
		It("info API should return siwtch running config", func() {
			// Create on-network api about
			serverURL := "http://localhost:8080"

			jsonBody := []byte(`{
				"endpoint": {
					"ip": "test",
					"username": "test",
					"password": "test",
					"switchType": "cisco"
				},
				"imageURL": "test",
				"switchModel": "Nexus3000 C3164PQ Chassis"
			}`)

			req, err := http.NewRequest("POST", serverURL+"/switchConfig", bytes.NewBuffer(jsonBody))
			Expect(err).ToNot(HaveOccurred())

			switchConfig := &models.Switch{}
			err = json.Unmarshal(jsonBody, switchConfig)
			Expect(err).ToNot(HaveOccurred())

			responder := MiddleWare(req, switchConfig)
			responder.WriteResponse(buff, prod)
			Expect(buff.Code).To(Equal(http.StatusNotFound))
		})
	})
})
