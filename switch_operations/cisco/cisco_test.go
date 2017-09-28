package cisco_test

import (
	"os"

	"github.com/RackHD/on-network/switch_operations/cisco"
	"github.com/RackHD/on-network/switch_operations/cisco/fake"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cisco", func() {
	var disruptiveSwitchModel = "Nexus3000 C3132QX Chassis"

	var nonDisruptiveSwitchModel = "Nexus3000 C3164PQ Chassis"

	BeforeEach(func() {
		os.Setenv("SWITCH_MODELS_FILE_PATH", "fake/switchModels.yml")
		os.Setenv("CISCO_RECONNECTION_TIMEOUT_IN_SECONDS", "12")
		os.Setenv("CISCO_BOOT_TIME_IN_SECONDS", "0")
		os.Setenv("CISCO_INSTALL_TIME_IN_MINUTES", "1")
	})

	Describe("Update", func() {
		Context("When copy command fails", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailCopyCommand: true},
				}

				err := ciscoSwitch.Update(disruptiveSwitchModel, "1.1.1.1/test.bin")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error copying image from remote"))
			})
		})

		Context("When install all command fails", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailInstallCommand: true},
				}

				err := ciscoSwitch.Update(disruptiveSwitchModel, "1.1.1.1/test.bin")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error install image"))
			})
		})

		Context("When fails to reconnect to the switch ", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailReconnecting: true},
				}

				err := ciscoSwitch.Update(disruptiveSwitchModel, "1.1.1.1/test.bin")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("connecting to the switch after update or failed to find the right version"))
			})
		})

		Context("When show version command fails", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailShowVersionCommand: true},
				}

				err := ciscoSwitch.Update(disruptiveSwitchModel, "1.1.1.1/test.bin")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("connecting to the switch after update or failed to find the right version"))
			})
		})

		Context("When update is successful", func() {
			It("shouldnt return any error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{SuccessShowVersion: true},
				}

				err := ciscoSwitch.Update(disruptiveSwitchModel, "1.1.1.1/test.bin")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the non-disruptive downgrade is not support", func() {
			It("shouldnt return any error ", func() {
				fakeRunner := &fake.FakeRunner{DowngradeNonDisruptive: true, TimeoutInstall: true}
				ciscoSwitch := cisco.Switch{
					Runner: fakeRunner,
				}

				err := ciscoSwitch.Update(nonDisruptiveSwitchModel, "1.1.1.1/test.bin")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeRunner.InstallCommand).To(ContainSubstring("non-disruptive"))
			})
		})

		Context("when the update type is non-disruptive", func() {
			It("should run install all with non-disruptive", func() {
				fakeRunner := &fake.FakeRunner{TimeoutInstall: true}
				ciscoSwitch := cisco.Switch{
					Runner: fakeRunner,
				}

				err := ciscoSwitch.Update(nonDisruptiveSwitchModel, "1.1.1.1/test.bin")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeRunner.InstallCommand).To(ContainSubstring("non-disruptive"))
			})
		})
	})


	Describe("GetConfig", func() {
		Context("when get running config command failed", func() {
			It("should return error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailShowConfigCommand: true},
				}

				_, err := ciscoSwitch.GetConfig()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error running show running-config command"))
			})
		})

		Context("when get running config command succeeded", func() {
			It("should return config", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{},
				}

				config, err := ciscoSwitch.GetConfig()
				Expect(err).ToNot(HaveOccurred())
				Expect(config).To(Equal("{\"config\":\"empty\"}"))
			})
		})
	})

	Describe("GetFirmware", func() {
		Context("when get firmware version command failed", func() {
			It("should return error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailShowFirmwareVersionCommand: true},
				}

				_, err := ciscoSwitch.GetFirmware()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error running show version command"))
			})
		})

		Context("when get firmware version command succeeded", func() {
			It("should return firmware version", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{},
				}

				firmware, err := ciscoSwitch.GetFirmware()
				Expect(err).ToNot(HaveOccurred())
				Expect(firmware).To(Equal("7.0(3)I5(2)"))
			})
		})
	})
})
