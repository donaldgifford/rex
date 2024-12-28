package cmd

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// var (
// 	defaultAdrPath               string
// 	defaultAdrIndexPage          string
// 	defaultAdrAddToIndex         bool
// 	defaultTemplatesPath         string
// 	defaultTemplatesEnabled      bool
// 	defaultTemplatesAdrDefault   string
// 	defaultTemplatesAdrIndex     string
// 	defaultEnabledGithubPages    bool
// 	defaultPagesIndex            string
// 	defaultPagesWebConfig        string
// 	defaultPagesWebLayoutAdr     string
// 	defaultPagesWebLayoutDefault string
// 	defaultExtras                bool
// 	defaultExtraPagesInstall     string
// 	defaultExtraPagesUsage       string
// )
//
// func viperSetHelper() {
// 	viper.Set("adr.path", defaultAdrPath)
// 	viper.Set("adr.index_page", defaultAdrIndexPage)
// 	viper.Set("adr.add_to_index", defaultAdrAddToIndex)
// 	viper.Set("templates.enabled", defaultTemplatesEnabled)
// 	viper.Set("templates.path", defaultTemplatesPath)
// 	viper.Set("templates.adr.default", defaultTemplatesAdrDefault)
// 	viper.Set("templates.adr.index", defaultTemplatesAdrIndex)
// 	viper.Set("enable_github_pages", defaultEnabledGithubPages)
// 	viper.Set("pages.index", defaultPagesIndex)
// 	viper.Set("pages.web.config", defaultPagesWebConfig)
// 	viper.Set("pages.web.layout.adr", defaultPagesWebLayoutAdr)
// 	viper.Set("pages.web.layout.default", defaultPagesWebLayoutDefault)
// 	viper.Set("extras", defaultExtras)
// 	viper.Set("extra_pages.install", defaultExtraPagesInstall)
// 	viper.Set("extra_pages.usage", defaultExtraPagesUsage)
// }
//
// type ViperConfig struct {
// 	defaultAdrPath               string `mapstructure:"REX_ADR_PATH"`
// 	defaultAdrIndexPage          string `mapstructure:"REX_ADR_INDEX_PAGE"`
// 	defaultAdrAddToIndex         bool   `mapstructure:"REX_ADR_ADD_TO_INDEX"`
// 	defaultTemplatesPath         string `mapstructure:"REX_TEMPLATES_PATH"`
// 	defaultTemplatesEnabled      bool   `mapstructure:"REX_TEMPLATES_ENABLED"`
// 	defaultTemplatesAdrDefault   string `mapstructure:"REX_TEMPLATES_ADR_DEFAULT"`
// 	defaultTemplatesAdrIndex     string `mapstructure:"REX_TEMPLATES_ADR_INDEX"`
// 	defaultEnabledGithubPages    bool   `mapstructure:"REX_ENABLED_GITHUB_PAGES"`
// 	defaultPagesIndex            string `mapstructure:"REX_PAGES_INDEX"`
// 	defaultPagesWebConfig        string `mapstructure:"REX_PAGES_WEB_CONFIG"`
// 	defaultPagesWebLayoutAdr     string `mapstructure:"REX_PAGES_WEB_LAYOUT_ADR"`
// 	defaultPagesWebLayoutDefault string `mapstructure:"REX_PAGES_WEB_LAYOUT_DEFAULT"`
// 	defaultExtras                bool   `mapstructure:"REX_EXTRAS"`
// 	defaultExtraPagesInstall     string `mapstructure:"REX_EXTRA_PAGES_INSTALL"`
// 	defaultExtraPagesUsage       string `mapstructure:"REX_EXTRA_PAGES_USAGE"`
// }
//
// func setOsEnvs() {
// 	os.Setenv("REX_ADR_PATH", defaultAdrPath)
// 	os.Setenv("REX_ADR_INDEX_PAGE", defaultAdrIndexPage)
// 	os.Setenv("REX_ADR_ADD_TO_INDEX", strconv.FormatBool(defaultAdrAddToIndex))
// 	os.Setenv("REX_TEMPLATES_PATH", defaultTemplatesPath)
// 	os.Setenv("REX_TEMPLATES_ENABLED", strconv.FormatBool(defaultTemplatesEnabled))
// 	os.Setenv("REX_TEMPLATES_ADR_DEFAULT", defaultTemplatesAdrDefault)
// 	os.Setenv("REX_TEMPLATES_ADR_INDEX", defaultTemplatesAdrIndex)
// 	os.Setenv("REX_ENABLED_GITHUB_PAGES", strconv.FormatBool(defaultEnabledGithubPages))
// 	os.Setenv("REX_PAGES_INDEX", defaultPagesIndex)
// 	os.Setenv("REX_PAGES_WEB_CONFIG", defaultPagesWebConfig)
// 	os.Setenv("REX_PAGES_WEB_LAYOUT_ADR", defaultPagesWebLayoutAdr)
// 	os.Setenv("REX_PAGES_WEB_LAYOUT_DEFAULT", defaultPagesWebLayoutDefault)
// 	os.Setenv("REX_EXTRAS", strconv.FormatBool(defaultExtras))
// 	os.Setenv("REX_EXTRA_PAGES_INSTALL", defaultExtraPagesInstall)
// 	os.Setenv("REX_EXTRA_PAGES_USAGE", defaultExtraPagesUsage)
// }
//
// func initDefaultConfig() {
// 	defaultAdrPath = "tests/docs/adr/"
// 	defaultAdrIndexPage = "README.md"
// 	defaultAdrAddToIndex = true
// 	defaultTemplatesPath = "templates/"
// 	defaultTemplatesEnabled = false
// 	defaultTemplatesAdrDefault = "adr.tmpl"
// 	defaultTemplatesAdrIndex = "index.tmpl"
// 	defaultEnabledGithubPages = true
// 	defaultPagesIndex = "index.md"
// 	defaultPagesWebConfig = "_config.yml"
// 	defaultPagesWebLayoutAdr = "adr.html"
// 	defaultPagesWebLayoutDefault = "default.html"
// 	defaultExtras = true
// 	defaultExtraPagesInstall = "install.md"
// 	defaultExtraPagesUsage = "usage.md"
// }
//
// // LoadConfig reads configuration from file or environment variables.
// func LoadConfig(path string) (config ViperConfig, err error) {
// 	viper.AddConfigPath(path)
// 	viper.SetConfigName("app")
// 	viper.SetConfigType("env")
//
// 	viper.AutomaticEnv()
//
// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		return
// 	}
//
// 	err = viper.Unmarshal(&config)
// 	return
// }

func ReadTestFile(file string) ([]byte, error) {
	t, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// path: "tests/docs/adr/"

func createConfigFile(file string) error {
	rexConfig := `adr:
  path: "tests/docs/adr/"
  index_page: "README.md"
  add_to_index: true # on rex create, a new record will be added to the index page
templates:
  enabled: false
  path: "templates/"
  adr:
    default: "adr.tmpl"
    index: "index.tmpl"
enable_github_pages: false
pages:
  index: "index.md"
  web:
    config: "_config.yml"
    layout:
      adr: "adr.html"
      default: "default.html"
extras: true
extra_pages:
  install: install.md
  usage: usage.md"`
	rc := []byte(rexConfig)
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(rc)
	if err != nil {
		return err
	}
	return nil
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

	// err = createConfigFile("tests/.rex.yaml")
	// if err != nil {
	// 	log.Print(err)
	// 	os.Exit(1)
	// }

	code := m.Run()

	err = removeTestFolder("tests")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// err = removeTestConfigFile("tests/.rex.yaml")
	// if err != nil {
	// 	log.Print(err)
	// 	os.Exit(1)
	// }

	os.Exit(code)
}
