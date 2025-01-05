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

	"github.com/donaldgifford/rex/internal/templates"
)

// fileExists returns checks if a file already exists on disk
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// writeTemplateFile writes a template to disk
func (r *RexConfig) writeTemplateFile(file []byte, templateType string) error {
	templateFile, err := os.Create(fmt.Sprintf("%s%s", r.Templates.Path, templateType))
	if err != nil {
		return err
	}

	_, err = templateFile.Write(file)
	if err != nil {
		return err
	}

	err = templateFile.Close()
	if err != nil {
		return err
	}
	return nil
}

// GenerateDefaultTemplates creates the default templates used for rex
//
// if force is set, it will overwrite the current template files if
// found with the defaults
func (r *RexConfig) GenerateDefaultTemplates(force bool) error {
	// create a templates setting to use
	eb := templates.EmbeddedTemplate{
		Settings: templates.Settings{
			TemplatePath:  "default/",
			AdrTemplate:   "adr.tmpl",
			IndexTemplate: "index.tmpl",
		},
	}

	// if force is true, overwrite current templates with the defaults
	if force {
		err := r.createDefaultTemplates(eb)
		if err != nil {
			return err
		}

		return nil
	}

	// check templates exist, if they do dont overwrite
	if fileExists(r.Settings().Templates.Path + eb.Settings.AdrTemplate) {
		return fmt.Errorf("ADR template file exists at: %s, please set --force to overwrite", r.Settings().Templates.Path+eb.Settings.AdrTemplate)
	}

	if fileExists(r.Settings().Templates.Path + eb.Settings.IndexTemplate) {
		return fmt.Errorf("index template file exists at: %s, please set --force to overwrite", r.Settings().Templates.Path+eb.Settings.IndexTemplate)
	}

	err := r.createDefaultTemplates(eb)
	if err != nil {
		return err
	}

	return nil
}

// createDefaultTemplates uses the template settings to create the default templates on disk
func (r *RexConfig) createDefaultTemplates(templateSettings templates.EmbeddedTemplate) error {
	// read default adr template from embedded template
	a, err := templateSettings.Read(templateSettings.Settings.TemplatePath + templateSettings.Settings.AdrTemplate)
	if err != nil {
		return err
	}

	// write adr template to file
	err = r.writeTemplateFile(a, templateSettings.Settings.AdrTemplate)
	if err != nil {
		return err
	}

	// read default Index template from embedded template
	t, err := templateSettings.Read(templateSettings.Settings.TemplatePath + templateSettings.Settings.IndexTemplate)
	if err != nil {
		return err
	}

	// write index template to file
	err = r.writeTemplateFile(t, templateSettings.Settings.IndexTemplate)
	if err != nil {
		return err
	}

	return nil
}
