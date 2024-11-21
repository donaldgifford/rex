package config

// IRexConf interface for creating and generating configs
type IRexConf interface {
	ReadYamlConfig() error
	GenerateYamlFile() error
}

// NewIRexConf generator for IRexConf interface
func NewIRexConf() IRexConf {
	return &RexConf{}
}

// RexConf contains the .rex.yaml configuration settings.
type RexConf struct {
	ADR       *ADRConfig
	Templates *TemplateConfig
}

// getRexConfig reads in the .rex.yaml settings and creates a RexConf to use.
func getRexConfig() RexConf {
	return RexConf{
		newADRConfig(),
		newTemplateConfig(),
	}
}

func (rc *RexConf) GetADRConfigSettings() map[string]any {
	return rc.ADR.GetSettings()
}

// ReadYamlConfig reads the rex.yaml config in.
// If a config is found it takes the settings in the config file and sets them in the RexConf
func (rc *RexConf) ReadYamlConfig() error {
	return nil
}

// GenerateRexYaml creates a default rex.yaml file in the current working directory
// If a .rex.yaml file is found, GenerateYamlFile will validate its settings to be able to use it in a RexConf
func (rc *RexConf) GenerateYamlFile() error {
	return nil
}

func (rc *RexConf) ReadConfig() {}
