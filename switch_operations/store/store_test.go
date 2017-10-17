package store_test

import (
	"os"

	. "github.com/RackHD/on-network/switch_operations/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store", func() {
	var switchFileDatabase *SwitchFileDatabase

	BeforeEach(func() {
		os.Setenv("SWITCH_MODELS_FILE_PATH", "../cisco/fake/switchModels.yml")
		path := os.Getenv("SWITCH_MODELS_FILE_PATH")
		Expect(path).ToNot(BeEmpty())
		switchFileDatabase = GetSwitchFileDatabase()
	})

	Describe("GetUpdateType", func() {
		Context("If the model exist", func() {
			It("returns the disruptive type of the model", func() {
				dType, dFirmware, err := switchFileDatabase.GetUpdateType("cisco", "3132")
				Expect(err).ToNot(HaveOccurred())
				Expect(dType).To(Equal("Disruptive"))
				Expect(dFirmware).To(Equal("6.0"))

				dType, dFirmware, err = switchFileDatabase.GetUpdateType("cisco", "3164")
				Expect(err).ToNot(HaveOccurred())
				Expect(dType).To(Equal("NonDisruptive"))
				Expect(dFirmware).To(Equal("7.0"))
			})
		})

		Context("If the model not exist", func() {
			It("return the an error", func() {
				_, _, err := switchFileDatabase.GetUpdateType("verizon", "3132")
				Expect(err).To(MatchError("couldn't find switch type"))
				_, _, err = switchFileDatabase.GetUpdateType("cisco", "7979")
				Expect(err).To(MatchError("couldn't find switch model"))
			})
		})
	})
})
