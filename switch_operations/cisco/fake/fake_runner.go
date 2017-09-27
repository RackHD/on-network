package fake

import (
	"strings"
	"time"

	"github.com/go-openapi/errors"
)

type FakeRunner struct {
	// Update
	FailCopyCommand        bool
	FailInstallCommand     bool
	InstallCommand         string
	FailReconnecting       bool
	FailShowVersionCommand bool
	SuccessShowVersion     bool
	ImageFileName          string
	// GetConfig
	FailShowConfigCommand bool

	//Get Firmware Version
	FailShowFirmwareVersionCommand bool

	TimeoutInstall		   bool
	IsShowVersion          bool

}

func (fr *FakeRunner) Run(command string,method string,   timeout time.Duration) (string, error) {


	if strings.Contains(command, "copy") {
		fr.ImageFileName = strings.Split((strings.Split(command, " ")[2]), ":")[1]

		if fr.FailCopyCommand {
			return "", errors.New(1, "fake copy command failed")
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
	}

	if strings.Contains(command, "show") {
		if fr.FailReconnecting {
			time.Sleep(10 * time.Second)
			return "", errors.New(2, "failed to reconnect")
		}

		if fr.FailShowVersionCommand {
			return "fake.bin", nil
		}

		if fr.SuccessShowVersion {
			return fr.ImageFileName, nil
		}
		if fr.TimeoutInstall {
			if (fr.IsShowVersion == false) {
				fr.IsShowVersion = true
				return "", errors.New(2, "failed as switch is rebooting")
			} else{
				return fr.ImageFileName, nil
			}
		}
	}

	if strings.Contains(command, "show running-config") {
		if fr.FailShowConfigCommand {
			return "", errors.New(4, "fake show config command failed")
		}
		return "{\"config\":\"empty\"}", nil
	}

	if strings.Contains(command, "show version") {
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

	return "", nil
}
