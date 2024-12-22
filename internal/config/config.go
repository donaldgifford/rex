package config

import (
	"fmt"
	"os"

	"github.com/donaldgifford/rex/internal/templates"
	"github.com/spf13/viper"
)

// IRexConf interface for creating and generating configs
// type IRexConf interface {
// 	ReadConfig() error
// 	GenerateDirectories(force bool, index bool) error
// 	Settings() *ConfigSettings
// 	CreateADR(adr *adr.ADR) error
// 	CreateIndex() error
// }
//
// // NewIRexConf generator for IRexConf interface
// func NewIRexConf() IRexConf {
// 	return &RexConf{
// 		ADR:       NewADRConfig(),
// 		Templates: templates.NewTemplate(),
// 	}
// }
//
// // RexConf contains the .rex.yaml configuration settings.
// type RexConf struct {
// 	ADR       *ADRConfig
// 	Templates templates.Template
// }

// type RexConfigurer interface {
// 	// ReadConfig() error
// 	WriteConfig(file string) error
// 	ConfigExists() bool
// 	GenerateConfig(force bool) error
// 	// GenerateDirectories(force bool, index bool) error
// 	Settings() *RexConfig
// }
//
// func NewRexConfigurer(install bool) RexConfigurer {
// 	if install {
// 		return NewRexConfigInstall()
// 	} else {
// 		return NewRexConfig()
// 	}
// }

// RexConfig holds configuration from .rex.yaml
type RexConfig struct {
	ADR               ADRConfig
	Templates         TemplateConfig
	EnableGithubPages bool
	Pages             PagesConfig
	Extras            bool
	ExtraPages        ExtrasConfig
}

type ADRConfig struct {
	Path       string
	IndexPage  string
	AddToIndex bool
}

type ADRTemplateConfig struct {
	Default string
	Index   string
}

type TemplateConfig struct {
	Enabled bool
	Path    string
	ADR     ADRTemplateConfig
}

type PagesConfig struct {
	Index string
	Web   PagesConfigWeb
}

type PagesConfigWeb struct {
	Config string
	Layout PagesConfigWebLayout
}

type PagesConfigWebLayout struct {
	ADR     string
	Default string
}

type ExtrasConfig struct {
	Install string
	Usage   string
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

	// check if config exists so not to accidently overwrite your config
	if rc.ConfigExists() {
		fmt.Println("Config already exists. Use --force option to overwrite it.")
		return nil
	}

	// write the config file
	err := rc.WriteConfig("default/rex.yaml")
	if err != nil {
		return err
	}
	return nil
}
func GenerateDirs()  {}
func GenerateIndex() {}

// func (rc *RexConf) GetADRConfigSettings() *ADRConfig {
// 	return rc.ADR.GetSettings()
// }
//
// func (rc *RexConf) GetTemplatesConfigSettings() *templates.Settings {
// 	return rc.Templates.GetSettings()
// }
//
// func (rc *RexConf) GenerateDirectories(force bool, index bool) error {
// 	// get current working directory
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}
//
// 	// create directory path string
// 	dirPath := fmt.Sprintf("%s/%s", cwd, rc.ADR.Path)
//
// 	// mkdirall with path string
// 	err = os.MkdirAll(dirPath, 0755)
// 	if err != nil && !os.IsExist(err) {
// 		// log.Fatal(err)
// 		return err
// 	}
//
// 	// if index true, create index file from rex.conf and template
// 	if index {
// 		err = rc.CreateIndex()
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// type ConfigSettings struct {
// 	ADR       *ADRConfig
// 	Templates *templates.Settings
// }
//
// // ReadYamlConfig reads the rex.yaml config in.
// // If a config is found it takes the settings in the config file and sets them in the RexConf
// func (rc *RexConf) ReadConfig() error {
// 	fmt.Println(rc.GetADRConfigSettings())
// 	fmt.Println(rc.GetTemplatesConfigSettings())
// 	return nil
// }
//
// func (rc *RexConf) Settings() *ConfigSettings {
// 	return &ConfigSettings{
// 		ADR:       rc.GetADRConfigSettings(),
// 		Templates: rc.GetTemplatesConfigSettings(),
// 	}
// }
//
// func (rc *RexConf) CreateADR(adr *adr.ADR) error {
// 	err := rc.Templates.CreateADR(adr)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (rc *RexConf) CreateIndex() error {
// 	idx := adr.NewIndex()
//
// 	err := idx.GetIndexAdrs()
// 	if err != nil {
// 		return err
// 	}
//
// 	err = rc.Templates.GenerateIndex(idx)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
