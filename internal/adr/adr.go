package adr

import (
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ADR struct {
	Content Content
	ID      int
	Path    string
}

type Content struct {
	Title  string
	Author string
	Status string
	Date   string
}

func getAdrFiles(path string) ([]string, error) {
	var files []string
	fileInfo, err := os.ReadDir(path)
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

func generateId(path string) (int, error) {
	adrs, err := getAdrFiles(path)
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

func NewADR(title string, author string, path string) (*ADR, error) {
	adrId, err := generateId(path)
	if err != nil {
		return &ADR{}, err
	}

	return &ADR{
		Content: Content{
			Title:  title,
			Author: author,
			Status: "pending",
			Date:   time.Now().Format(time.DateOnly),
		},
		ID:   adrId,
		Path: viper.GetString("adr.path"),
	}, nil
}

// func (a *ADR) GetADRFiles() ([]string, error) {
// 	var files []string
// 	fileInfo, err := os.ReadDir(a.Settings.Path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, file := range fileInfo {
// 		if file.Name() != "README.md" {
// 			files = append(files, file.Name())
// 		}
// 	}
//
// 	return files, nil
// }
//
// func (a *ADR) GenerateId() (int, error) {
// 	adrs, err := a.GetADRFiles()
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if len(adrs) == 0 {
// 		return 1, nil
// 	}
//
// 	var fileNames []int
// 	for _, v := range adrs {
// 		k := strings.Split(v, "-")[0]
// 		a, err := strconv.Atoi(k)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fileNames = append(fileNames, a)
// 	}
//
// 	lastestID := slices.Max(fileNames)
//
// 	return lastestID + 1, nil
// }
