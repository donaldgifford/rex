package adr

import (
	"fmt"
	"os"

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
			Title: "README",
			Adrs:  getIndexAdrs(viper.GetString("adr.path")),
		},
	}
}

func getIndexAdrs(idxPath string) []*IndexAdrs {
	var myAdrs []*IndexAdrs

	entries, err := os.ReadDir(idxPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for i, e := range entries {
		myAdrs = append(myAdrs, &IndexAdrs{
			Id:    i,
			Title: e.Name(),
		})
		fmt.Println(e.Name())
	}

	return myAdrs
}
