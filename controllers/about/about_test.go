package about_test

import (
	. "github.com/RackHD/on-network/controllers/about"

	"encoding/json"
	"github.com/RackHD/on-network/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
)

// TestProducer is ...
type TestProducer struct{}

//Produce is ...
func (t TestProducer) Produce(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

var _ = Describe("About", func() {
	var prod TestProducer
	var buff *httptest.ResponseRecorder
	BeforeEach(func() {
		// Set up receiver to mock out where response would go
		prod = TestProducer{}
		buff = httptest.NewRecorder()
	})
	Context("When a message is routed to the /api/about handler", func() {
		It("INTEGRATION info API should return an About object containing name='on-network'", func() {
			// Create on-network api about
			serverURL := "http://localhost:8080"
			req, err := http.NewRequest("GET", serverURL+"/about", nil)

			Expect(err).ToNot(HaveOccurred())
			// Put HTTP Request into router
			responder := MiddleWare(req)
			responder.WriteResponse(buff, prod)

			Expect(buff.Code).To(Equal(http.StatusOK))
			about := models.About{}
			err = json.Unmarshal(buff.Body.Bytes(), &about)
			Expect(err).ToNot(HaveOccurred())

			Expect(*about.Name).To(Equal("on-network"))
		})
	})
})
