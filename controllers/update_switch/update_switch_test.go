// Copyright 2017, Dell EMC, Inc.

package update_switch_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/RackHD/on-network/controllers/update_switch"
	"github.com/RackHD/on-network/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

type TestProducer struct{}

//Produce is ...
func (t TestProducer) Produce(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

var _ = Describe("UpdateSwitch", func() {
	var prod TestProducer
	var buff *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up receiver to mock out where response would go
		prod = TestProducer{}
		buff = httptest.NewRecorder()

		os.Setenv("SWITCH_MODELS_FILE_PATH", "../../switch_operations/cisco/fake/switchModels.yml")
	})

	Context("When a message is routed to the /updateSwitch handler", func() {
		It("info API should return 'status OK'", func() {
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

			req, err := http.NewRequest("POST", serverURL+"/updateSwitch", bytes.NewBuffer(jsonBody))
			Expect(err).ToNot(HaveOccurred())

			updateSwitch := &models.UpdateSwitch{}
			err = json.Unmarshal(jsonBody, updateSwitch)
			Expect(err).ToNot(HaveOccurred())

			responder := MiddleWare(req, updateSwitch)
			responder.WriteResponse(buff, prod)
			Expect(buff.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
