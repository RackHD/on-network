package cisco_test

import (
	"os"

	"github.com/RackHD/on-network/switch_operations/cisco"
	"github.com/RackHD/on-network/switch_operations/cisco/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cisco", func() {
	var defaultUpdateSwitchType = "non-interruptive"

	BeforeEach(func() {
		os.Setenv("CISCO_RECONNECTION_TIMEOUT_IN_SECONDS", "8")
		os.Setenv("CISCO_BOOT_TIME_IN_SECONDS", "0")
	})

	Context("When copy command fails", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{
				Runner: &fake.FakeRunner{FailCopyCommand: true},
			}

			err := ciscoSwitch.Update(defaultUpdateSwitchType, "1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error copying image from remote"))
		})
	})

	Context("When install all command fails", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{
				Runner: &fake.FakeRunner{FailInstallCommand: true},
			}

			err := ciscoSwitch.Update(defaultUpdateSwitchType, "1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error install image"))
		})
	})

	Context("When fails to reconnect to the switch ", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{
				Runner: &fake.FakeRunner{FailReconnecting: true},
			}

			err := ciscoSwitch.Update(defaultUpdateSwitchType, "1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("timeout connecting to switch after update"))
		})
	})

	Context("When show version command fails", func() {
		It("should return an error", func() {
			ciscoSwitch := cisco.Switch{
				Runner: &fake.FakeRunner{FailShowVersionCommand: true},
			}

			err := ciscoSwitch.Update(defaultUpdateSwitchType, "1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to find the expected version"))
		})
	})

	Context("When update is successful", func() {
		It("shouldnt return any error", func() {
			ciscoSwitch := cisco.Switch{
				Runner: &fake.FakeRunner{SuccessShowVersion: true},
			}

			err := ciscoSwitch.Update(defaultUpdateSwitchType, "1.1.1.1/test.bin")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when the update type is interruptive", func() {
		It("should run install all with interruptive", func() {
			fakeRunner := &fake.FakeRunner{SuccessShowVersion: true}
			ciscoSwitch := cisco.Switch{
				Runner: fakeRunner,
			}

			err := ciscoSwitch.Update("interruptive", "1.1.1.1/test.bin")
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeRunner.InstallCommand).To(ContainSubstring(" interruptive"))
		})
	})
})
