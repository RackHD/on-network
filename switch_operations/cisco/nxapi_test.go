package cisco_test

import (
	"github.com/RackHD/on-network/switch_operations/cisco"
	"github.com/RackHD/on-network/switch_operations/cisco/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nxapi", func() {
	Context("When copy command fails", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{&fake.FakeRunner{FailCopyCommand: true}}
			err := ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error copying image from remote"))
		})
	})

	Context("When install all command fails", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{&fake.FakeRunner{FailInstallCommand: true}}
			err := ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error install image"))
		})
	})

	Context("When fails to reconnect to the switch ", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{&fake.FakeRunner{FailReconnecting: true}}
			err := ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("timeout connecting to switch after update"))
		})
	})

	Context("When show version command fails", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{&fake.FakeRunner{FailShowVersionCommand: true}}
			err := ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to find the expected version"))
		})
	})

	Context("When update is successful", func() {
		It("shouldnt return any error", func() {
			ciscoSwitch := cisco.Switch{&fake.FakeRunner{SuccessShowVersion: true}}
			err := ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
