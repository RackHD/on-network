// Copyright 2017, Dell EMC, Inc.

package cisco_test

import (
	"os"

	"github.com/RackHD/on-network/switch_operations/cisco"
	"github.com/RackHD/on-network/switch_operations/cisco/fake"

	"github.com/RackHD/on-network/models"
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

	Describe("UpdateFirmware", func() {
		Context("When copy command fails for 6.0 firmware", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailCopyCommand: true},
				}

				var firmwareImages []*models.FirmwareImage

				imageTypeKickstart := "kickstart"
				imageURLKickstart := "1.1.1.1/kickstart.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeKickstart, ImageURL: &imageURLKickstart})

				imageTypeSystem := "system"
				imageURLSystem := "1.1.1.1/system.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeSystem, ImageURL: &imageURLSystem})

				err := ciscoSwitch.Update(disruptiveSwitchModel, firmwareImages)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error copying image from remote"))
			})
		})

		Context("When copy command fails for 7.0 firmware ", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailCopyCommand: true},
				}

				var firmwareImages []*models.FirmwareImage

				imageType := "system"
				imageURL := "1.1.1.1/test.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageType, ImageURL: &imageURL})
				err := ciscoSwitch.Update(nonDisruptiveSwitchModel, firmwareImages)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error copying image from remote"))
			})
		})

		Context("When passing wrong parameters to copy for 6.0 ", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailCopyCommandWrongParam: true},
				}

				var firmwareImages []*models.FirmwareImage

				imageType := "test"
				imageURL := "1.1.1.1/test.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageType, ImageURL: &imageURL})

				err := ciscoSwitch.Update(disruptiveSwitchModel, firmwareImages)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Missing required image type"))
			})
		})

		Context("When install all command fails", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailInstallCommand: true},
				}

				var firmwareImages []*models.FirmwareImage

				imageTypeKickstart := "kickstart"
				imageURLKickstart := "1.1.1.1/kickstart.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeKickstart, ImageURL: &imageURLKickstart})

				imageTypeSystem := "system"
				imageURLSystem := "1.1.1.1/system.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeSystem, ImageURL: &imageURLSystem})

				err := ciscoSwitch.Update(disruptiveSwitchModel, firmwareImages)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Installation failed"))
			})
		})

		Context("When fails to reconnect to the switch ", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailReconnecting: true},
				}
				var firmwareImages []*models.FirmwareImage

				imageTypeKickstart := "kickstart"
				imageURLKickstart := "1.1.1.1/kickstart.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeKickstart, ImageURL: &imageURLKickstart})

				imageTypeSystem := "system"
				imageURLSystem := "1.1.1.1/system.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeSystem, ImageURL: &imageURLSystem})
				err := ciscoSwitch.Update(disruptiveSwitchModel, firmwareImages)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("connecting to the switch after update or failed to find the right version"))
			})
		})

		Context("When show version command fails", func() {
			It("should return an error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailShowVersionCommand: true},
				}

				var firmwareImages []*models.FirmwareImage

				imageTypeKickstart := "kickstart"
				imageURLKickstart := "1.1.1.1/kickstart.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeKickstart, ImageURL: &imageURLKickstart})

				imageTypeSystem := "system"
				imageURLSystem := "1.1.1.1/system.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeSystem, ImageURL: &imageURLSystem})
				err := ciscoSwitch.Update(disruptiveSwitchModel, firmwareImages)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("connecting to the switch after update or failed to find the right version"))
			})
		})

		Context("When update is successful for ", func() {
			It("shouldnt return any error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{SuccessShowVersion: true},
				}

				var firmwareImages []*models.FirmwareImage

				imageTypeKickstart := "kickstart"
				imageURLKickstart := "1.1.1.1/kickstart.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeKickstart, ImageURL: &imageURLKickstart})

				imageTypeSystem := "system"
				imageURLSystem := "1.1.1.1/system.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageTypeSystem, ImageURL: &imageURLSystem})
				err := ciscoSwitch.Update(disruptiveSwitchModel, firmwareImages)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("When the non-disruptive downgrade is not supported", func() {
			It("shouldn't return any error ", func() {
				fakeRunner := &fake.FakeRunner{DowngradeNonDisruptive: true, TimeoutInstall: true}
				ciscoSwitch := cisco.Switch{
					Runner: fakeRunner,
				}

				var firmwareImages []*models.FirmwareImage
				imageType := "system"
				imageURL := "1.1.1.1/test.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageType, ImageURL: &imageURL})
				err := ciscoSwitch.Update(nonDisruptiveSwitchModel, firmwareImages)
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeRunner.InstallCommand).To(ContainSubstring("non-disruptive"))
			})
		})

		Context("When the update type is non-disruptive", func() {
			It("should run install all with non-disruptive", func() {
				fakeRunner := &fake.FakeRunner{TimeoutInstall: true}
				ciscoSwitch := cisco.Switch{
					Runner: fakeRunner,
				}

				var firmwareImages []*models.FirmwareImage
				imageType := "system"
				imageURL := "1.1.1.1/test.bin"

				firmwareImages = append(firmwareImages, &models.FirmwareImage{ImageType: &imageType, ImageURL: &imageURL})
				err := ciscoSwitch.Update(nonDisruptiveSwitchModel, firmwareImages)
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

	Describe("CheckVlan", func() {
		Context("when check vlan command failed", func() {
			It("should return error", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{FailCheckVlanCommand: true},
				}

				_, err := ciscoSwitch.CheckVlan(0)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error running show vlan command"))
			})
		})

		Context("when show vlan command succeeded with null result", func() {
			It("should return false", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{SuccessShowVlanNullResult: true},
				}

				isVlan, err := ciscoSwitch.CheckVlan(5)
				Expect(err).ToNot(HaveOccurred())
				Expect(isVlan).To(Equal(false))
			})
		})

		Context("when show vlan command succeeded with result", func() {
			It("should return true", func() {
				ciscoSwitch := cisco.Switch{
					Runner: &fake.FakeRunner{SuccessShowVlanNullResult: false},
				}

				isVlan, err := ciscoSwitch.CheckVlan(5)
				Expect(err).ToNot(HaveOccurred())
				Expect(isVlan).To(Equal(true))
			})
		})
	})
})
