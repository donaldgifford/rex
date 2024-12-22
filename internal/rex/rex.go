package rex

import (
	"github.com/donaldgifford/rex/internal/adr"
	"github.com/donaldgifford/rex/internal/config"
	"github.com/donaldgifford/rex/internal/templates"
)

type Rex struct {
	ADR      adr.IADR
	Index    adr.IIndex
	Template templates.Template
	Config   config.RexConfigurer
}

func NewRex(install bool) *Rex {
	return &Rex{
		ADR:      adr.NewIADR(),
		Index:    adr.NewIIndex(),
		Template: templates.NewTemplate(),
		Config:   config.NewRexConfigurer(install),
	}
}

func (r *Rex) Settings() *config.RexConfig {
	return r.Config.Settings()
}

func (r *Rex) NewADR(content *adr.Content) error {
	// create adr
	adr, err := r.ADR.Create(content)
	if err != nil {
		return err
	}

	// write ADR to disk using template
	err = r.Template.CreateADR(adr)
	if err != nil {
		return err
	}

	return nil
}

func (r *Rex) UpdateIndex() error {
	err := r.Index.ADRs()
	if err != nil {
		return err
	}

	idx := r.Index.Execute()
	err = r.Template.GenerateIndex(idx)
	return nil
}
