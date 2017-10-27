package nexus

import "time"

type CommandRunner interface {
	Run(string, string, time.Duration) (string, error)
}
