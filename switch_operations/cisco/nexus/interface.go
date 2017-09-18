package nexus

import "time"

type CommandRunner interface {
	Run(string, time.Duration) (string, error)
}
