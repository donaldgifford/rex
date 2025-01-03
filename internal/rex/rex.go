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

package rex

import (
	"fmt"
	"os"

	"github.com/donaldgifford/rex/internal/adr"
	"github.com/donaldgifford/rex/internal/config"
	"github.com/donaldgifford/rex/internal/templates"
)

type Rex struct {
	ADR      adr.IADR
	Index    adr.IIndex
	Template templates.Template
	Config   config.RexConfigure
}

func New() *Rex {
	return &Rex{
		ADR:      adr.NewIADR(),
		Index:    adr.NewIIndex(),
		Template: templates.NewTemplate(),
		Config:   config.NewRexConfigure(),
	}
}

func (r *Rex) Settings() *config.RexConfig {
	return r.Config.Settings()
}

func (r *Rex) NewADR(content *adr.Content) error {
	// create adr
	adr, err := r.ADR.Create(content)
	if err != nil {
		return err
	}

	// write ADR to disk using template
	err = r.Template.CreateADR(adr)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rex) UpdateIndex() error {
	err := r.Index.ADRs()
	if err != nil {
		return err
	}

	idx := r.Index.Execute()
	err = r.Template.GenerateIndex(idx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rex) GenerateDirectories() error {
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create directory path string
	dirPath := fmt.Sprintf("%s/%s", cwd, r.ADR.GetSettings().Path)

	// mkdirall with path string
	err = os.MkdirAll(dirPath, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

func (r *Rex) GenerateIndex() error {
	err := r.Template.GenerateIndex(r.Index.Execute())
	if err != nil {
		return err
	}

	return nil
}
