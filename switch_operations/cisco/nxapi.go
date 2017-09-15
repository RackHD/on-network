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
	imageFileName := uuid.New().String() + path.Base(imageURL)
	copyCmd := fmt.Sprintf("copy %s bootflash:%s vrf management", imageURL, imageFileName)
	fmt.Println("starting copy")
	_, err := c.Runner.Run(copyCmd)
	if err != nil {
		return fmt.Errorf("error copying image from remote: %+v", err)
	}

	installCmd := fmt.Sprintf("install all nxos bootflash:%s non-interruptive", imageFileName)
	fmt.Println("starting installation")
	_, err = c.Runner.Run(installCmd)
	if err != nil {
		return fmt.Errorf("error install image: %+v", err)
	}

	//After installation, the switch takes around 10 seconds to reboot, so we need to wait before we run show version
	time.Sleep(20 * time.Second)

	fmt.Println("Verifying management connection and version upgrade")
	timeout := time.After(15 * time.Second)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeout:
			return errors.New(2, "timeout connecting to switch after upgrade.")

		case <-tick:
			body, err := c.Runner.Run("show version")
			if err == nil {
				if strings.Contains(body, imageFileName) == true {
					return nil
				} else {
					return errors.New(3, "failed to find the expected version")
				}
			}
		}
	}
}
