package switch_version_test

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
	"github.com/RackHD/on-network/controllers/switch_firmware"
)

type TestProducer struct{}

// Produce is ...
func (t TestProducer) Produce(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

var _ = Describe("SwitchVersion", func() {
	var prod TestProducer
	var buff *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up receiver to mock out where response would go
		prod = TestProducer{}
		buff = httptest.NewRecorder()

		os.Setenv("SWITCH_MODELS_FILE_PATH", "../../switch_operations/cisco/fake/switchModels.yml")
	})

	Context("when a message is routed to the /switchVersion handler", func() {
		It("info API should return switch firmware", func() {
			// Create on-network api about
			serverURL := "http://localhost:8080"

			jsonBody := []byte(`{
				"endpoint": {
					"ip": "test",
					"username": "test",
					"password": "test",
					"switchType": "cisco"
				}
			}`)

			req, err := http.NewRequest("POST", serverURL+"/switchVersion", bytes.NewBuffer(jsonBody))
			Expect(err).ToNot(HaveOccurred())

			switchFirmware := &models.Switch{}
			err = json.Unmarshal(jsonBody, switchFirmware)

			Expect(err).ToNot(HaveOccurred())

			responder := switch_firmware.MiddleWare(req, switchFirmware)
			responder.WriteResponse(buff, prod)
			Expect(buff.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
