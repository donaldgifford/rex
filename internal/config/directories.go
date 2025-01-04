package config

import (
	"os"
)

// GenerateDirectories creates the default directories used for rex
// force is used to overwrite the templates if found
//
// if "templates.enabled: true" is the .rex.yaml config file
// then this function creates the default templates directory
func (r *RexConfig) GenerateDirectories() error {
	// create default templates if "templates.enabled"
	if r.Templates.Enabled {
		// create the templates directory from settings
		err := os.MkdirAll(r.Templates.Path, 0750)
		if err != nil && !os.IsExist(err) {
			return err
		}

		return nil
	}

	// mkdirall with path string
	err := os.MkdirAll(r.ADR.Path, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}
