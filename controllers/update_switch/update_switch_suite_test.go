package update_switch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUpdateSwitch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UpdateSwitch Suite")
}
