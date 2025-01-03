/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>

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
package install

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/donaldgifford/rex/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func defaultValueMap() map[string]any {
	return map[string]any{
		"defaultAdrPath":               "docs/adr/",
		"defaultAdrIndexPage":          "README.md",
		"defaultAdrAddToIndex":         true,
		"defaultTemplatesPath":         "templates/",
		"defaultTemplatesEnabled":      false,
		"defaultTemplatesAdrDefault":   "adr.tmpl",
		"defaultTemplatesAdrIndex":     "index.tmpl",
		"defaultEnabledGithubPages":    true,
		"defaultPagesIndex":            "index.md",
		"defaultPagesWebConfig":        "_config.yml",
		"defaultPagesWebLayoutAdr":     "adr.html",
		"defaultPagesWebLayoutDefault": "default.html",
		"defaultExtras":                true,
		"defaultExtraPagesInstall":     "install.md",
		"defaultExtraPagesUsage":       "usage.md",
	}
}

func removeTestConfigFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		return err
	}
	return nil
}

func readYamlFile(file string) (*config.RexConfig, error) {
	var r *config.RexConfig

	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return r, err
	}

	err = yaml.Unmarshal(yamlFile, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func TestMain(m *testing.M) {
	code := m.Run()

	err := removeTestConfigFile(".rex.yaml")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	os.Exit(code)
}

func TestCreateRexConfigFile(t *testing.T) {
	tests := map[string]struct {
		err           bool
		defaultValues map[string]any
	}{
		"good": {
			err:           false,
			defaultValues: defaultValueMap(),
		},
		// "embedded": {
		// 	settings: Settings{
		// 		TemplatePath:  "default/",
		// 		AdrTemplate:   "adr.tmpl",
		// 		IndexTemplate: "index.tmpl",
		// 	},
		// 	enabled: false,
		// 	err:     false,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := CreateRexConfigFile()
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}

			c, err := readYamlFile(".rex.yaml")
			if err != nil {
				log.Fatal(err)
			}

			// adr
			assert.Equal(t, test.defaultValues["defaultAdrPath"], c.ADR.Path, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultAdrIndexPage"], c.ADR.IndexPage, "ADR IndexPage")
			assert.Equal(t, test.defaultValues["defaultAdrAddToIndex"], c.ADR.AddToIndex, "ADR AddToIndex")
			// template
			assert.Equal(t, test.defaultValues["defaultTemplatesPath"], c.Templates.Path, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultTemplatesEnabled"], c.Templates.Enabled, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultTemplatesAdrDefault"], c.Templates.ADR.Default, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultTemplatesAdrIndex"], c.Templates.ADR.Index, "ADR Path")
			// pages
			assert.Equal(t, test.defaultValues["defaultEnabledGithubPages"], c.EnableGithubPages, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultPagesIndex"], c.Pages.Index, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultPagesWebConfig"], c.Pages.Web.Config, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultPagesWebLayoutDefault"], c.Pages.Web.Layout.Default, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultPagesWebLayoutAdr"], c.Pages.Web.Layout.ADR, "ADR Path")
			// extras
			assert.Equal(t, test.defaultValues["defaultExtras"], c.Extras, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultExtraPagesInstall"], c.ExtraPages.Install, "ADR Path")
			assert.Equal(t, test.defaultValues["defaultExtraPagesUsage"], c.ExtraPages.Usage, "ADR Path")
		})
	}
}
