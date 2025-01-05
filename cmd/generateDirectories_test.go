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
