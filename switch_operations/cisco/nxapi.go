package cisco

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/RackHD/on-network/switch_operations/cisco/nexus_interface"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"
)

type Switch struct {
	Runner nexus_interface.CommandRunner
}

func (c *Switch) Update(imageURL string) error {
	imageFileName := fmt.Sprintf("%s-%s", uuid.New().String(), path.Base(imageURL))
	fmt.Println("filename", imageFileName)
	copyCmd := fmt.Sprintf("copy %s bootflash:%s vrf management", imageURL, imageFileName)
	fmt.Println("starting copy")
	_, err := c.Runner.Run(copyCmd, 0)
	if err != nil {
		return fmt.Errorf("error copying image from remote: %+v", err)
	}

	installCmd := fmt.Sprintf("install all nxos bootflash:%s non-interruptive", imageFileName)
	fmt.Println("starting installation")
	_, err = c.Runner.Run(installCmd, 0)
	if err != nil {
		return fmt.Errorf("error install image: %+v", err)
	}

	// After installation, the switch takes around 10 seconds to reboot, so we need to wait before we run show version
	fmt.Println("Sleeping for 20 seconds")
	time.Sleep(20 * time.Second)

	fmt.Println("Verifying management connection and version update")

	timeout := time.NewTimer(3 * time.Minute).C
	tick := time.NewTicker(5 * time.Second).C
	for {
		select {
		case <-timeout:
			return errors.New(2, "timeout connecting to switch after update.")

		case <-tick:
			body, err := c.Runner.Run("show version", time.Duration(2*time.Second))
			if err == nil {
				if strings.Contains(body, imageFileName) == true {
					return nil
				}
				return errors.New(3, "failed to find the expected version")
			}
		}
	}
}
