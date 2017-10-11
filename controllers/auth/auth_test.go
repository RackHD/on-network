package auth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/RackHD/on-network/controllers/auth"
	"github.com/RackHD/on-network/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

type TestProducer struct{}

// Produce is ...
func (t TestProducer) Produce(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

var _ = Describe("loginFunction", func() {
	var prod TestProducer
	var buff *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up receiver to mock out where response would go
		prod = TestProducer{}
		buff = httptest.NewRecorder()
		fmt.Println("Inside login test")

		os.Setenv("SERVICE_USERNAME", "admin")
		os.Setenv("SERVICE_PASSWORD", "Password123!")
	})

	Context("when a message is routed to the /login handler", func() {
		It("info API should return a generated token", func() {

			serverURL := "http://localhost:8080"

			jsonBody := []byte(`{
					"username": "test",
					"password": "test"
			}`)

			req, err := http.NewRequest("POST", serverURL+"/login", bytes.NewBuffer(jsonBody))
			Expect(err).ToNot(HaveOccurred())

			login := &models.Login{}
			err = json.Unmarshal(jsonBody, login)
			Expect(err).ToNot(HaveOccurred())

			responder := MiddleWare(req, login)
			responder.WriteResponse(buff, prod)
			Expect(buff.Code).To(Equal(http.StatusUnauthorized))

		})
	})
})
