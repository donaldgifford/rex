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

	"github.com/donaldgifford/rex/internal/adr"
)

// TODO: add custom template generation later
var configDefaultAdrTemplate string = "templates.adr.default"

type RexTemplate struct {
	Settings Settings
}

// GetSettings returns the settings for ADR
func (rt *RexTemplate) GetSettings() *Settings {
	return &rt.Settings
}

// "default/rex.yaml"
// TODO: Change to read from disk - settings in .rex.yaml
func (rt *RexTemplate) Read(file string) ([]byte, error) {
	t, err := DefaultRexTemplates.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(t))
	return t, nil
}

func (rt *RexTemplate) Execute() {}

func (rt *RexTemplate) CreateADR(adr *adr.ADR) error {
	return nil
}

func (rt *RexTemplate) GenerateIndex(idx *adr.Index) error {
	return nil
}
