package fake

import (
	"strings"
	"time"

	"github.com/go-openapi/errors"
)

type FakeRunner struct {
	FailCopyCommand        bool
	FailInstallCommand     bool
	InstallCommand         string
	FailReconnecting       bool
	FailShowVersionCommand bool
	SuccessShowVersion     bool
	ImageFileName          string
}

func (fr *FakeRunner) Run(command string, timeout time.Duration) (string, error) {
	if strings.Contains(command, "copy") && fr.FailCopyCommand {
		return "", errors.New(1, "fake copy command failed")
	}

	if strings.Contains(command, "install") {
		fr.InstallCommand = command

		if fr.FailInstallCommand {
			return "", errors.New(2, "fake install command failed")
		}
	}

	if strings.Contains(command, "show") && fr.FailReconnecting {
		time.Sleep(10 * time.Second)
		return "", errors.New(2, "failed to reconnect")
	}

	if strings.Contains(command, "show") && fr.FailShowVersionCommand {
		return "fake.bin", nil
	}

	if strings.Contains(command, "copy") && fr.SuccessShowVersion {
		fr.ImageFileName = strings.Split((strings.Split(command, " ")[2]), ":")[1]
	}

	if strings.Contains(command, "show") && fr.SuccessShowVersion {
		return fr.ImageFileName, nil
	}

	return "", nil
}
