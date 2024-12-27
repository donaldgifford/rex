package rex

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/donaldgifford/rex/internal/adr"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// var adrDocsPath string

// func viperSetHelper(path string, indexPage string, addToIndex bool) {
// 	viper.Set("adr.path", path)
// 	viper.Set("adr.index_page", indexPage)
// 	viper.Set("adr.add_to_index", addToIndex)
// }

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

func TestNewRex(t *testing.T) {
	tests := map[string]struct {
		path     string
		index    string
		add      bool
		install  bool
		expected string
		err      bool
	}{
		"good": {
			path:     "tests/docs/adr/",
			expected: "tests/docs/adr/",
			install:  false,
			err:      false,
		},
		"bad_path": {
			path:     "",
			expected: "",
			err:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("adr.path", test.path)
			r := NewRex(test.install)
			actual := r.Settings().ADR.Path
			assert.Equal(t, test.expected, actual, "")
		})
	}
}

func TestRexSettings(t *testing.T) {
	tests := map[string]struct {
		path     string
		index    string
		add      bool
		install  bool
		expected string
		err      bool
	}{
		"good": {
			path:     "tests/docs/adr/",
			expected: "tests/docs/adr/",
			install:  false,
			err:      false,
		},
		"bad_path": {
			path:     "",
			expected: "",
			install:  false,
			err:      true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("adr.path", test.path)
			r := NewRex(test.install)
			actual := r.Settings().ADR.Path
			assert.Equal(t, test.expected, actual, "")
		})
	}
}

func TestRexNewAdr(t *testing.T) {
	// validates an adr create call can happen without err and
	// with errors being handled. Eventually will have to get something
	// to read in an ADR file and parse correctness.
	d := time.Now().Format(time.DateOnly)
	tests := map[string]struct {
		configPath  string
		configIndex string
		configAdd   bool
		content     adr.Content
		id          int
		err         bool
	}{
		"good": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			content: adr.Content{
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
			content: adr.Content{
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

		// create content
		u := adr.Content{
			Title:  "Test 3",
			Author: "Author",
			Status: "Draft",
			Date:   d,
		}

		r := NewRex(false)
		err := r.NewADR(&u)

		t.Run(name, func(t *testing.T) {
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestRexUpdateIndex(t *testing.T) {
	// not really a good test on most these, need to setup mocks
	tests := map[string]struct {
		configPath  string
		configIndex string
		configAdd   bool
		err         bool
	}{
		"good": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			err:         false,
		},
		"error": {
			configPath:  "/path/to/adr",
			configIndex: "README.md",
			configAdd:   true,
			err:         true,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.configPath)
		viper.Set("adr.index_page", test.configIndex)
		viper.Set("adr.add_to_index", test.configAdd)

		r := NewRex(false)
		err := r.UpdateIndex()

		t.Run(name, func(t *testing.T) {
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}
