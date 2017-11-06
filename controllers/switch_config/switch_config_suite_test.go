// Copyright 2017, Dell EMC, Inc.

package switch_config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSwitchConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SwitchConfig Suite")
}
