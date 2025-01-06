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
package config

import (
	"fmt"
	"log"
	"os"
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

	err = createTestFolder(templatesPath + "poop/adr/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestADRFile(
		fmt.Sprintf("%s%s", templatesPath+"poop/adr/", "adr.tmpl"),
	)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestFolder(templatesPath + "poop/index/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestADRFile(
		fmt.Sprintf("%s%s", templatesPath+"poop/index/", "index.tmpl"),
	)
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

/*
* Acutal Tests
 */

func TestRexConfig_NewRexConfigure(t *testing.T) {
	tests := map[string]struct {
		cwd      string
		expected string
		err      bool
	}{
		"output": {
			cwd:      "",
			expected: "adr:\n    path: tests/docs/adr/\n    index_page: README.md\n    add_to_index: true\ntemplates:\n    enabled: false\n    path: tests/docs/templates/\n    adr:\n        default: adr.tmpl\n        index: index.tmpl\nenable_github_pages: true\npages:\n    index: index.md\n    web:\n        config: _config.yml\n        layout:\n            adr: adr.html\n            default: default.html\nextras: true\nextra_pages:\n    install: install.md\n    usage: usage.md\n",
			err:      false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewRexConfigure()

			b, err := r.YamlOut()
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
				assert.Equal(t, test.expected, string(b))

			}
		})
	}
}

func TestNewRexConfig(t *testing.T) {
	c := &RexConfig{
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

	config := NewRexConfig()

	if config.ADR != c.ADR {
		t.Errorf("ADR Settings dont match: %v, %v", config.ADR, c.ADR)
	}

	if config.Templates != c.Templates {
		t.Errorf(
			"Templates settings dont match: %v, %v",
			config.Templates,
			c.Templates,
		)
	}

	if config.Pages != c.Pages {
		t.Errorf("Pages settings dont match: %v, %v", config.Pages, c.Pages)
	}
}

func TestRexConfigSettings(t *testing.T) {
	c := &RexConfig{
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

	rc := NewRexConfig()
	config := rc.Settings()

	if config.ADR != c.ADR {
		t.Errorf("ADR Settings dont match: %v, %v", config.ADR, c.ADR)
	}

	if config.Templates != c.Templates {
		t.Errorf(
			"Templates settings dont match: %v, %v",
			config.Templates,
			c.Templates,
		)
	}

	if config.Pages != c.Pages {
		t.Errorf("Pages settings dont match: %v, %v", config.Pages, c.Pages)
	}
}

func TestRexConfig_YamlOut(t *testing.T) {
	tests := map[string]struct {
		cwd      string
		expected string
		err      bool
	}{
		"output": {
			cwd:      "",
			expected: "adr:\n    path: tests/docs/adr/\n    index_page: README.md\n    add_to_index: true\ntemplates:\n    enabled: false\n    path: tests/docs/templates/\n    adr:\n        default: adr.tmpl\n        index: index.tmpl\nenable_github_pages: true\npages:\n    index: index.md\n    web:\n        config: _config.yml\n        layout:\n            adr: adr.html\n            default: default.html\nextras: true\nextra_pages:\n    install: install.md\n    usage: usage.md\n",
			err:      false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewRexConfig()

			b, err := r.YamlOut()
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
				assert.Equal(t, test.expected, string(b))

			}
		})
	}
}
