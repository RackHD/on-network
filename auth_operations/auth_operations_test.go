package auth_operations_test

import (
	"os"

	"github.com/RackHD/on-network/auth_operations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuthOperations", func() {


	BeforeEach(func() {
		os.Setenv("SERVICE_USERNAME", "admin")
		os.Setenv("SERVICE_PASSWORD", "Password123!")
	})

	Describe("Validate Login", func() {
		Context("When login is invalid", func() {
			It("should return false", func() {
				claim := auth_operations.Claims{

				}

				value, err := claim.ValidateLogin("root", "password")
				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(Equal(false))

			})
		})

		Context("When login is valid", func() {
			It("should return true", func() {
				claim := auth_operations.Claims{

				}

				value, err := claim.ValidateLogin("admin", "Password123!")
				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(Equal(true))

			})
		})

		Context("When login is empty", func() {
			It("should return false", func() {
				claim := auth_operations.Claims{

				}

				os.Unsetenv("SERVICE_USERNAME")
				os.Unsetenv("SERVICE_PASSWORD")

				_, err := claim.ValidateLogin("", "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Service Username or Password not set"))

			})
		})
	})

	Describe("Validate Token", func() {
		Context("When Token is valid", func() {
			It("should return true", func() {
				claim := auth_operations.Claims{

				}
				signedToken := claim.SetToken("admin")
				validToken := auth_operations.ValidateToken("Bearer " +signedToken)
				Expect(validToken).To(Equal(true))
			})
		})

		Context("When Token is invalid", func() {
			It("should return false", func() {
				claim := auth_operations.Claims{

				}
				signedToken := claim.SetToken("admin")
				validToken := auth_operations.ValidateToken("Bearer " + signedToken +"xyz")
				Expect(validToken).To(Equal(false))
			})
		})
	})
})

