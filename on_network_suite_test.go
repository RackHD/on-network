package on_network_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOnNetwork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OnNetwork Suite")
}
