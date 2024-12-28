package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/donaldgifford/rex/internal/templates"
)

type RexConfigurer interface {
	ConfigExists() bool
	ReadConfig() error
	GenerateConfig(force bool) error
	GenerateIndex(force bool) error
	GenerateDirectories(force bool, index bool) error
	Settings() *RexConfig
}

func NewRexConfigurer(install bool) RexConfigurer {
	if install {
		return NewRexConfigInstall()
	} else {
		return NewRexConfig()
	}
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

func (rc *RexConfig) ConfigExists() bool {
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Config file not found, `run rex config create` to generate one", viper.ConfigFileUsed())
		return false
	}
	return true
}

// ReadYamlConfig reads the rex.yaml config in.
// If a config is found it takes the settings in the config file and sets them in the RexConf
// Used to re read the config in before any changes.
func (rc *RexConfig) ReadConfig() error {
	rc.ADR.Path = viper.GetString("adr.path")
	rc.ADR.IndexPage = viper.GetString("adr.index_page")
	rc.ADR.AddToIndex = viper.GetBool("adr.add_to_index")

	rc.Templates = TemplateConfig{
		Enabled: viper.GetBool("templates.enabled"),
		Path:    viper.GetString("templates.path"),
		ADR: ADRTemplateConfig{
			Default: viper.GetString("templates.adr.default"),
			Index:   viper.GetString("templates.adr.index"),
		},
	}

	rc.EnableGithubPages = viper.GetBool("enable_github_pages")

	rc.Pages = PagesConfig{
		Index: viper.GetString("pages.index"),
		Web: PagesConfigWeb{
			Config: viper.GetString("pages.web.config"),
			Layout: PagesConfigWebLayout{
				ADR:     viper.GetString("pages.web.layout.adr"),
				Default: viper.GetString("pages.web.layout.default"),
			},
		},
	}
	rc.Extras = viper.GetBool("extras")
	rc.ExtraPages = ExtrasConfig{
		Install: viper.GetString("extra_pages.install"),
		Usage:   viper.GetString("extra_pages.usage"),
	}

	fmt.Println(rc)

	return nil
}

func (rc *RexConfig) WriteConfig(file string) error {
	fmt.Println("Creating new config file at .rex.yaml")

	// get template to be used
	rexConf, err := templates.DefaultRexTemplates.ReadFile(file)
	if err != nil {
		return err
	}
	fmt.Println(string(rexConf))

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create the file
	fileName := cwd + "/" + ".rex.yaml"
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	// write the file
	_, err = f.Write(rexConf)
	if err != nil {
		return err
	}

	return nil
}

// GenerateRexYaml creates a default rex.yaml file in the current working directory
// If a .rex.yaml file is found, GenerateYamlFile will validate its settings to be able to use it in a RexConf
func (rc *RexConfig) GenerateConfig(force bool) error {
	// if force is true, overwrite the config file
	if force {
		err := rc.WriteConfig("default/rex.yaml")
		if err != nil {
			return err
		}
		return nil
	}

	// check if config exists so not to accidentally overwrite your config
	if rc.ConfigExists() {
		fmt.Println("Config already exists. Use --force option to overwrite it.")
		return errors.New("config already exists. Use --force option to overwrite it")
	}

	// write the config file
	err := rc.WriteConfig("default/rex.yaml")
	if err != nil {
		return err
	}
	return nil
}

func (rc *RexConfig) GenerateDirectories(force bool, index bool) error {
	return nil
}

func (rc *RexConfig) GenerateIndex(force bool) error {
	return nil
}

func (rc *RexConfig) YamlOut() error {
	yamlData, err := yaml.Marshal(&rc)
	if err != nil {
		return err
	}

	fmt.Println(string(yamlData))
	return nil
}

func (rc *RexConfig) YamlIn() error {
	t, err := os.ReadFile(".rex.yaml")
	if err != nil {
		return err
	}

	var config RexConfig

	err = yaml.Unmarshal(t, &config)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	fmt.Println(&config)

	return nil
}
