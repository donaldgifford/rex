/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>

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
package templates

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/donaldgifford/rex/internal/adr"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestEmbeddedGetSettings(t *testing.T) {
	tests := map[string]struct {
		settings Settings
		err      bool
	}{
		"embedded": {
			settings: Settings{
				TemplatePath:  "default/",
				AdrTemplate:   "adr.tmpl",
				IndexTemplate: "index.tmpl",
			},
			err: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("templates.enabled", false)
			tmp := NewTemplate()
			a := tmp.GetSettings()
			assert.Equal(t, test.settings.AdrTemplate, a.AdrTemplate, "")
			assert.Equal(t, test.settings.TemplatePath, a.TemplatePath, "")
			assert.Equal(t, test.settings.IndexTemplate, a.IndexTemplate, "")
		})
	}
}

func TestEmbeddedRead(t *testing.T) {
	tests := map[string]struct {
		file     string
		contents string
		err      bool
	}{
		"adr.tmpl": {
			file:     "adr.tmpl",
			contents: "# {{ .Content.Title }}\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| {{ .Content.Status }} | {{ .Content.Author }} | {{ .Content.Date }} | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome\n",
			err:      false,
		},
		"index.tmpl": {
			file:     "index.tmpl",
			contents: "# {{ .Content.Title }}\n\n## ADRs\n\n| ID | Title | Link |\n| -- | ----- | ---- |\n{{- range .Content.Adrs }}\n| {{ .Id }} | {{ .Title }} | link |\n{{- end }}\n",
			err:      false,
		},
		"index_readme.tmpl": {
			file:     "index_readme.tmpl",
			contents: "# ADR List\n\n## ADRs\n\n| ID | Title | Link |\n| -- | ----- | ---- |\n",
			err:      false,
		},
		// "rex.yaml": {
		// 	file:     "rex.yaml",
		// 	contents: "adr:\n  path: \"docs/adr/\"\n  index_page: \"README.md\"\n  add_to_index: true # on rex create, a new record will be added to the index page\ntemplates:\n  enabled: false # uses embedded templates by default. If true reference the paths\n  path: \"templates/\"\n  adr:\n    default: \"adr.tmpl\"\n    index: \"index.tmpl\"\nenable_github_pages: true\npages:\n  index: \"index.md\"\n  web:\n    config: \"_config.yml\"\n    layout:\n      adr: \"adr.html\"\n      default: \"default.html\"\nextras: true\nextra_pages:\n  install: install.md\n  usage: usage.md\n",
		// 	err:      false,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("templates.enabled", false)
			tmp := NewTemplate()
			a, err := tmp.Read(fmt.Sprintf("%s%s", "default/", test.file))
			assert.Equal(t, test.contents, string(a), "")
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}
func TestEmbeddedExecute(t *testing.T) {}
func TestEmbeddedCreateADR(t *testing.T) {
	d := time.Now().Format(time.DateOnly)
	tests := map[string]struct {
		file    string
		content string
		adr     *adr.ADR
		err     bool
	}{
		"adr": {
			file:    "3-Test-3.md",
			content: "# Test 3\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| Draft | Author | 2024-12-29 | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome\n",
			adr: &adr.ADR{
				Content: adr.Content{
					Title:  "Test 3",
					Author: "Author",
					Status: "Draft",
					Date:   d,
				},
				ID: 3,
				Config: adr.ADRConfig{
					Path:       defaultAdrPath,
					IndexPage:  "README.md",
					AddToIndex: true,
				},
			},
			err: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("templates.enabled", false)

			tmp := NewTemplate()
			err := tmp.CreateADR(test.adr)
			if err != nil {
				t.Errorf("error creating test file: %v, err: %v", test.file, err.Error())
			}

			b, err := ReadTestFile(fmt.Sprintf("%s%s", defaultAdrPath, test.file))
			if err != nil {
				t.Errorf("error opening test file: %v, err: %v", test.file, err.Error())
			}
			assert.Equal(t, test.content, string(b), "")

			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func ReadTestFile(file string) ([]byte, error) {
	t, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func TestEmbeddedGenereateIndex(t *testing.T) {
	tests := map[string]struct {
		file    string
		content string
		idx     *adr.Index
		err     bool
	}{
		"idx": {
			file:    defaultTemplatesAdrIndex,
			content: "# ADR Index\n\n## ADRs\n\n| ID | Title | Link |\n| -- | ----- | ---- |\n| 1 | test1 | link |\n| 2 | test2 | link |\n| 3 | Test-3 | link |\n",
			idx: &adr.Index{
				DocPath:       defaultAdrPath,
				IndexFileName: defaultTemplatesAdrIndex,
				Content: adr.IndexContent{
					Title: "ADR Index",
					Adrs: []*adr.IndexAdr{
						{Id: 1, Title: "test1"},
						{Id: 2, Title: "test2"},
						{Id: 3, Title: "Test-3"},
					},
				},
			},
			err: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("templates.enabled", false)

			tmp := NewTemplate()
			err := tmp.GenerateIndex(test.idx)
			if err != nil {
				t.Errorf("error creating test file: %v, err: %v", test.file, err.Error())
			}

			b, err := ReadTestFile(fmt.Sprintf("%s%s", defaultAdrPath, test.file))
			if err != nil {
				t.Errorf("error opening test file: %v, err: %v", test.file, err.Error())
			}
			assert.Equal(t, test.content, string(b), "")

			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}
