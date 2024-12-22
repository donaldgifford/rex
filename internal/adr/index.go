/*
* Eventually index will work like terraform-docs where you put an anchor in the readme or whatever markdown file and then
* it will parse the template into that file between those anchors.
*
*
 */
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
		return nil
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
