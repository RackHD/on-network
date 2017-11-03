// Copyright 2017, Dell EMC, Inc.

package auth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSwitchConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Login Function Suite")
}
