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

	"github.com/donaldgifford/rex/internal/templates"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRexConfig_GenerateDefaultTemplates(t *testing.T) {
	tests := map[string]struct {
		configPath    string
		templatesPath string
		configIndex   string
		configAdd     bool
		cwd           string
		expected      []string
		force         bool
		err           bool
	}{
		"good_force": {
			configPath:    defaultAdrPath,
			templatesPath: defaultTemplatesPath,
			configIndex:   "README.md",
			configAdd:     true,
			cwd:           "",
			expected:      []string{"1-test1.md", "2-test2.md"},
			force:         true,
			err:           false,
		},

		"error_adr_found_no_force": {
			configPath:    defaultAdrPath,
			templatesPath: defaultTemplatesPath + "poop/adr/",
			configIndex:   "README.md",
			configAdd:     true,
			cwd:           "",
			expected:      []string{"1-test1.md", "2-test2.md"},
			force:         false,
			err:           true,
		},
		"error_index_found_no_force": {
			configPath:    defaultAdrPath,
			templatesPath: defaultTemplatesPath + "poop/index/",
			configIndex:   "README.md",
			configAdd:     true,
			cwd:           "",
			expected:      []string{"1-test1.md", "2-test2.md"},
			force:         false,
			err:           true,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.configPath)
		viper.Set("adr.index_page", test.configIndex)
		viper.Set("adr.add_to_index", test.configAdd)
		viper.Set("templates.path", test.templatesPath)

		r := NewRexConfig()

		err := r.GenerateDefaultTemplates(test.force)
		t.Run(name, func(t *testing.T) {
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestRexConfig_createDefaultTemplates(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		configIndex      string
		configAdd        bool
		templateSettings templates.EmbeddedTemplate
		cwd              string
		err              bool
	}{
		"good_force": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			templateSettings: templates.EmbeddedTemplate{
				Settings: templates.Settings{
					TemplatePath:  "default/",
					AdrTemplate:   "adr.tmpl",
					IndexTemplate: "index.tmpl",
				},
			},
			cwd: "",
			err: false,
		},
		"error_no_embedded_adr_template": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			templateSettings: templates.EmbeddedTemplate{
				Settings: templates.Settings{
					TemplatePath:  "default/",
					AdrTemplate:   "adr_fail.tmpl",
					IndexTemplate: "index.tmpl",
				},
			},
			cwd: "",
			err: true,
		},
		"error_no_embedded_index_template": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			templateSettings: templates.EmbeddedTemplate{
				Settings: templates.Settings{
					TemplatePath:  "default/",
					AdrTemplate:   "adr.tmpl",
					IndexTemplate: "index_fail.tmpl",
				},
			},
			cwd: "",
			err: true,
		},
		"error_embedded_template_path": {
			configPath:  defaultAdrPath,
			configIndex: "README.md",
			configAdd:   true,
			templateSettings: templates.EmbeddedTemplate{
				Settings: templates.Settings{
					TemplatePath:  "default_fail/",
					AdrTemplate:   "adr.tmpl",
					IndexTemplate: "index.tmpl",
				},
			},
			cwd: "",
			err: true,
		},
	}

	for name, test := range tests {
		viper.Set("adr.path", test.configPath)
		viper.Set("adr.index_page", test.configIndex)
		viper.Set("adr.add_to_index", test.configAdd)

		viper.Set("templates.path", defaultTemplatesPath)

		r := NewRexConfig()

		err := r.createDefaultTemplates(test.templateSettings)
		t.Run(name, func(t *testing.T) {
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func gotFile(path string, file string) string {
	f, err := os.ReadFile(path + file)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(f)
}

func TestRexConfig_writeTemplateFile(t *testing.T) {
	// get cwd
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed on Getwd, err: %v\n", err.Error())
	}

	tests := map[string]struct {
		configPath   string
		file         []byte
		templateType string
		cwd          string
		expected     string
		err          bool
	}{
		"good_adr": {
			configPath:   defaultTemplatesPath,
			file:         []byte("# {{ .Content.Title }}\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| {{ .Content.Status }} | {{ .Content.Author }} | {{ .Content.Date }} | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome\n"),
			templateType: "adr2.tmpl",
			cwd:          "",
			expected:     "# {{ .Content.Title }}\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| {{ .Content.Status }} | {{ .Content.Author }} | {{ .Content.Date }} | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome\n",
			err:          false,
		},
		"adr_empty_file": {
			configPath:   defaultTemplatesPath,
			file:         []byte(""),
			templateType: "adr3.tmpl",
			cwd:          "",
			expected:     "",
			err:          false,
		},

		"good_index": {
			configPath:   defaultTemplatesPath,
			file:         []byte("# {{ .Content.Title }}\n\n## ADRs\n\n| ID | Title | Link |\n| -- | ----- | ---- |\n{{- range .Content.Adrs }}\n| {{ .Id }} | {{ .Title }} | link |\n{{- end }}\n"),
			templateType: "index2.tmpl",
			cwd:          cwd,
			expected:     "# {{ .Content.Title }}\n\n## ADRs\n\n| ID | Title | Link |\n| -- | ----- | ---- |\n{{- range .Content.Adrs }}\n| {{ .Id }} | {{ .Title }} | link |\n{{- end }}\n",
			err:          false,
		},
	}

	for name, test := range tests {

		viper.Set("templates.path", defaultTemplatesPath)

		r := NewRexConfig()

		err := r.writeTemplateFile(test.file, test.templateType)
		t.Run(name, func(t *testing.T) {
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
				got := gotFile(test.configPath, test.templateType)
				assert.Equal(t, test.expected, got, "")
			}
		})
	}
}
