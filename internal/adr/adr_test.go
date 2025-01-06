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
package adr

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// var adrDocsPath string

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
	defaultTemplatesPath = "templates/"
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

func TestMain(m *testing.M) {
	adrDocsPath := "tests/docs/adr/"

	err := createTestFolder(adrDocsPath)
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
	os.Exit(code)
}

/*
* Acutal Tests
 */

func TestNewADRConfig(t *testing.T) {
	// Test that NewADRConfig reads configuration settings correctly
	// viperSetHelper()
	config := NewADRConfig()
	if config.Path != defaultAdrPath ||
		config.IndexPage != defaultAdrIndexPage ||
		config.AddToIndex != defaultAdrAddToIndex {
		t.Errorf("NewADRConfig returned incorrect settings: %v", config)
	}
}

func TestNewADR(t *testing.T) {
	// test NewADR creates a valid config and ADR
	// viperSetHelper()
	adrT := NewADR()

	if adrT.Config.Path != defaultAdrPath ||
		adrT.Config.IndexPage != defaultAdrIndexPage ||
		adrT.Config.AddToIndex != defaultAdrAddToIndex {
		t.Errorf("NewADRConfig returned incorrect settings: %v", adrT.Config)
	}
}

// Tests ADRGetAdrFilesNames
func TestADRGetAdrFilesNames(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected []string
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			expected: []string{"1-test1.md", "2-test2.md"},
			err:      false,
		},
		"bad_path": {
			path:     "path/to/adrs",
			expected: []string(nil),
			err:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("adr.path", test.path)
			a := NewADR()
			actual, err := a.GetAdrFilesNames()
			assert.Equal(t, test.expected, actual, "")
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestNewADRCreate(t *testing.T) {
	d := time.Now().Format(time.DateOnly)
	tests := map[string]struct {
		configPath  string
		configIndex string
		configAdd   bool
		content     Content
		id          int
		err         bool
	}{
		"good": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			content: Content{
				Title:  "Test 3",
				Author: "Author",
				Status: "Draft",
				Date:   d,
			},
			id:  3,
			err: false,
		},
		"error": {
			configPath:  "/path/to/adr",
			configIndex: "README.md",
			configAdd:   true,
			content: Content{
				Title:  "",
				Author: "",
				Status: "",
				Date:   "",
			},
			id:  3,
			err: true,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.configPath)
		viper.Set("adr.index_page", test.configIndex)
		viper.Set("adr.add_to_index", test.configAdd)

		a := NewADR()

		// create content
		u := Content{
			Title:  "Test 3",
			Author: "Author",
			Status: "Draft",
			Date:   d,
		}

		c, err := a.Create(&u)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.content.Title, c.Content.Title, "")
			assert.Equal(t, test.content.Author, c.Content.Author, "")
			assert.Equal(t, test.content.Status, c.Content.Status, "")
			assert.Equal(t, test.content.Date, c.Content.Date, "")
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestADRGetSettings(t *testing.T) {
	tests := map[string]struct {
		path  string
		index string
		add   bool
	}{
		"good": {
			path:  "/path/to/adr",
			index: "index.md",
			add:   true,
		},
		"good2": {
			path:  "/path/to/adr",
			index: "README.md",
			add:   false,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.path)
		viper.Set("adr.index_page", test.index)
		viper.Set("adr.add_to_index", test.add)

		config := NewADRConfig()
		settings := config.GetSettings()
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.path, settings.Path, "")
			assert.Equal(t, test.index, settings.IndexPage, "")
			assert.Equal(t, test.add, settings.AddToIndex, "")
		})
	}
}

func TestId(t *testing.T) {
	tests := map[string]struct {
		configPath  string
		configIndex string
		configAdd   bool
		id          int
		err         bool
	}{
		"good": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			id:          3,
			err:         false,
		},
		"error": {
			configPath:  "/path/to/adr",
			configIndex: "index.md",
			configAdd:   true,
			id:          0,
			err:         true,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.configPath)
		viper.Set("adr.index_page", test.configIndex)
		viper.Set("adr.add_to_index", test.configAdd)

		a := NewADR()
		id, err := a.Id()
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.id, id, "")
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestNewIADR(t *testing.T) {
	tests := map[string]struct {
		path  string
		index string
		add   bool
	}{
		"good": {
			path:  defaultAdrPath,
			index: "README.md",
			add:   true,
		},
		"error": {
			path:  "/path/to/adr",
			index: "index.md",
			add:   true,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.path)
		viper.Set("adr.index_page", test.index)
		viper.Set("adr.add_to_index", test.add)

		a := NewIADR()
		settings := a.GetSettings()
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.path, settings.Path, "")
			assert.Equal(t, test.index, settings.IndexPage, "")
			assert.Equal(t, test.add, settings.AddToIndex, "")
		})
	}
}
