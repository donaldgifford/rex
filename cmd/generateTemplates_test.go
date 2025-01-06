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
			setArgs: []string{
				"--config=tests/.rex.yaml",
				"config",
				"generate",
				"templates",
			},
			err: false,
		},
		"templates_enabled": {
			configPath:       "tests/dirs/docs/adr/",
			templatesEnabled: true,
			templatesPath:    "tests/dirs/docs/templates/",
			content:          "",
			indexFile:        "README.md",
			setArgs: []string{
				"--config=tests/.dirs-enabled-rex.yaml",
				"config",
				"generate",
				"templates",
				"-f",
			},
			err: false,
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
				assert.Equal(
					t,
					true,
					fileExists(test.templatesPath+"index.tmpl"),
				)
			}
		})
	}
}
