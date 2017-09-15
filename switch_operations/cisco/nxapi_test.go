package cisco_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"github.com/go-openapi/errors"
	"time"
	"github.com/RackHD/on-network/switch_operations/cisco"
)

type FakeRunner struct {
	FailCopyCommand bool
	FailInstallCommand bool
	FailReconnecting bool
	FailShowVersionCommand bool
	SuccessShowVersion bool
	ImageFileName string
}


func (fr *FakeRunner) Run(command string) (string, error) {
	if strings.Contains(command, "copy") && fr.FailCopyCommand {
		return "", errors.New(1, "fake copy command failed")
	}
	if strings.Contains(command, "install") && fr.FailInstallCommand{
		return "", errors.New(2, "fake install command failed")
	}
	if strings.Contains(command, "show") && fr.FailReconnecting{
		time.Sleep(20 * time.Second)
		return "", errors.New(2, "failed to reconnect")
	}
	if strings.Contains(command, "show") && fr.FailShowVersionCommand{
		return "fake.bin", nil
	}
	if strings.Contains(command, "copy") && fr.SuccessShowVersion {
		fr.ImageFileName = strings.Split((strings.Split(command, " ")[2]),":")[1]
	}
	if strings.Contains(command, "show") && fr.SuccessShowVersion{
		return fr.ImageFileName, nil
	}

	return "", nil
}

var _ = Describe("Nxapi", func() {
	Context("When copy command fails", func() {
		It("should return an error", func(){
 			ciscoSwitch := cisco.Switch{ &FakeRunner{FailCopyCommand:true}}
			err:= ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error copying image from remote"))
		})
	})

	Context("When install all command fails", func() {
		It("should return an error", func(){
			ciscoSwitch := cisco.Switch{ &FakeRunner{FailInstallCommand:true}}
			err:= ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error install image"))
		})
	})

	Context("When fails to reconnect to the switch ", func() {
		It("should return an error", func(){
			ciscoSwitch := cisco.Switch{ &FakeRunner{FailReconnecting:true}}
			err:= ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("timeout connecting to switch after upgrade"))
		})
	})

	Context("When show version command fails", func() {
		It("should return an error", func(){
			ciscoSwitch := cisco.Switch{ &FakeRunner{FailShowVersionCommand:true}}
			err:= ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to find the expected version"))
		})
	})

	Context("When upgrade is successful", func() {
		It("shouldnt return any error", func(){
			ciscoSwitch := cisco.Switch{ &FakeRunner{SuccessShowVersion:true}}
			err:= ciscoSwitch.Update("1.1.1.1/test.bin")
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
