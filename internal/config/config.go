package config

import (
	"fmt"
	"os"

	"github.com/donaldgifford/rex/internal/templates"
)

// IRexConf interface for creating and generating configs
type IRexConf interface {
	ReadConfig() error
	GenerateDirectories(force bool) error
}

// NewIRexConf generator for IRexConf interface
func NewIRexConf() IRexConf {
	return &RexConf{
		ADR:       NewADRConfig(),
		Templates: templates.NewTemplate(),
	}
}

// RexConf contains the .rex.yaml configuration settings.
type RexConf struct {
	ADR       *ADRConfig
	Templates templates.Template
}

func (rc *RexConf) GetADRConfigSettings() map[string]any {
	return rc.ADR.GetSettings()
}

func (rc *RexConf) GetTemplatesConfigSettings() *templates.Settings {
	return rc.Templates.GetSettings()
}

func (rc *RexConf) GenerateDirectories(force bool) error {
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create directory path string
	dirPath := fmt.Sprintf("%s/%s", cwd, rc.ADR.path)

	// mkdirall with path string
	err = os.MkdirAll(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		// log.Fatal(err)
		return err
	}

	return nil
}

// ReadYamlConfig reads the rex.yaml config in.
// If a config is found it takes the settings in the config file and sets them in the RexConf
func (rc *RexConf) ReadConfig() error {
	fmt.Println(rc.GetADRConfigSettings())
	fmt.Println(rc.GetTemplatesConfigSettings())
	return nil
}
