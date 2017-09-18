package nexus_interface

import "time"

type CommandRunner interface {
	Run(string, time.Duration) (string, error)
}
