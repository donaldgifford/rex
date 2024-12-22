package config

import (
	"fmt"
	"os"

	"github.com/donaldgifford/rex/internal/adr"
	"github.com/donaldgifford/rex/internal/templates"
	"github.com/spf13/viper"
)

func NewRexConfigInstall() *RexConfInstall {
	return &RexConfInstall{}
}

type RexConfInstall struct{}

func (rc *RexConfInstall) ConfigExists() bool {
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Config file not found, `run rex config create` to generate one", viper.ConfigFileUsed())
		return false
	}
	return true
}

// ReadYamlConfig reads the rex.yaml config in.
// If a config is found it takes the settings in the config file and sets them in the RexConf
func (rc *RexConfInstall) ReadConfig() error {
	return nil
}

func (rc *RexConfInstall) Settings() *RexConfig {
	return nil
}

func (rc *RexConfInstall) WriteConfig(file string) error {
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
func (rc *RexConfInstall) GenerateConfig(force bool) error {
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

// GenerateDirectories will create directories in the default setting
func (rc *RexConfInstall) GenerateDirectories(force bool, index bool) error {
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create directory path string
	dirPath := fmt.Sprintf("%s/%s", cwd, "docs/adr/")

	// mkdirall with path string
	err = os.MkdirAll(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		// log.Fatal(err)
		return err
	}

	// if index true, create index file from rex.conf and template
	if index {
		err = rc.GenerateIndex("default/index_readme.tmpl")
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *RexConfInstall) GenerateIndex(file string) error {
	// get template to be used
	idx, err := templates.DefaultRexTemplates.ReadFile(file)
	if err != nil {
		return err
	}
	fmt.Println(string(idx))

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create the file
	fileName := cwd + "/docs/adr/" + "README.md"
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	// write the file
	_, err = f.Write(idx)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RexConfInstall) CreateIndex() error {
	t := templates.NewTemplate()

	idx := adr.NewIndex()
	err := t.GenerateIndex(idx)
	if err != nil {
		return err
	}

	return nil
}
