/*
Copyright © 2024 Donald Gifford <dgifford06@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// The install package is responsible for providing the ability to initialize
// and setup the rex cli tool.
//
// This is mostly a hack and workaround cobra-cli. When you install the binary
// for rex, on first run you don't have a rex.yaml configuration file. This
// will allow a simple way to run `rex --install` on first run to create a
// rex.yaml file. After that is created, you can then run other commands to
// continue with the setup. This is cleaner than privious workaround.

package install

import (
	"os"

	"github.com/donaldgifford/rex/internal/config"
)

// These all map to a specific setting in the rex.yaml file for configuration
var (
	defaultAdrPath               string = "docs/adr/"
	defaultAdrIndexPage          string = "README.md"
	defaultAdrAddToIndex         bool   = true
	defaultTemplatesPath         string = "templates/"
	defaultTemplatesEnabled      bool   = false
	defaultTemplatesAdrDefault   string = "adr.tmpl"
	defaultTemplatesAdrIndex     string = "index.tmpl"
	defaultEnabledGithubPages    bool   = true
	defaultPagesIndex            string = "index.md"
	defaultPagesWebConfig        string = "_config.yml"
	defaultPagesWebLayoutAdr     string = "adr.html"
	defaultPagesWebLayoutDefault string = "default.html"
	defaultExtras                bool   = true
	defaultExtraPagesInstall     string = "install.md"
	defaultExtraPagesUsage       string = "usage.md"
)

// CreateRexConfigFile creates a rex.yaml file using the default settings
// returns error if the .rex.yaml file cannot be created or if the yaml
// output is incorrect.
func CreateRexConfigFile() error {
	// create RexConfig from default settings
	rex := &config.RexConfig{
		ADR: config.ADRConfig{
			Path:       defaultAdrPath,
			IndexPage:  defaultAdrIndexPage,
			AddToIndex: defaultAdrAddToIndex,
		},
		Templates: config.TemplateConfig{
			Enabled: defaultTemplatesEnabled,
			Path:    defaultTemplatesPath,
			ADR: config.ADRTemplateConfig{
				Default: defaultTemplatesAdrDefault,
				Index:   defaultTemplatesAdrIndex,
			},
		},
		EnableGithubPages: defaultEnabledGithubPages,
		Pages: config.PagesConfig{
			Index: defaultPagesIndex,
			Web: config.PagesConfigWeb{
				Config: defaultPagesWebConfig,
				Layout: config.PagesConfigWebLayout{
					ADR:     defaultPagesWebLayoutAdr,
					Default: defaultPagesWebLayoutDefault,
				},
			},
		},
		Extras: defaultExtras,
		ExtraPages: config.ExtrasConfig{
			Install: defaultExtraPagesInstall,
			Usage:   defaultExtraPagesUsage,
		},
	}

	// generate yaml output to use for configuration file
	yaml, err := rex.YamlOut()
	if err != nil {
		return err
	}

	// create the file
	f, err := os.Create(".rex.yaml")
	if err != nil {
		return err
	}

	defer f.Close()

	// write the file
	_, err = f.Write(yaml)
	if err != nil {
		return err
	}

	return nil
}
