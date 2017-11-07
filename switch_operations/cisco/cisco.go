package cisco

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"github.com/RackHD/on-network/switch_operations/cisco/nexus"
	"github.com/RackHD/on-network/switch_operations/store"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"

	"github.com/RackHD/on-network/models"
)

type Switch struct {
	Runner nexus.CommandRunner
}

func (c *Switch) Update(switchModel string, firmwareImages []*models.FirmwareImage) error {
	switchesDatabase := store.GetSwitchFileDatabase() //switch this if use database

	updateType, firmware, err := switchesDatabase.GetUpdateType("cisco", switchModel)
	if err != nil {
		return err
	}

	if firmware == "7.0" {
		imageURL := ""
		for _, firmwareImage := range firmwareImages {

			if *firmwareImage.ImageType == "system" {
				imageURL = *firmwareImage.ImageURL
			}
		}

		if imageURL == "" {
			return fmt.Errorf("Missing required image type: system")
		}

		imageFileName, err := c.copyImage(imageURL)
		if err != nil {
			return err
		}

		if updateType == "Disruptive" {
			err = c.disruptiveInstall(imageFileName)
		} else if updateType == "NonDisruptive" {
			err = c.nonDisruptiveInstall(imageFileName)
		}

		if err != nil {
			return fmt.Errorf("Installation failed: %+v", err)
		}
		return nil

	} else {
		kickstartURL := ""
		systemURL := ""

		for _, firmwareImage := range firmwareImages {

			if *firmwareImage.ImageType == "kickstart" {
				kickstartURL = *firmwareImage.ImageURL
			} else if *firmwareImage.ImageType == "system" {
				systemURL = *firmwareImage.ImageURL
			}
		}

		if kickstartURL == "" || systemURL == "" {
			return fmt.Errorf("Missing required image type: kickstart or system")
		}

		kickstartFileName, err := c.copyImage(kickstartURL)
		if err != nil {
			return err
		}

		systemFileName, err := c.copyImage(systemURL)
		if err != nil {
			return err
		}

		if updateType == "Disruptive" {
			err = c.disruptiveInstall6(kickstartFileName, systemFileName)
		}
		if err != nil {
			return fmt.Errorf("Installation failed: %+v", err)
		}
		return nil
	}

	return nil
}

func (c *Switch) copyImage(fileURL string) (string, error) {
	fileName := fmt.Sprintf("%s-%s", uuid.New().String(), path.Base(fileURL))
	fmt.Println("filename", fileName)

	copyCmd := fmt.Sprintf("copy %s bootflash:%s vrf management", fileURL, fileName)
	fmt.Println("starting copy")
	_, err := c.Runner.Run(copyCmd, "cli", 0)
	if err != nil {
		return "", fmt.Errorf("error copying image from remote: %+v", err)
	}

	return fileName, nil
}

func (c *Switch) disruptiveInstall6(kickstartFileName, systemFileName string) error {
	installCmd := fmt.Sprintf("install all kickstart bootflash:%s system  bootflash:%s non-interruptive", kickstartFileName, systemFileName)
	fmt.Println("starting disruptive installation on 6.0")

	_, err := c.Runner.Run(installCmd, "cli", 0)
	err = checkUpgradeErrorStatus(err)
	if err != nil {
		return err
	}

	err = c.checkNewVersion(kickstartFileName, true)
	if err != nil {
		return err
	}

	err = c.checkNewVersion(systemFileName, false)
	if err != nil {
		return err
	}

	return nil
}

func (c *Switch) disruptiveInstall(imageFileName string) error {
	installCmd := fmt.Sprintf("install all nxos bootflash:%s non-interruptive", imageFileName)
	fmt.Println("starting disruptive installation")

	_, err := c.Runner.Run(installCmd, "cli", 0)
	err = checkUpgradeErrorStatus(err)
	if err != nil {
		return err
	}

	rebootCmd := fmt.Sprintf("reload force")
	fmt.Println("Force rebooting the switch...")
	_, err = c.Runner.Run(rebootCmd, "cli", 0)
	err = c.checkNewVersion(imageFileName, true)
	if err != nil {
		return err
	}

	return nil
}

func (c *Switch) nonDisruptiveInstall(imageFileName string) error {
	installCmd := fmt.Sprintf("install all nxos bootflash:%s non-disruptive non-interruptive", imageFileName)
	fmt.Println("starting non-disruptive installation ")

	//Setting Installation Timeout
	i, err := strconv.Atoi(os.Getenv("CISCO_INSTALL_TIME_IN_MINUTES"))
	if err != nil {
		panic("CISCO_INSTALL_TIME_IN_MINUTES was not set as an interger!")
	}
	installTimeDuration := time.Duration(i) * time.Minute
	_, err = c.Runner.Run(installCmd, "cli", installTimeDuration)
	if err != nil {
		if strings.Contains(err.Error(), "failed to get expected string") {
			fmt.Println("Non-disruptive not supported, so starting with Disruptive install.....")
			err = c.disruptiveInstall(imageFileName)
			if err != nil {
				return err
			}
			return nil
		}

		err = c.checkNewVersion(imageFileName, true)
		if err != nil {
			return fmt.Errorf("error while checking version: %+v", err)
		}
	} else {

		rebootCmd := fmt.Sprintf("reload force")
		fmt.Println("Force rebooting the switch...")
		_, err := c.Runner.Run(rebootCmd, "cli", 0)
		err = c.checkNewVersion(imageFileName, false)
		if err != nil {
			return fmt.Errorf("error while checking version: %+v", err)
		}
	}

	return nil
}

