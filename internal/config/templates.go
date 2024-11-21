package config

import (
	"github.com/spf13/viper"
)

type TemplateConfig struct {
	path string
	adr  string
}

func newTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		path: viper.GetString(""),
		adr:  viper.GetString(""),
	}
}
