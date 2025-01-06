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
package templates

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/viper"

	"github.com/donaldgifford/rex/internal/adr"
)

// create an embedded file system to hold all the default config files

//go:embed default/adr.tmpl
//go:embed default/index.tmpl
//go:embed default/index_readme.tmpl
var DefaultRexTemplates embed.FS

// EmbeddedTemplate holds the Settings data
type EmbeddedTemplate struct {
	Settings Settings
}

// GetSettings returns the settings for ADR
func (et *EmbeddedTemplate) GetSettings() *Settings {
	return &et.Settings
}

// Read will read a template from the embedded fs and return it
func (et *EmbeddedTemplate) Read(file string) ([]byte, error) {
	t, err := DefaultRexTemplates.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Execute not implemented
func (et *EmbeddedTemplate) Execute() {}

// CreateADR creates adr files using the default embedded template
func (et *EmbeddedTemplate) CreateADR(adr *adr.ADR) error {
	// get the default template from settings and parse it
	tmpl, err := template.ParseFS(
		DefaultRexTemplates,
		fmt.Sprintf("%s%s", et.Settings.TemplatePath, et.Settings.AdrTemplate),
	)
	if err != nil {
		return err
	}

	// strip the ADR content to create a file name to use
	strippedTitle := strings.Join(
		strings.Split(strings.Trim(adr.Content.Title, "\n \t"), " "),
		"-",
	)
	fileName := fmt.Sprintf("%d-%s.md", adr.ID, strippedTitle)

	// create file on disk
	var f *os.File
	f, err = os.Create(viper.GetString("adr.path") + fileName)
	if err != nil {
		log.Fatal(err)
	}

	// execute template with ADR content
	err = tmpl.Execute(f, adr)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// GenerateIndex creates the index of adrs using the embedded index template
//
// force: if an index already exists, this option will overwrite it.
func (et *EmbeddedTemplate) GenerateIndex(idx *adr.Index, force bool) error {
	if force {
		err := et.writeIndex(idx)
		if err != nil {
			return err
		}

		return nil
	}

	// check index exists
	if fileExists(idx.DocPath + idx.IndexFileName) {
		return fmt.Errorf(
			"index file found at %s, to overwrite please pass --force flag",
			idx.DocPath+idx.IndexFileName,
		)
	}

	err := et.writeIndex(idx)
	if err != nil {
		return err
	}

	return nil
}

// writeIndex writes the index to disk using the default embedded template
func (et *EmbeddedTemplate) writeIndex(idx *adr.Index) error {
	// parse template from Settings
	tmpl, err := template.ParseFS(
		DefaultRexTemplates,
		fmt.Sprintf(
			"%s%s",
			et.Settings.TemplatePath,
			et.Settings.IndexTemplate,
		),
	)
	if err != nil {
		return err
	}

	// create file on disk
	var f *os.File
	f, err = os.Create(idx.DocPath + idx.IndexFileName)
	if err != nil {
		log.Fatal(err)
	}

	// execute index with index template
	err = tmpl.Execute(f, idx)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
