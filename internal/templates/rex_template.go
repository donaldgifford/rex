package templates

import (
	"fmt"

	"github.com/donaldgifford/rex/internal/adr"
)

// TODO: add custom template generation later
var configDefaultAdrTemplate string = "templates.adr.default"

type RexTemplate struct {
	Settings Settings
}

// GetSettings returns the settings for ADR
func (rt *RexTemplate) GetSettings() *Settings {
	return &rt.Settings
}

// "default/rex.yaml"
// TODO: Change to read from disk - settings in .rex.yaml
func (rt *RexTemplate) Read(file string) ([]byte, error) {
	t, err := DefaultRexTemplates.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(t))
	return t, nil
}

func (rt *RexTemplate) Execute() {}

func (rt *RexTemplate) CreateADR(adr *adr.ADR) error {
	return nil
}
