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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"

	"github.com/donaldgifford/rex/internal/adr"
)

type RexTemplate struct {
	Settings Settings
}

// GetSettings returns the settings for ADR
func (rt *RexTemplate) GetSettings() *Settings {
	return &rt.Settings
}

// "default/rex.yaml"
// TODO: Change to read from disk - settings in .rex.yaml
// still not sure why i need this
func (rt *RexTemplate) Read(file string) ([]byte, error) {
	cleanFile := filepath.Clean(file)
	t, err := os.ReadFile(cleanFile)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (rt *RexTemplate) Execute() {}

func (rt *RexTemplate) CreateADR(adr *adr.ADR) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s%s", rt.Settings.TemplatePath, rt.Settings.AdrTemplate))
	if err != nil {
		return err
	}

	strippedTitle := strings.Join(strings.Split(strings.Trim(adr.Content.Title, "\n \t"), " "), "-")
	fileName := fmt.Sprintf("%d-%s.md", adr.ID, strippedTitle)

	var f *os.File
	cleanFile := filepath.Clean(fmt.Sprintf("%s%s", viper.GetString("adr.path"), fileName))
	f, err = os.Create(cleanFile)
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

func (rt *RexTemplate) GenerateIndex(idx *adr.Index, force bool) error {
	if force {
		err := rt.writeIndex(idx)
		if err != nil {
			return err
		}

		return nil
	}

	if fileExists(idx.DocPath + idx.IndexFileName) {
		return fmt.Errorf("index file found at %s, to overwrite please pass --force flag", idx.DocPath+idx.IndexFileName)
	}

	err := rt.writeIndex(idx)
	if err != nil {
		return err
	}

	return nil
}

func (rt *RexTemplate) writeIndex(idx *adr.Index) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s%s", rt.Settings.TemplatePath, rt.Settings.IndexTemplate))
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
