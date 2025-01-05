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
