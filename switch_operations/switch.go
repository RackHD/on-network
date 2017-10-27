package switch_operations

import "github.com/RackHD/on-network/models"

// ISwitch is an interface for switch
type ISwitch interface {
	Update(string, []*models.FirmwareImage) error
	GetConfig() (string, error)
	GetFirmware() (string, error)
	GetFullVersion() (map[string]interface{}, error)
}
