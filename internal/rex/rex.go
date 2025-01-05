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

// Package rex is the top level data structure that allows the cli commands to
// perform actions.
package rex

import (
	"github.com/donaldgifford/rex/internal/adr"
	"github.com/donaldgifford/rex/internal/config"
	"github.com/donaldgifford/rex/internal/templates"
)

// Rex holds the interfaces and data for performing actions needed
// by the cli calls.
type Rex struct {
	ADR      adr.IADR
	Index    adr.IIndex
	Template templates.Template
	Config   config.RexConfigure
}

// New creates a new Rex to use.
//
// Each part of the struct calls below its interfaces to
// use.
func New() *Rex {
	return &Rex{
		ADR:      adr.NewIADR(),
		Index:    adr.NewIIndex(),
		Template: templates.NewTemplate(),
		Config:   config.NewRexConfigure(),
	}
}

// Settings is a helper to return current settings
func (r *Rex) Settings() *config.RexConfig {
	return r.Config.Settings()
}

// NewADR creates a new ADR from content on disk and updates
// the current index configured.
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

func (r *Rex) UpdateIndex(force bool) error {
	err := r.Index.ADRs()
	if err != nil {
		return err
	}

	idx := r.Index.Execute()
	err = r.Template.GenerateIndex(idx, force)
	if err != nil {
		return err
	}
	return nil
}

// GenerateDirectories creates the default directories used for rex
// force is used to overwrite the templates if found
//
// if "templates.enabled: true" is the .rex.yaml config file
// then this function creates the default templates directory
func (r *Rex) GenerateDirectories() error {
	err := r.Config.Settings().GenerateDirectories()
	if err != nil {
		return err
	}
	return nil
}

// GenerateDefaultTemplates creates the default templates used for rex
//
// if force is set, it will overwrite the current template files if
// found with the defaults
func (r *Rex) GenerateTemplates(force bool) error {
	err := r.Config.Settings().GenerateDefaultTemplates(force)
	if err != nil {
		return err
	}
	return nil
}

// GenerateIndex updates the current index
//
// if force is set, it will overwrite the current index file if
// found.
func (r *Rex) GenerateIndex(force bool) error {
	err := r.Template.GenerateIndex(r.Index.Execute(), force)
	if err != nil {
		return err
	}

	return nil
}
