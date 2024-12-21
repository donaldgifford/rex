package adr

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Index struct {
	DocPath       string
	IndexFileName string
	Content       IndexContent
}

type IndexContent struct {
	Title string
	Adrs  []*IndexAdrs
}

type IndexAdrs struct {
	Id    int
	Title string
}

func NewIndex() *Index {
	return &Index{
		DocPath:       viper.GetString("adr.path"),
		IndexFileName: viper.GetString("adr.index_page"),
		Content: IndexContent{
			Title: "ADR Index",
		},
	}
}

func (idx *Index) GetIndexAdrs() error {
	var myAdrs []*IndexAdrs

	entries, err := os.ReadDir(idx.DocPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, e := range entries {
		if e.Name() != idx.IndexFileName {
			id, title := idx.ProcessIndexAdrs(e.Name())
			myAdrs = append(myAdrs, &IndexAdrs{
				Id:    id,
				Title: title,
			})
			fmt.Println(e.Name())
		}
	}
	idx.Content.Adrs = myAdrs
	return nil
}

func (idx *Index) ProcessIndexAdrs(file string) (int, string) {
	idTitle := strings.SplitN(file, "-", 2)
	id, _ := strconv.Atoi(idTitle[0])
	title := strings.TrimSuffix(idTitle[1], ".md")

	return id, title
}
