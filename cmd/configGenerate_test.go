package cmd

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGenerateCmdInstall(t *testing.T) {
	tests := map[string]struct {
		file    string
		content string
		setArgs []string
		err     bool
	}{
		"cmd": {
			file:    ".rex.yaml",
			content: "adr:\n  path: \"docs/adr/\"\n  index_page: \"README.md\"\n  add_to_index: true # on rex create, a new record will be added to the index page\ntemplates:\n  enabled: false # uses embedded templates by default. If true reference the paths\n  path: \"templates/\"\n  adr:\n    default: \"adr.tmpl\"\n    index: \"index.tmpl\"\nenable_github_pages: true\npages:\n  index: \"index.md\"\n  web:\n    config: \"_config.yml\"\n    layout:\n      adr: \"adr.html\"\n      default: \"default.html\"\nextras: true\nextra_pages:\n  install: install.md\n  usage: usage.md\n",
			// setArgs: []string{"--config=", "--force", "--index", "--install", "--directories"},
			// setArgs: []string{"--config=", "config", "generate", "--install"},
			setArgs: []string{"--install", "config", "generate", "f"},
			err:     false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			// configGenerateCmd.SetOut(buf)
			rootCmd.SetArgs(test.setArgs)

			// configGenerateCmd.SetArgs(test.setArgs)

			err := rootCmd.Execute()
			// err := configGenerateCmd.Execute()
			if err != nil {
				t.Errorf("error executing cmd: %v", err.Error())
			}
			b, err := ReadTestFile(test.file)
			if err != nil {
				t.Errorf("error opening test file: %v, err: %v", test.file, err.Error())
			}
			assert.Equal(t, test.content, string(b), "")

			// if test.err {
			// 	assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			// } else {
			// 	assert.Nil(t, err, "")
			// }
		})
	}

	err := removeTestConfigFile(".rex.yaml")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
func TestConfigGenerateCmd_Config_Force(t *testing.T) {}
func TestConfigGenerateCmd_Dirs_Index(t *testing.T)   {}
func TestConfigGenerateCmd_Index(t *testing.T)        {}
