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
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func TestRexConfig_GenereateDirectories(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		configIndex      string
		configAdd        bool
		templatesEnabled bool
		templatesPath    string
		cwd              string
		expected         []string
		err              bool
	}{
		"templates_disabled": {
			configPath:       defaultAdrPath + "gendir",
			configIndex:      "README.md",
			configAdd:        true,
			templatesEnabled: false,
			templatesPath:    "",
			cwd:              "",
			expected:         []string{"1-test1.md", "2-test2.md"},
			err:              false,
		},
		"templates_enabled": {
			configPath:       defaultAdrPath + "gendir2",
			configIndex:      "README.md",
			configAdd:        true,
			templatesEnabled: true,
			templatesPath:    defaultAdrPath + "gendir2" + "/templates",
			cwd:              "",
			expected:         []string{"1-test1.md", "2-test2.md"},
			err:              false,
		},
		"templates_disabled_error": {
			configPath:       defaultAdrPath + "README.md",
			configIndex:      "README.md",
			configAdd:        true,
			templatesEnabled: false,
			templatesPath:    "",
			cwd:              "",
			expected:         []string{"1-test1.md", "2-test2.md"},
			err:              true,
		},
		"templates_enabled_error": {
			configPath:       defaultAdrPath,
			configIndex:      "README.md",
			configAdd:        true,
			templatesEnabled: true,
			templatesPath:    defaultTemplatesPath + "poop/adr/adr.tmpl",
			cwd:              "",
			expected:         []string{"1-test1.md", "2-test2.md"},
			err:              true,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.configPath)
		viper.Set("adr.index_page", test.configIndex)
		viper.Set("adr.add_to_index", test.configAdd)
		viper.Set("templates.enabled", test.templatesEnabled)
		viper.Set("templates.path", test.templatesPath)
		t.Run(name, func(t *testing.T) {
			r := NewRexConfig()

			err := r.GenerateDirectories()
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
