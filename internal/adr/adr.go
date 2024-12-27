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

type ADR struct {
	Content Content
	ID      int
	Config  ADRConfig
}

type Content struct {
	Title  string
	Author string
	Status string
	Date   string
}

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
			// Date:   time.Now().Format(time.DateOnly),
		},
		ID: adrId,
	}, nil
}

func (adr *ADR) Revision(id int) (*ADR, error) {
	return nil, nil
}

func (adr *ADR) GetSettings() *ADRConfig {
	return adr.Config.GetSettings()
}

func NewADR() *ADR {
	return &ADR{
		Config: *NewADRConfig(),
	}
}
