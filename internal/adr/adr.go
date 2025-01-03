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
package adr

import (
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type IADR interface {
	Create(content *Content) (*ADR, error)
	Revision(id int) (*ADR, error)
	Id() (int, error)
	GetSettings() *ADRConfig
}

func NewIADR() IADR {
	return NewADR()
}

// ADR is the data that is sent to the templates to be created
type ADR struct {
	Content Content
	ID      int
	Config  ADRConfig
}

// Content is the input for creating a new ADR
type Content struct {
	Title  string
	Author string
	Status string
	Date   string
}

// ADRConfig holds configuration for where ADR's are written to, what
// the index page is, and if ADR's are added to the index page.
type ADRConfig struct {
	Path       string
	IndexPage  string
	AddToIndex bool
}

// newADRConfig reads the configuration settings under "adr"
func NewADRConfig() *ADRConfig {
	return &ADRConfig{
		Path:       viper.GetString("adr.path"),
		IndexPage:  viper.GetString("adr.index_page"),
		AddToIndex: viper.GetBool("adr.add_to_index"),
	}
}

// GetSettings returns the settings for ADR
func (a *ADRConfig) GetSettings() *ADRConfig {
	return a
}

// GetAdrFilesNames returns a slice of strings containing the
// file names found in the ADR Path minus the index page.
// Returns error if path cannot be found.
func (adr *ADR) GetAdrFilesNames() ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(adr.Config.Path)
	if err != nil {
		return nil, err
	}
	for _, file := range fileInfo {
		if file.Name() != "README.md" {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

// Id is a helper function that returns the next int
// to use as the ID for an ADR
func (adr *ADR) Id() (int, error) {
	adrs, err := adr.GetAdrFilesNames()
	if err != nil {
		return 0, err
	}

	if len(adrs) == 0 {
		return 1, nil
	}

	var fileNames []int
	for _, v := range adrs {
		k := strings.Split(v, "-")[0]
		a, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal(err)
		}
		fileNames = append(fileNames, a)
	}

	lastestID := slices.Max(fileNames)

	return lastestID + 1, nil
}

// Create takes a content pointer and returns an ADR pointer and error.
func (adr *ADR) Create(content *Content) (*ADR, error) {
	adrId, err := adr.Id()
	if err != nil {
		return &ADR{}, err
	}

	return &ADR{
		Content: Content{
			Title:  content.Title,
			Author: content.Author,
			Status: content.Status,
			Date:   content.Date,
		},
		ID: adrId,
	}, nil
}

// TODO: Revision takes the current ADR and creates a
// revision for it.
func (adr *ADR) Revision(id int) (*ADR, error) {
	return nil, nil
}

// GetSettings returns the ADRConfig settings for the
// ADR.
func (adr *ADR) GetSettings() *ADRConfig {
	return adr.Config.GetSettings()
}

// NewADR returns an ADR with a NewADRConfig set in
// the Config.
func NewADR() *ADR {
	return &ADR{
		Config: *NewADRConfig(),
	}
}
