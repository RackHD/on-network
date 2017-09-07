package about_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAbout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "About Suite")
}
