/* templates.go
*
* templates package is responsible for embedding the default template files used to created the rex config yaml file and other optionial templates. The embedded files can be found in the 'default' directory inside this package.
*
* templates provides an interface to allow creating new custom ADR's, docs, or other files to be generated and used with rex.
*
* IE:
*
* default templates for ADR, github pages, for output to json, markdown, html
* create custom templates for ADR or load them from a specified directory via config or cli flag
 */
package templates

import (
	"github.com/donaldgifford/rex/internal/adr"
	"github.com/spf13/viper"
)

type Template interface {
	Read(file string) ([]byte, error) // read in template file from embedded directory
	Execute()                         // Execute the template with passed in configuration variables
	GetSettings() *Settings
	CreateADR(adr *adr.ADR) error
}

func NewTemplate() Template {
	if viper.GetBool("templates.enabled") {
		return &RexTemplate{
			Settings: Settings{
				TemplatePath: viper.GetString("templates.path"),
				AdrTemplate:  viper.GetString("templates.adr.default"),
			},
		}
	} else {
		return &EmbeddedTemplate{
			Settings: Settings{
				TemplatePath: "default/",
				AdrTemplate:  "adr.tmpl",
			},
		}
	}
}

type Settings struct {
	TemplatePath string
	AdrTemplate  string
}

// import (
// 	"embed"
// 	"fmt"
// )
//
// // TODO: add custom template generation later
// var configDefaultAdrTemplate string = "templates.adr.default"
//
// //go:embed default/adr.tmpl
// //go:embed default/rex.yaml
// //go:embed default/gh/adr.tmpl
// var defaultRexTemplates embed.FS
//
// type Template interface {
// 	Read()    // read in template file from embedded directory
// 	Execute() // Execute the template with passed in configuration variables
// }
//
// type ADR struct {
// 	Title  string
// 	Author string
// 	Status string
// 	Date   string
// }
//
// func (a *ADR) Read() {
// 	t, _ := defaultRexTemplates.ReadFile("default/adr.tmpl")
// 	fmt.Println(string(t))
// }
// func (a *ADR) Execute() {}
//
// type GhPages struct {
// 	ADR
// 	WebTitle string
// 	Layout   string
// }
