package store

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type ISwitchDatabase interface {
	GetUpdateType(switchType, switchModel string) (string, error)
}

type Switch struct {
	Name   string `yaml: "name"`
	Models []struct {
		Name       string `yaml: "name"`
		Disruptive bool   `yaml:"disruptive"`
		Firmware   string `yaml:"firmware"`
	} `yaml :"models"`
}

type SwitchType struct {
	Switches []Switch `yaml: "switches"`
}

type SwitchFileDatabase struct {
	Switches []Switch
}

var switchFileDatabaseInstance *SwitchFileDatabase

func GetSwitchFileDatabase() *SwitchFileDatabase {
	if switchFileDatabaseInstance == nil {
		switches, err := GetSwitches()
		if err != nil {
			panic(err)
		}

		switchFileDatabaseInstance = &SwitchFileDatabase{
			Switches: switches,
		}
	}
	return switchFileDatabaseInstance
}

func (so *SwitchFileDatabase) GetUpdateType(switchType, switchModel string) (string,string, error) {
	for _, stype := range so.Switches {
		if stype.Name == switchType {
			for _, smodels := range stype.Models {
				if strings.Contains(strings.ToLower(switchModel), strings.ToLower(smodels.Name)) {
					if smodels.Disruptive == true {
						return "Disruptive", smodels.Firmware,  nil
					}
					return "NonDisruptive",smodels.Firmware,  nil
				}
			}
			return "","", errors.New("couldn't find switch model")
		}
	}
	return "","", errors.New("couldn't find switch type")
}

func GetSwitches() ([]Switch, error) {
	path := os.Getenv("SWITCH_MODELS_FILE_PATH")
	if path == "" {
		return nil, errors.New("SWITCH_MODELS_FILE_PATH was not set")
	}

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switchType := SwitchType{}
	err = yaml.Unmarshal([]byte(fileData), &switchType)
	if err != nil {
		return nil, err
	}

	return switchType.Switches, nil
}
