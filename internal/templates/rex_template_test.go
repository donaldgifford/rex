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
	"testing"
	"time"

	"github.com/donaldgifford/rex/internal/adr"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRexTemplate_Read(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		file    string
		want    string
		wantErr bool
	}{
		{
			name: "adr_tmpl",
			file: "tests/docs/templates/adr.tmpl",
			want: "# {{ .Content.Title }}\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| {{ .Content.Status }} | {{ .Content.Author }} | {{ .Content.Date }} | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome",
			// want:    []byte{},
			wantErr: false,
		},
		{
			name: "index_tmpl",
			file: "tests/docs/templates/index.tmpl",
			want: "# {{ .Content.Title }}\n\n## ADRs\n\n| ID | Title | Link |\n| -- | ----- | ---- |\n{{- range .Content.Adrs }}\n| {{ .Id }} | {{ .Title }} | link |\n{{- end }}",
			// want:    []byte{},
			wantErr: false,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			var rt RexTemplate
			got, gotErr := rt.Read(tt.file)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Read() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Read() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				assert.Equal(t, tt.want, string(got), "")
				// t.Errorf("Read() = %v, want %v", got, string(tt.want))
			}
		})
	}
}

func TestRexTemplate_GetSettings(t *testing.T) {
	tests := map[string]struct {
		settings Settings
		err      bool
	}{
		"template": {
			settings: Settings{
				TemplatePath:  "tests/docs/templates/",
				AdrTemplate:   "adr.tmpl",
				IndexTemplate: "index.tmpl",
			},
			err: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("templates.enabled", true)
			viper.Set("templates.path", defaultTemplatesPath)
			tmp := NewTemplate()
			a := tmp.GetSettings()
			assert.Equal(t, test.settings.AdrTemplate, a.AdrTemplate, "")
			assert.Equal(t, test.settings.TemplatePath, a.TemplatePath, "")
			assert.Equal(t, test.settings.IndexTemplate, a.IndexTemplate, "")
		})
	}
}
func TestRexTemplate_Execute(t *testing.T) {}

func TestRexTemplate_CreateADR(t *testing.T) {
	d := time.Now().Format(time.DateOnly)
	tests := map[string]struct {
		file    string
		content string
		adr     *adr.ADR
		err     bool
	}{
		"adr": {
			file:    "3-Test-3.md",
			content: parseContentWithDate("# Test 3\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| Draft | Author | %s | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome"),
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
			viper.Set("templates.enabled", true)
			viper.Set("templates.path", defaultTemplatesPath)

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

func TestRexTemplate_GenerateIndex(t *testing.T) {
	tests := map[string]struct {
		file    string
		content string
		idx     *adr.Index
		force   bool
		err     bool
	}{
		"create": {
			file:    defaultTemplatesAdrIndex,
			content: "# ADR Index\n\n## ADRs\n\n| ID | Title | Link |\n| -- | ----- | ---- |\n| 1 | test1 | link |\n| 2 | test2 | link |\n| 3 | Test-3 | link |\n",
			idx: &adr.Index{
				DocPath:       defaultAdrPath,
				IndexFileName: "rex_" + defaultAdrIndexPage,
				Content: adr.IndexContent{
					Title: "ADR Index",
					Adrs: []*adr.IndexAdr{
						{Id: 1, Title: "test1"},
						{Id: 2, Title: "test2"},
						{Id: 3, Title: "Test-3"},
					},
				},
			},
			force: false,
			err:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			viper.Set("templates.enabled", true)
			viper.Set("templates.path", defaultTemplatesPath)
			viper.Set("adr.index_page", "rex_"+defaultAdrIndexPage)

			tmp := NewTemplate()
			err := tmp.GenerateIndex(test.idx, test.force)
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
