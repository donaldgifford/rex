/*
Copyright © 2024-2025 Donald Gifford <dgifford06@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

var (
	defaultAdrPath               string
	defaultAdrIndexPage          string
	defaultAdrAddToIndex         bool
	defaultTemplatesPath         string
	defaultTemplatesEnabled      bool
	defaultTemplatesAdrDefault   string
	defaultTemplatesAdrIndex     string
	defaultEnabledGithubPages    bool
	defaultPagesIndex            string
	defaultPagesWebConfig        string
	defaultPagesWebLayoutAdr     string
	defaultPagesWebLayoutDefault string
	defaultExtras                bool
	defaultExtraPagesInstall     string
	defaultExtraPagesUsage       string
)

func viperSetHelper() {
	viper.Set("adr.path", defaultAdrPath)
	viper.Set("adr.index_page", defaultAdrIndexPage)
	viper.Set("adr.add_to_index", defaultAdrAddToIndex)
	viper.Set("templates.enabled", defaultTemplatesEnabled)
	viper.Set("templates.path", defaultTemplatesPath)
	viper.Set("templates.adr.default", defaultTemplatesAdrDefault)
	viper.Set("templates.adr.index", defaultTemplatesAdrIndex)
	viper.Set("enable_github_pages", defaultEnabledGithubPages)
	viper.Set("pages.index", defaultPagesIndex)
	viper.Set("pages.web.config", defaultPagesWebConfig)
	viper.Set("pages.web.layout.adr", defaultPagesWebLayoutAdr)
	viper.Set("pages.web.layout.default", defaultPagesWebLayoutDefault)
	viper.Set("extras", defaultExtras)
	viper.Set("extra_pages.install", defaultExtraPagesInstall)
	viper.Set("extra_pages.usage", defaultExtraPagesUsage)
}

func initDefaultConfig() {
	defaultAdrPath = "tests/docs/adr/"
	defaultAdrIndexPage = "README.md"
	defaultAdrAddToIndex = true
	defaultTemplatesPath = "tests/docs/templates/"
	defaultTemplatesEnabled = false
	defaultTemplatesAdrDefault = "adr.tmpl"
	defaultTemplatesAdrIndex = "index.tmpl"
	defaultEnabledGithubPages = true
	defaultPagesIndex = "index.md"
	defaultPagesWebConfig = "_config.yml"
	defaultPagesWebLayoutAdr = "adr.html"
	defaultPagesWebLayoutDefault = "default.html"
	defaultExtras = true
	defaultExtraPagesInstall = "install.md"
	defaultExtraPagesUsage = "usage.md"
}

func createTestADRFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func createTestTemplates(path string) error {
	cleanFile := filepath.Clean(path + "adr.tmpl")
	a, err := os.Create(cleanFile)
	if err != nil {
		return err
	}
	_, err = a.WriteString(`# {{ .Content.Title }}

| Status | Author         |  Created | Last Update | Current Version |
| ------ | -------------- | -------- | ----------- | --------------- |
| {{ .Content.Status }} | {{ .Content.Author }} | {{ .Content.Date }} | N/A | v0.0.1 |

## Context and Problem Statement

## Decision Drivers

## Considered Options

## Decision Outcome`)
	if err != nil {
		return err
	}
	err = a.Close()
	if err != nil {
		return err
	}

	cleanFile = filepath.Clean(path + "index.tmpl")
	i, err := os.Create(cleanFile)
	if err != nil {
		return err
	}
	_, err = i.WriteString(`# {{ .Content.Title }}

## ADRs

| ID | Title | Link |
| -- | ----- | ---- |
{{- range .Content.Adrs }}
| {{ .Id }} | {{ .Title }} | link |
{{- end }}`)
	if err != nil {
		return err
	}
	err = i.Close()
	if err != nil {
		return err
	}

	return nil
}

func createTestFolder(name string) error {
	err := os.MkdirAll(name, 0755)
	if err != nil {
		return err
	}
	return nil
}

func removeTestFolder(name string) error {
	err := os.RemoveAll(name)
	if err != nil {
		return err
	}
	return nil
}

func removeTestConfigFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	adrDocsPath := "tests/docs/adr/"
	templatesPath := "tests/docs/templates/"

	err := createTestFolder(adrDocsPath)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestFolder(templatesPath)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestTemplates(templatesPath)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// setup some test files
	err = createTestADRFile(fmt.Sprintf("%s%s", adrDocsPath, "1-test1.md"))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestADRFile(fmt.Sprintf("%s%s", adrDocsPath, "2-test2.md"))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestADRFile(fmt.Sprintf("%s%s", adrDocsPath, "README.md"))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	initDefaultConfig()
	viperSetHelper()

	code := m.Run()

	err = removeTestFolder("tests")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = removeTestFolder("docs")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	os.Exit(code)
}

func TestNewTemplate(t *testing.T) {
	tests := map[string]struct {
		settings Settings
		enabled  bool
		err      bool
	}{
		"rex": {
			settings: Settings{
				TemplatePath:  defaultTemplatesPath,
				AdrTemplate:   defaultTemplatesAdrDefault,
				IndexTemplate: defaultTemplatesAdrIndex,
			},
			enabled: true,
			err:     false,
		},
		"embedded": {
			settings: Settings{
				TemplatePath:  "default/",
				AdrTemplate:   "adr.tmpl",
				IndexTemplate: "index.tmpl",
			},
			enabled: false,
			err:     false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.enabled {
				viper.Set("templates.enabled", true)
			} else {
				viper.Set("templates.enabled", false)
			}
			a := NewTemplate()
			assert.Equal(t, test.settings.AdrTemplate, a.GetSettings().AdrTemplate, "")
			assert.Equal(t, test.settings.TemplatePath, a.GetSettings().TemplatePath, "")
			assert.Equal(t, test.settings.IndexTemplate, a.GetSettings().IndexTemplate, "")
		})
	}
}
