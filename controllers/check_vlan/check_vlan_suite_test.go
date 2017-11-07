// Copyright 2017, Dell EMC, Inc.

package check_vlan_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCheckVlan(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Check Vlan Suite")
}
