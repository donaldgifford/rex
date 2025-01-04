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

// templates.go
//
//  templates package is responsible for embedding the default template files used to created the rex config yaml file and other optionial templates. The embedded files can be found in the 'default' directory inside this package.
//
//  templates provides an interface to allow creating new custom ADR's, docs, or other files to be generated and used with rex.
//
//  IE:
//
//  default templates for ADR, github pages, for output to json, markdown, html
//  create custom templates for ADR or load them from a specified directory via config or cli flag

package templates

import (
	"os"

	"github.com/spf13/viper"

	"github.com/donaldgifford/rex/internal/adr"
)

type Template interface {
	Read(file string) ([]byte, error) // read in template file from embedded directory
	Execute()                         // Execute the template with passed in configuration variables
	GetSettings() *Settings
	CreateADR(adr *adr.ADR) error
	GenerateIndex(idx *adr.Index, force bool) error
}

func NewTemplate() Template {
	if viper.GetBool("templates.enabled") {
		return &RexTemplate{
			Settings: Settings{
				TemplatePath:  viper.GetString("templates.path"),
				AdrTemplate:   viper.GetString("templates.adr.default"),
				IndexTemplate: viper.GetString("templates.adr.index"),
			},
		}
	} else {
		return &EmbeddedTemplate{
			Settings: Settings{
				TemplatePath:  "default/",
				AdrTemplate:   "adr.tmpl",
				IndexTemplate: "index.tmpl",
			},
		}
	}
}

type Settings struct {
	TemplatePath  string
	AdrTemplate   string
	IndexTemplate string
}

// fileExists returns checks if a file already exists on disk
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
