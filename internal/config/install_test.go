package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
* Acutal Tests
 */

func TestNewRexConfInstall(t *testing.T) {
	tests := map[string]struct {
		path     string
		file     string
		force    bool
		expected *RexConfInstall
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			file:     "default/rex.yaml",
			force:    false,
			expected: &RexConfInstall{},
			err:      false,
		},
		"bad_path": {
			path:     "path/to/adrs",
			file:     "default/rex-error.yaml",
			force:    true,
			expected: &RexConfInstall{},
			err:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := NewRexConfigInstall()
			assert.Equal(t, test.expected, a, "")
			// if test.err {
			// 	assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			// } else {
			// 	assert.Nil(t, err, "")
			// }
		})
	}
}

func TestRexConfInstallConfigExists(t *testing.T) {}
func TestRexConfInstallReadConfig(t *testing.T) {
	tests := map[string]struct {
		path     string
		file     string
		force    bool
		expected *RexConfInstall
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			file:     "default/rex.yaml",
			force:    false,
			expected: &RexConfInstall{},
			err:      false,
		},
		// "bad_path": {
		// 	path:     "path/to/adrs",
		// 	file:     "default/rex-error.yaml",
		// 	force:    true,
		// 	expected: &RexConfInstall{},
		// 	err:      true,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := NewRexConfigInstall()
			err := a.ReadConfig()

			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestRexConfInstallWriteConfig(t *testing.T) {
	tests := map[string]struct {
		path     string
		file     string
		force    bool
		expected []string
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			file:     "default/rex.yaml",
			force:    false,
			expected: []string{"1-test1.md", "2-test2.md"},
			err:      false,
		},
		"bad_path": {
			path:     "path/to/adrs",
			file:     "default/rex-error.yaml",
			force:    true,
			expected: []string(nil),
			err:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := NewRexConfigInstall()
			err := a.WriteConfig(test.file)
			// assert.Equal(t, test.expected, actual, "")
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestRexConfInstallGenereateConfig(t *testing.T) {
	tests := map[string]struct {
		path     string
		force    bool
		expected []string
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			force:    false,
			expected: []string{"1-test1.md", "2-test2.md"},
			err:      false,
		},
		// "bad_path": {
		// 	path:     "path/to/adrs",
		// 	force:    true,
		// 	expected: []string(nil),
		// 	err:      true,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := NewRexConfigInstall()
			err := a.GenerateConfig(test.force)
			// assert.Equal(t, test.expected, actual, "")
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestRexConfInstallGenereateIndex(t *testing.T) {
	tests := map[string]struct {
		path     string
		force    bool
		index    bool
		expected []string
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			force:    true,
			index:    true,
			expected: []string{"1-test1.md", "2-test2.md"},
			err:      false,
		},
		"bad_path": {
			path:     "path/to/adrs",
			force:    false,
			index:    true,
			expected: []string(nil),
			err:      true,
		},
	}

	for name, test := range tests {
		// removeTestConfigFile("docs/adr/README.md")
		createTestFolder("docs/adr/")
		t.Run(name, func(t *testing.T) {
			a := NewRexConfigInstall()
			err := a.GenerateIndex(test.force)
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestRexConfInstallGenereateDirectories(t *testing.T) {
	tests := map[string]struct {
		path     string
		force    bool
		index    bool
		expected []string
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			force:    true,
			index:    true,
			expected: []string{"1-test1.md", "2-test2.md"},
			err:      false,
		},
		// "bad_path": {
		// 	path:     "path/to/adrs",
		// 	force:    true,
		// 	index:    true,
		// 	expected: []string(nil),
		// 	err:      true,
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := NewRexConfigInstall()
			err := a.GenerateDirectories(test.force, test.index)
			if test.err {
				assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			} else {
				assert.Nil(t, err, "")
			}
		})
	}
}

func TestRexConfInstallSettings(t *testing.T) {
	tests := map[string]struct {
		path     string
		file     string
		force    bool
		expected *RexConfig
		err      bool
	}{
		"good": {
			path:     defaultAdrPath,
			file:     "default/rex.yaml",
			force:    false,
			expected: nil,
			err:      false,
		},
		"bad_path": {
			path:     "path/to/adrs",
			file:     "default/rex-error.yaml",
			force:    true,
			expected: nil,
			err:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := NewRexConfigInstall()
			err := a.Settings()

			assert.Equal(t, test.expected, err, "")

			// if test.err {
			// 	assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			// } else {
			// 	assert.Nil(t, err, "")
			// }
		})
	}
}
