package templates

import (
	"embed"
	"fmt"

	"github.com/donaldgifford/rex/internal/adr"
)

//go:embed default/adr.tmpl
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
func (et *EmbeddedTemplate) CreateADR(adr *adr.ADR) error
