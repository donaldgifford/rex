package templates

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/donaldgifford/rex/internal/adr"
)

//go:embed default/adr.tmpl
//go:embed default/index.tmpl
//go:embed default/rex.yaml
var DefaultRexTemplates embed.FS

type EmbeddedTemplate struct {
	Settings Settings
}

// GetSettings returns the settings for ADR
func (et *EmbeddedTemplate) GetSettings() *Settings {
	return &et.Settings
}

func (et *EmbeddedTemplate) Read(file string) ([]byte, error) {
	t, err := DefaultRexTemplates.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(t))
	return t, nil
}

func (et *EmbeddedTemplate) Execute()             {}
func (et *EmbeddedTemplate) GenerateDirectories() {}
func (et *EmbeddedTemplate) CreateADR(adr *adr.ADR) error {
	// get settings for file path and name
	// parse the template with Settings
	// write template to file

	tmpl, err := template.ParseFS(DefaultRexTemplates, fmt.Sprintf("%s%s", et.Settings.TemplatePath, et.Settings.AdrTemplate))
	if err != nil {
		return err
	}

	// tmpl, err := template.ParseFiles(fmt.Sprintf("%s%s", et.Settings.TemplatePath, et.Settings.AdrTemplate))
	// if err != nil {
	// 	return err
	// }

	strippedTitle := strings.Join(strings.Split(strings.Trim(adr.Content.Title, "\n \t"), " "), "-")
	fileName := fmt.Sprintf("%d-%s.md", adr.ID, strippedTitle)

	var f *os.File
	f, err = os.Create(adr.Path + fileName)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(f, adr)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (et *EmbeddedTemplate) GenerateIndex(idx *adr.Index) error {
	tmpl, err := template.ParseFS(DefaultRexTemplates, fmt.Sprintf("%s%s", et.Settings.TemplatePath, et.Settings.IndexTemplate))
	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(idx.DocPath + idx.IndexFileName)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(f, idx)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
