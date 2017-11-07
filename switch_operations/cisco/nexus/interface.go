// Copyright 2017, Dell EMC, Inc.

package nexus

import "time"

type CommandRunner interface {
	Run(string, string, time.Duration) (string, error)
}
