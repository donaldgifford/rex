package config

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type RexConfigure interface {
	Settings() *RexConfig
	YamlOut() ([]byte, error)
}

func NewRexConfigure() RexConfigure {
	return NewRexConfig()
}

// RexConfig holds configuration from .rex.yaml
type RexConfig struct {
	ADR               ADRConfig      `yaml:"adr"`
	Templates         TemplateConfig `yaml:"templates"`
	EnableGithubPages bool           `yaml:"enable_github_pages"`
	Pages             PagesConfig    `yaml:"pages"`
	Extras            bool           `yaml:"extras"`
	ExtraPages        ExtrasConfig   `yaml:"extra_pages"`
}

type ADRConfig struct {
	Path       string `yaml:"path"`
	IndexPage  string `yaml:"index_page"`
	AddToIndex bool   `yaml:"add_to_index"`
}

type ADRTemplateConfig struct {
	Default string `yaml:"default"`
	Index   string `yaml:"index"`
}

type TemplateConfig struct {
	Enabled bool              `yaml:"enabled"`
	Path    string            `yaml:"path"`
	ADR     ADRTemplateConfig `yaml:"adr"`
}

type PagesConfig struct {
	Index string         `yaml:"index"`
	Web   PagesConfigWeb `yaml:"web"`
}

type PagesConfigWeb struct {
	Config string               `yaml:"config"`
	Layout PagesConfigWebLayout `yaml:"layout"`
}

type PagesConfigWebLayout struct {
	ADR     string `yaml:"adr"`
	Default string `yaml:"default"`
}

type ExtrasConfig struct {
	Install string `yaml:"install"`
	Usage   string `yaml:"usage"`
}

// NewRexConfig creates an empty config object
func NewRexConfig() *RexConfig {
	return &RexConfig{
		ADR: ADRConfig{
			Path:       viper.GetString("adr.path"),
			IndexPage:  viper.GetString("adr.index_page"),
			AddToIndex: viper.GetBool("adr.add_to_index"),
		},
		Templates: TemplateConfig{
			Enabled: viper.GetBool("templates.enabled"),
			Path:    viper.GetString("templates.path"),
			ADR: ADRTemplateConfig{
				Default: viper.GetString("templates.adr.default"),
				Index:   viper.GetString("templates.adr.index"),
			},
		},
		EnableGithubPages: viper.GetBool("enable_github_pages"),
		Pages: PagesConfig{
			Index: viper.GetString("pages.index"),
			Web: PagesConfigWeb{
				Config: viper.GetString("pages.web.config"),
				Layout: PagesConfigWebLayout{
					ADR:     viper.GetString("pages.web.layout.adr"),
					Default: viper.GetString("pages.web.layout.default"),
				},
			},
		},
		Extras: viper.GetBool("extras"),
		ExtraPages: ExtrasConfig{
			Install: viper.GetString("extra_pages.install"),
			Usage:   viper.GetString("extra_pages.usage"),
		},
	}
}

// Settings exposes settings out to use in other calls
func (rc *RexConfig) Settings() *RexConfig {
	return rc
}

func (rc *RexConfig) YamlOut() ([]byte, error) {
	yamlData, err := yaml.Marshal(&rc)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}
