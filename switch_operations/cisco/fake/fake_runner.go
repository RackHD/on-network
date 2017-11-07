package fake

import (
	"strings"
	"time"

	"github.com/go-openapi/errors"
)

type FakeRunner struct {
	// Update
	FailCopyCommand           bool
	FailCopyCommandWrongParam bool
	FailInstallCommand        bool
	InstallCommand            string
	FailReconnecting          bool
	FailShowVersionCommand    bool
	SuccessShowVersion        bool
	ImageFileName             []string
	DowngradeNonDisruptive    bool

	// GetConfig
	FailShowConfigCommand bool

	//Get Firmware Version
	FailShowFirmwareVersionCommand bool

	//Check Vlan
	FailCheckVlanCommand bool
	SuccessShowVlanNullResult bool

	TimeoutInstall bool
	IsShowVersion  bool
	CopyCounter    int
}

func (fr *FakeRunner) Run(command string, method string, timeout time.Duration) (string, error) {

	if strings.Contains(command, "copy") {
		fr.ImageFileName = append(fr.ImageFileName, strings.Split((strings.Split(command, " ")[2]), ":")[1])

		if fr.FailCopyCommand {
			return "", errors.New(1, "fake copy command failed")
		}

		if fr.FailCopyCommandWrongParam {
			return "", errors.New(1, "Missing required image type")
		}
	}

	if strings.Contains(command, "install") {
		fr.InstallCommand = command

		if fr.FailInstallCommand {
			return "", errors.New(2, "fake install command failed")
		}

		if fr.TimeoutInstall {
			return "", errors.New(2, "fake install command timedout")
		}

		if fr.DowngradeNonDisruptive && strings.Contains(command, "non-disruptive") {
			return "", errors.New(2, "failed to get expected string. status code: 500")
		}
	}

	if strings.Contains(command, "show version") {
		if fr.FailReconnecting {
			time.Sleep(10 * time.Second)
			return "", errors.New(2, "failed to reconnect")
		}

		if fr.FailShowVersionCommand {
			return "fake.bin", nil
		}

		if fr.SuccessShowVersion {
			imageFileName := fr.ImageFileName[fr.CopyCounter]
			fr.CopyCounter++
			return imageFileName, nil
		}
		if fr.TimeoutInstall {
			if fr.IsShowVersion == false {
				fr.IsShowVersion = true
				return "", errors.New(2, "failed as switch is rebooting")
			} else {
				imageFileName := fr.ImageFileName[fr.CopyCounter]
				fr.CopyCounter++
				return imageFileName, nil
			}
		}
		if fr.FailShowFirmwareVersionCommand {
			return "", errors.New(4, "fake show version command failed")
		}
		return `{
		"jsonrpc": "2.0",
		"result": {
			"body": {
				"rr_sys_ver": "7.0(3)I5(2)"
			}
		},
		"id": 1
		}`, nil
	}

	if strings.Contains(command, "show running-config") {
		if fr.FailShowConfigCommand {
			return "", errors.New(4, "fake show config command failed")
		}
		return "{\"config\":\"empty\"}", nil
	}

	if strings.Contains(command, "show vlan id") {
		if fr.FailCheckVlanCommand {
			return "", errors.New(4, "Invalid value/range")
		}
		if fr.SuccessShowVlanNullResult {
			return `{
				"jsonrpc": "2.0",
					"result": null,
					"id": 1
			}`, nil
		}
		return `{
			  "jsonrpc": "2.0",
			  "result": {
				"body": {
				  "vlanshowrspan-vlantype": "notrspan"
				}
			  },
			  "id": 1
			}`,nil
	}

	return "", nil
}
