package cisco

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/RackHD/on-network/switch_operations/cisco/nexus"
	"github.com/RackHD/on-network/switch_operations/store"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"
	//"github.com/revel/config"
)

type Switch struct {
	Runner nexus.CommandRunner
}

func (c *Switch) Update(switchModel, imageURL string) error {
	switchesDatabase := store.GetSwitchFileDatabase() //switch this if use database

	updateType, err := switchesDatabase.GetUpdateType("cisco", switchModel)
	if err != nil {
		return err
	}

	imageFileName := fmt.Sprintf("%s-%s", uuid.New().String(), path.Base(imageURL))
	fmt.Println("filename", imageFileName)

	copyCmd := fmt.Sprintf("copy %s bootflash:%s vrf management", imageURL, imageFileName)
	fmt.Println("starting copy")
	_, err = c.Runner.Run(copyCmd, "cli",0)
	if err != nil {
		return fmt.Errorf("error copying image from remote: %+v", err)
	}

	var installCmd string

	if updateType == "Disruptive" {
		installCmd = fmt.Sprintf("install all nxos bootflash:%s non-interruptive", imageFileName)
		fmt.Println("starting disruptive installation")
		_, err = c.Runner.Run(installCmd,"cli", 0)
		if err != nil {
			return fmt.Errorf("error install image: %+v", err)
		}

	} else if updateType == "NonDisruptive" {
		installCmd = fmt.Sprintf("install all nxos bootflash:%s non-disruptive non-interruptive", imageFileName)
		fmt.Println("starting non-disruptive installation ")
		_, err = c.Runner.Run(installCmd, "cli",2 * time.Second)
		if err != nil {

			i, err := strconv.Atoi(os.Getenv("CISCO_INSTALL_TIME_IN_MINUTES"))
			if err != nil {
				panic("CISCO_INSTALL_TIME_IN_MINUTES was not set as an interger!")
			}
			installTimeDuration := time.Duration(i) * time.Minute

			rebootTimeout :=time.NewTimer(installTimeDuration)
			rebootTick := time.NewTicker(5 *time.Second)
			isBreak := false
			for {
				select {
				case <-rebootTimeout.C:
					return errors.New(2, "Something went wrong during installation, switch never rebooted" )
				case <-rebootTick.C:
					_, err := c.Runner.Run("show version", "cli",time.Duration(6*time.Second))

					if err != nil {
						fmt.Println("Installation completed, and switch is rebooting.")
						rebootTimeout.Stop()
						rebootTick.Stop()
						isBreak=true
					}
				}
				if (isBreak) {
					break
				}
			}
		} else {
			return errors.New(2, "Something went wrong during installation." )
		}
	}
	b, err := strconv.Atoi(os.Getenv("CISCO_BOOT_TIME_IN_SECONDS"))
	if err != nil {
		panic("CISCO_BOOT_TIME_IN_SECONDS was not set as an interger!")
	}
	bootTimeDuration := time.Duration(b) * time.Second

	// After installation, the switch takes around 10 seconds to reboot, so we need to wait before we run show version
	fmt.Printf("Sleeping for %+v\n", bootTimeDuration)
	time.Sleep(bootTimeDuration)

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
			return errors.New(2, "timeout connecting to switch after update.")

		case <-tick:
			body, err := c.Runner.Run("show version","cli", time.Duration(2*time.Second))
			if err == nil {
				if strings.Contains(body, imageFileName) == true {
					fmt.Println("Successfully updgraded to the right version.")
					return nil
				}
				return errors.New(3, "failed to find the expected version")
			}
		}
	}
}

// GetConfig returns running-config of given switch
func (c *Switch) GetConfig() (string, error) {
	result, err := c.Runner.Run("show running-config", "cli_ascii", 0)
	//result, err := c.Runner.Run("show version", "cli", 0)
	config := result
	if err != nil {
		return "", fmt.Errorf("error running show running-config command: %+v", err)
	}

	return config, nil
}
