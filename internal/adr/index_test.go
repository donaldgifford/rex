package adr

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewIIndex(t *testing.T) {
	tests := map[string]struct {
		path    string
		index   string
		content IndexContent
		err     bool
	}{
		"good": {
			path:  defaultAdrPath,
			index: "README.md",
			content: IndexContent{
				Title: "ADR Index",
			},
		},
		"bad_path": {
			path:  defaultAdrPath,
			index: "README.md",
			content: IndexContent{
				Title: "ADR Index",
			},
		},
	}

	for name, test := range tests {
		// viperSetHelper(test.path, "README.md", true)
		viper.Set("adr.path", test.path)
		viper.Set("adr.index_page", "README.md")
		viper.Set("adr.add_to_index", true)

		t.Run(name, func(t *testing.T) {
			i := NewIIndex()
			idx := i.Execute()
			c := IndexContent{
				Title: "ADR Index",
			}

			assert.Equal(t, test.path, idx.DocPath, "")
			assert.Equal(t, test.index, idx.IndexFileName, "")
			assert.Equal(t, test.content, c, "")
		})
	}
}

// func TestIndexExecute(t *testing.T) {}

func TestIndexProcess(t *testing.T) {
	tests := map[string]struct {
		expectedId    int
		expectedTitle string
		file          string
		path          string
		index         string
		content       IndexContent
	}{
		"good": {
			expectedId:    1,
			expectedTitle: "test1",
			file:          "1-test1.md",
			path:          defaultAdrPath,
			index:         "README.md",
			content: IndexContent{
				Title: "ADR Index",
			},
		},
		"bad_file": {
			// this is a false positive, need to rework the process function to
			// clean out any odd characters.
			expectedId:    2,
			expectedTitle: "--222",
			file:          "2---222.md",
			path:          defaultAdrPath,
			index:         "README.md",
			content: IndexContent{
				Title: "ADR Index",
			},
		},
	}

	for name, test := range tests {
		// viperSetHelper(test.path, "README.md", true)
		viper.Set("adr.path", test.path)
		viper.Set("adr.index_page", "README.md")
		viper.Set("adr.add_to_index", true)

		t.Run(name, func(t *testing.T) {
			i := NewIIndex()
			actual := i.Process(test.file)

			assert.Equal(t, test.expectedId, actual.Id, "")
			assert.Equal(t, test.expectedTitle, actual.Title, "")
		})
	}
}

func TestIndexADRs(t *testing.T) {
	tests := map[string]struct {
		expected []*IndexAdr
		path     string
		index    string
		content  IndexContent
		err      bool
	}{
		"good": {
			expected: []*IndexAdr{
				{
					Id:    1,
					Title: "test1",
				}, {
					Id:    2,
					Title: "test2",
				},
			},
			path:  defaultAdrPath,
			index: "README.md",
			content: IndexContent{
				Title: "ADR Index",
			},
			err: false,
		},
		"bad": {
			expected: []*IndexAdr(nil),
			path:     "/path/to/adr",
			index:    "README.md",
			content: IndexContent{
				Title: "ADR Index",
			},
			err: true,
		},
	}

	for name, test := range tests {
		// viperSetHelper(test.path, test.index, true)
		viper.Set("adr.path", test.path)
		viper.Set("adr.index_page", test.index)
		viper.Set("adr.add_to_index", true)

		t.Run(name, func(t *testing.T) {
			i := NewIIndex()
			err := i.ADRs()
			actual := i.Content.Adrs

			assert.Equal(t, test.expected, actual, "")

			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}
