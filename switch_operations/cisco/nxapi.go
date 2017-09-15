package cisco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/RackHD/on-network/switch_operations/cisco/api"
	"github.com/go-openapi/errors"
	"github.com/google/uuid"
)

type Switch struct {
	Runner api.CommandRunner
}

type Params struct {
	Command string `json:"cmd"`
	Version int    `json:"version"`
}

type CommandRunnerBody struct {
	JsonRpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}

type CopyCommand struct {
	Src string
	Dst string
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

type NexusRunner struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (nr *NexusRunner) Run(command string) (string, error) {
	endpoint := fmt.Sprintf("http://%s/ins", nr.IP)

	commandParam := Params{command, 1}
	postBody := CommandRunnerBody{"2.0", "cli", commandParam, 1}
	bodyBytes, err := json.Marshal(postBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling command: %+v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error making request: %+v", err)
	}

	req.Header.Set("Content-Type", "application/json-rpc")
	req.SetBasicAuth(nr.Username, nr.Password)

	fmt.Println("making request to nxos...")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %+v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return string(body), fmt.Errorf("error reading response body: %+v", err)
	}

	if resp.StatusCode != 200 {
		return string(body), errors.New(1, "failed to get expected string. status code: %d\nbody: %+v", resp.StatusCode, string(body))
	}

	return string(body), nil
}
