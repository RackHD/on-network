package api

type CommandRunner interface {
	Run(string) (string, error)
}
