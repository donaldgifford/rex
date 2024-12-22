package rex

import (
	"github.com/donaldgifford/rex/internal/config"
)

type Rex struct {
	Config *config.RexConfig
}

func NewRex() *Rex {
	return &Rex{
		Config: config.NewRexConfig(),
	}
}

func (r *Rex) Settings() *config.RexConfig {
	return r.Config.Settings()
}
