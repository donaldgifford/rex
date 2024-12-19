package config

import (
	"fmt"
	"os"

	"github.com/donaldgifford/rex/internal/templates"
	"github.com/spf13/viper"
)

func NewRexConfInstall() *RexConfInstall {
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
