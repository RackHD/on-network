// Copyright 2017, Dell EMC, Inc.

package switch_firmware_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSwitchFirmware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SwitchFirmware Suite")
}
