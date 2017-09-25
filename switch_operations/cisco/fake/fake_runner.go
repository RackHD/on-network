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
}

func (fr *FakeRunner) Run(command string, timeout time.Duration) (string, error) {
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
	}

	if strings.Contains(command, "show running-config") {
		if fr.FailShowConfigCommand {
			return "", errors.New(4, "fake show config command failed")
		}
		return "{\"config\":\"empty\"}", nil
	}

	return "", nil
}
