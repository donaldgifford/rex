/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/
package src

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

type ADR struct {
	Title  string
	Author string
	Status string
	Date   string
}

func CreateADR(adr ADR) {
	// docs/adr/
	adrPath := viper.GetString("adr.path")

	templatesPath := viper.GetString("templates.path")
	adrTemplate := viper.GetString("templates.adr.default")

	// tmplFile := "templates/adr.tmpl"
	tmplFile := fmt.Sprintf("%s/%s", templatesPath, adrTemplate)
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		log.Fatal(err)
	}

	id := getNewADRId()
	strippedTitle := strings.Join(strings.Split(strings.Trim(adr.Title, "\n \t"), " "), "-")
	fileName := fmt.Sprintf("%s-%s.md", id, strippedTitle)

	var f *os.File
	f, err = os.Create(adrPath + fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(f, adr)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func getADRFiles(path string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(path)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func getNewADRId() string {
	adrPath := viper.GetString("adr.path")
	files, err := getADRFiles(adrPath)
	if err != nil {
		log.Fatal(err)
	}
	var adrs []string
	for _, file := range files {
		if file == "README.md" {
		} else {
			adrs = append(adrs, file)
		}
	}

	if len(adrs) == 0 {
		return "0001"
	}

	var fileNames []int
	for _, v := range adrs {
		k := strings.Split(v, "-")[0]
		a, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal(err)
		}
		fileNames = append(fileNames, a)
	}

	lastestID := slices.Max(fileNames)
	return fmt.Sprintf("%04d", lastestID+1)
}
