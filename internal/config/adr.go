package config

import "github.com/spf13/viper"

var (
	rexConfigFile        string = ".rex.yaml"
	rexConfigFileContent string = `adr:
  path: "docs/adr/"
  index_page: "README.md"
  add_to_index: true # on rex create, a new record will be added to the index page
templates:
  path: "templates/"
  adr:
    default: "adr.tmpl"`
)

type ADRConfig struct {
	path       string
	indexPage  string
	addToIndex bool
}

// newADRConfig reads the configuration settings under "adr"
func newADRConfig() *ADRConfig {
	return &ADRConfig{
		path:       viper.GetString(""),
		indexPage:  viper.GetString(""),
		addToIndex: viper.GetBool(""),
	}
}

// GetSettings returns the settings for ADR
func (a *ADRConfig) GetSettings() map[string]any {
	settings := make(map[string]any)
	settings["path"] = a.path
	settings["indexPage"] = a.indexPage
	settings["addToIndex"] = a.addToIndex

	return settings
}
