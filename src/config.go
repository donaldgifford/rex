package src

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config interface {
	Init()
	createDir(string)
	generateFile(string, RexFile)
	generateADR()
	check(error)
}

func Init() {
	config, err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	runConfig(config)
}

func checkConfigFileExists() {
	cwd := getCwd()
	file, err := os.Open(cwd + ".rex.yaml")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("No .rex.yaml file exists. Creating a new one...")
		createDefaultConfigFile(cwd)
	} else if err != nil {
		fmt.Println("Error opening .rex.yaml file: ", err)
	} else {
		fmt.Printf("File exists: &s\n", file.Name())
	}
}

func runConfig(c Config) {
	c.Init()
}

func getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return cwd
}

func InitConfig() (Config, error) {
	if viper.GetBool("enable_github_pages") {
		return &ConfigGithubPages{
			adrPath:       viper.GetString("adr.path"),
			templatesPath: viper.GetString("templates.path"),
			cwd:           getCwd(),
		}, nil
	} else {
		return &ConfigAdr{
			adrPath:       viper.GetString("adr.path"),
			templatesPath: viper.GetString("templats.path"),
			cwd:           getCwd(),
		}, nil
	}
}

func CreateConfig() {
	cwd := getCwd()

	createDefaultConfigFile(cwd)
}

func createDefaultConfigFile(cwd string) {
	rexConfigFile := ".rex.yaml"

	rexConfig := `adr:
  path: "docs/adr/"
  index_page: "README.md"
  add_to_index: true # on rex create, a new record will be added to the index page
templates:
  path: "templates/"
  adr:
    default: "adr.tmpl"
enable_github_pages: false`
	rc := []byte(rexConfig)
	fileName := cwd + "/" + rexConfigFile
	fmt.Println(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(rc)
	if err != nil {
		log.Fatal(err)
	}
}