func (c *Switch) checkNewVersion(imageFileName string, isSleep bool) error {

	if isSleep {
		b, err := strconv.Atoi(os.Getenv("CISCO_BOOT_TIME_IN_SECONDS"))
		if err != nil {
			panic("CISCO_BOOT_TIME_IN_SECONDS was not set as an interger!")
		}
		bootTimeDuration := time.Duration(b) * time.Second

		// After installation, the switch takes around 10 seconds to reboot, so we need to wait before we run show version
		fmt.Printf("Sleeping for %+v\n", bootTimeDuration)
		time.Sleep(bootTimeDuration)
	}
	fmt.Println("Verifying management connection and version update")

	t, err := strconv.Atoi(os.Getenv("CISCO_RECONNECTION_TIMEOUT_IN_SECONDS"))
	if err != nil {
		panic("CISCO_RECONNECTION_TIMEOUT_IN_SECONDS was not set as an integer!")
	}
	timeoutDuration := time.Duration(t) * time.Second

	timeout := time.NewTimer(timeoutDuration).C
	tick := time.NewTicker(5 * time.Second).C
	for {
		select {
		case <-timeout:
			return errors.New(2, "error while connecting to the switch after update or failed to find the right version")

		case <-tick:
			body, err := c.Runner.Run("show version", "cli", time.Duration(2*time.Second))
			if err == nil {
				if strings.Contains(body, imageFileName) == true {
					fmt.Println("Successfully updgraded to the right version.")
					return nil
				}

			}
		}
	}

}

// GetConfig returns running-config of given switch
func (c *Switch) GetConfig() (string, error) {
	result, err := c.Runner.Run("show running-config", "cli_ascii", 0)
	config := result
	if err != nil {
		return "", fmt.Errorf("error running show running-config command: %+v", err)
	}

	return config, nil
}

// GetFirmware returns Firmware Version of given switch
func (c *Switch) GetFirmware() (string, error) {
	result, err := c.Runner.Run("show version", "cli", 0)
	if err != nil {
		return "", fmt.Errorf("error running show version command: %+v", err)
	}

	var versionBody nexus.CommandRunnerResponseBody

	err = json.Unmarshal([]byte(result), &versionBody)
	if err != nil {
		return "", fmt.Errorf("error getting the result of show version: %+v", err)
	}

	respInterface := versionBody.Result.Body.(map[string]interface{})

	version := respInterface["rr_sys_ver"]
	return version.(string), nil
}

// GetFirmware returns Firmware Version of given switch
func (c *Switch) GetFullVersion() (map[string]interface{}, error) {
	result, err := c.Runner.Run("show version", "cli", 0)
	if err != nil {
		return nil, fmt.Errorf("error running show version command: %+v", err)
	}
	var versionBody nexus.CommandRunnerResponseBody
	err = json.Unmarshal([]byte(result), &versionBody)
	if err != nil {
		return nil, fmt.Errorf("error getting the result of show version: %+v", err)
	}
	respInterface := versionBody.Result.Body.(map[string]interface{})
	delete(respInterface, "TABLE_package_list")
	return respInterface, nil
}

func (c *Switch) CheckVlan( vlanID int64 ) (bool, error){

	command := fmt.Sprintf("show vlan id %v" , vlanID)
	result, err := c.Runner.Run(command, "cli", 0)
	if err != nil {
		return false, fmt.Errorf("error running show vlan command: %+v", err)
	}

	var checkVlanBody nexus.CommandRunnerResponseBody
	err = json.Unmarshal([]byte(result), &checkVlanBody)
	if err != nil {
		return false, fmt.Errorf("error getting the result of show vlan: %+v", err)
	}

	if checkVlanBody.Result.Body != nil  {
		return true, nil
	} else{
		return false, nil
	}
}


func checkUpgradeErrorStatus(err error) error {
	if err != nil {
		// This is workaround for a cisco bug. During upgrade process, the nxapi return timeout error although it successfully upgraded the switch. This was noticed on FW version:  6.0(2)U6(9)
		if !strings.Contains(err.Error(), "Finishing the upgrade, switch will reboot in 10 seconds") {
			return err
		}
	}
	return nil
}
