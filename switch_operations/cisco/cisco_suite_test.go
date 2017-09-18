package cisco_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCisco(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cisco Suite")
}
