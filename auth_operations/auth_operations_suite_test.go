// Copyright 2017, Dell EMC, Inc.

package auth_operations_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAuthOperations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthOperations Suite")
}
