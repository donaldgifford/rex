package config

import (
	"fmt"
	"os"

	"github.com/donaldgifford/rex/internal/adr"
	"github.com/donaldgifford/rex/internal/templates"
)

// IRexConf interface for creating and generating configs
type IRexConf interface {
	ReadConfig() error
	GenerateDirectories(force bool, index bool) error
	Settings() *ConfigSettings
	CreateADR(adr *adr.ADR) error
	CreateIndex(idx *adr.Index) error
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

func (rc *RexConf) GetADRConfigSettings() *ADRConfig {
	return rc.ADR.GetSettings()
}

func (rc *RexConf) GetTemplatesConfigSettings() *templates.Settings {
	return rc.Templates.GetSettings()
}

func (rc *RexConf) GenerateDirectories(force bool, index bool) error {
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create directory path string
	dirPath := fmt.Sprintf("%s/%s", cwd, rc.ADR.Path)

	// mkdirall with path string
	err = os.MkdirAll(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		// log.Fatal(err)
		return err
	}

	// if index true, create index file from rex.conf and template
	idx := adr.NewIndex()

	err = rc.Templates.CreateIndex(idx)
	if err != nil {
		return err
	}

	return nil
}

type ConfigSettings struct {
	ADR       *ADRConfig
	Templates *templates.Settings
}

// ReadYamlConfig reads the rex.yaml config in.
// If a config is found it takes the settings in the config file and sets them in the RexConf
func (rc *RexConf) ReadConfig() error {
	fmt.Println(rc.GetADRConfigSettings())
	fmt.Println(rc.GetTemplatesConfigSettings())
	return nil
}

func (rc *RexConf) Settings() *ConfigSettings {
	return &ConfigSettings{
		ADR:       rc.GetADRConfigSettings(),
		Templates: rc.GetTemplatesConfigSettings(),
	}
}

func (rc *RexConf) CreateADR(adr *adr.ADR) error {
	err := rc.Templates.CreateADR(adr)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RexConf) CreateIndex(idx *adr.Index) error {
	err := rc.Templates.CreateIndex(idx)
	if err != nil {
		return err
	}

	return nil
}
