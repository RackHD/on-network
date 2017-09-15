package nexus_interface

type CommandRunner interface {
	Run(string) (string, error)
}
