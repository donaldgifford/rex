package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func TestGenerateDirectories_Cmd(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		templatesEnabled bool
		templatesPath    string
		content          string
		setArgs          []string
		err              bool
	}{
		"default": {
			configPath:       "tests/dirs/docs/adr/",
			templatesEnabled: false,
			templatesPath:    "",
			content:          "",
			setArgs:          []string{"--config=tests/.dirs-rex.yaml", "config", "generate", "directories"},
			err:              false,
		},
		"templates_enabled": {
			configPath:       "tests/dirs/docs/adr/",
			templatesEnabled: true,
			templatesPath:    "tests/dirs/docs/templates/",
			content:          "",
			setArgs:          []string{"--config=tests/.dirs-enabled-rex.yaml", "config", "generate", "directories"},
			err:              false,
		},
		// TODO: somehow test an error from mkdirall
		// "default_error": {
		// 	configPath:       "tests/dir-errors/docs/adr/",
		// 	templatesEnabled: false,
		// 	templatesPath:    "",
		// 	content:          "",
		// 	setArgs:          []string{"--config=tests/.dirs-error-rex.yaml", "config", "generate", "directories"},
		// 	err:              true,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(test.setArgs)

			err := rootCmd.Execute()
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
				assert.Equal(t, true, directoryExists(test.configPath))

				if test.templatesEnabled {
					assert.Equal(t, true, directoryExists(test.templatesPath))
				}
			}
		})
	}
}
