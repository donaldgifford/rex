/* package src

init.go

Methods and functions for initializing the rex tool.

Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/

package src

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// create directory structure with default files

func Init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//	createDefaultConfigFile(cwd)
	createDir(cwd, "adr.path")
	createDir(cwd, "templates.path")
	createADRDefaultTemplate(cwd, "templates.path")
}

func CreateConfig() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	createDefaultConfigFile(cwd)
}

func createDir(cwd string, configPath string) {
	dirConfigPath := viper.GetString(configPath)

	dirPath := fmt.Sprintf("%s/%s", cwd, dirConfigPath)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func createADRDefaultTemplate(cwd string, templatesPath string) {
	t := viper.GetString(templatesPath)
	path := fmt.Sprintf("%s/%s", cwd, t)
	fileName := viper.GetString("templates.adr.default")

	fmt.Println(fileName)
	d := `# {{ .Title }}

Status: {{ .Status }}
Author: {{ .Author }}
Date: {{ .Date }}

## Context and Problem Statement

## Decision Drivers


## Considererd Options


## Decision Outcome
  `
	dd := []byte(d)
	// os.WriteFile(path+fileName, dd, 0644)
	fmt.Println(path + fileName)
	f, err := os.Create(path + fileName)
	check(err)
	defer f.Close()

	_, err = f.Write(dd)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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
    default: "adr.tmpl"`
	rc := []byte(rexConfig)
	fileName := cwd + "/" + rexConfigFile
	fmt.Println(fileName)
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()

	_, err = f.Write(rc)
	check(err)
}
