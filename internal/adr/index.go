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

// Eventually index will work like terraform-docs where you put an anchor in the readme or whatever markdown file and then
// it will parse the template into that file between those anchors.
package adr

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type IIndex interface {
	Execute() *Index
	ADRs() error
	Process(fileName string) *IndexAdr
}

type Index struct {
	DocPath       string
	IndexFileName string
	Content       IndexContent
}

type IndexContent struct {
	Title string
	Adrs  []*IndexAdr
}

type IndexAdr struct {
	Id    int
	Title string
}

func NewIIndex() *Index {
	return &Index{
		DocPath:       viper.GetString("adr.path"),
		IndexFileName: viper.GetString("adr.index_page"),
		Content: IndexContent{
			Title: "ADR Index",
		},
	}
}

func (idx *Index) Execute() *Index {
	return idx
}

func (idx *Index) ADRs() error {
	var myAdrs []*IndexAdr

	entries, err := os.ReadDir(idx.DocPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, e := range entries {
		if e.Name() != idx.IndexFileName {
			adr := idx.Process(e.Name())
			myAdrs = append(myAdrs, adr)
			fmt.Println(e.Name())
		}
	}
	idx.Content.Adrs = myAdrs
	return nil
}

func (idx *Index) Process(file string) *IndexAdr {
	idTitle := strings.SplitN(file, "-", 2)
	id, _ := strconv.Atoi(idTitle[0])
	title := strings.TrimSuffix(idTitle[1], ".md")

	return &IndexAdr{
		Id:    id,
		Title: title,
	}
}
