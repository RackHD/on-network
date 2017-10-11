package switch_operations

// ISwitch is an interface for switch
type ISwitch interface {
	Update(string, string) error
	GetConfig() (string, error)
	GetFirmware()(string, error)
	GetFullVersion()(map[string] interface{}, error)
}
