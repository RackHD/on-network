package switch_version_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSwitchFirmware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Switch Version Suite")
}
