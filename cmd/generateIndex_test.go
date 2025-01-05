package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fileExists returns checks if a file already exists on disk
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func TestGenerateIndex_Cmd(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		templatesEnabled bool
		templatesPath    string
		content          string
		indexFile        string
		setArgs          []string
		err              bool
	}{
		"default_embedded": {
			configPath:       "tests/docs/adr/",
			templatesEnabled: false,
			templatesPath:    "",
			content:          "",
			indexFile:        "README.md",
			setArgs:          []string{"--config=tests/.rex.yaml", "config", "generate", "index"},
			err:              false,
		},
		"templates_enabled": {
			configPath:       "tests/dirs/docs/adr/",
			templatesEnabled: true,
			templatesPath:    "tests/dirs/docs/templates/",
			content:          "",
			indexFile:        "README.md",
			setArgs:          []string{"--config=tests/.dirs-enabled-rex.yaml", "config", "generate", "index"},
			err:              false,
		},
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

				fmt.Println(test.configPath + test.indexFile)
				assert.Equal(t, true, fileExists(test.configPath+test.indexFile))
				//
				// if test.templatesEnabled {
				// 	assert.Equal(t, true, directoryExists(test.templatesPath))
				// }
			}
		})
	}
}
