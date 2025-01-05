package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTemplates_Cmd(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		templatesEnabled bool
		templatesPath    string
		content          string
		indexFile        string
		setArgs          []string
		err              bool
	}{
		"default_not_enabled": {
			configPath:       "tests/docs/adr/",
			templatesEnabled: false,
			templatesPath:    "tests/docs/templates/",
			content:          "",
			indexFile:        "README.md",
			setArgs:          []string{"--config=tests/.rex.yaml", "config", "generate", "templates"},
			err:              false,
		},
		"templates_enabled": {
			configPath:       "tests/dirs/docs/adr/",
			templatesEnabled: true,
			templatesPath:    "tests/dirs/docs/templates/",
			content:          "",
			indexFile:        "README.md",
			setArgs:          []string{"--config=tests/.dirs-enabled-rex.yaml", "config", "generate", "templates", "-f"},
			err:              false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(test.setArgs)

			rootCmd.Execute()
			if test.templatesEnabled {

				assert.Equal(t, true, fileExists(test.templatesPath+"adr.tmpl"))
				assert.Equal(t, true, fileExists(test.templatesPath+"index.tmpl"))
			}
		})
	}
}
