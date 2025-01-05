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

//go:embed default/adr.tmpl
//go:embed default/index.tmpl
//go:embed default/index_readme.tmpl
var DefaultRexTemplates embed.FS

type EmbeddedTemplate struct {
	Settings Settings
}

// GetSettings returns the settings for ADR
func (et *EmbeddedTemplate) GetSettings() *Settings {
	return &et.Settings
}

func (et *EmbeddedTemplate) Read(file string) ([]byte, error) {
	t, err := DefaultRexTemplates.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (et *EmbeddedTemplate) Execute() {}

func (et *EmbeddedTemplate) CreateADR(adr *adr.ADR) error {
	// get settings for file path and name
	// parse the template with Settings
	// write template to file

	tmpl, err := template.ParseFS(DefaultRexTemplates, fmt.Sprintf("%s%s", et.Settings.TemplatePath, et.Settings.AdrTemplate))
	if err != nil {
		return err
	}

	strippedTitle := strings.Join(strings.Split(strings.Trim(adr.Content.Title, "\n \t"), " "), "-")
	fileName := fmt.Sprintf("%d-%s.md", adr.ID, strippedTitle)

	var f *os.File
	f, err = os.Create(viper.GetString("adr.path") + fileName)
	if err != nil {
		log.Fatal(err)
	}

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

func (et *EmbeddedTemplate) GenerateIndex(idx *adr.Index, force bool) error {
	if force {
		err := et.writeIndex(idx)
		if err != nil {
			return err
		}

		return nil
	}

	if fileExists(idx.DocPath + idx.IndexFileName) {
		return fmt.Errorf("index file found at %s, to overwrite please pass --force flag", idx.DocPath+idx.IndexFileName)
	}

	err := et.writeIndex(idx)
	if err != nil {
		return err
	}

	return nil
}

func (et *EmbeddedTemplate) writeIndex(idx *adr.Index) error {
	tmpl, err := template.ParseFS(DefaultRexTemplates, fmt.Sprintf("%s%s", et.Settings.TemplatePath, et.Settings.IndexTemplate))
	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(idx.DocPath + idx.IndexFileName)
	if err != nil {
		log.Fatal(err)
	}

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
