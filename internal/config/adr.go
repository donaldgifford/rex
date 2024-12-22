package config

var (
	rexConfigFile        string = ".rex.yaml"
	rexConfigFileContent string = `adr:
  path: "docs/adr/"
  index_page: "README.md"
  add_to_index: true # on rex create, a new record will be added to the index page
templates:
  path: "templates/"
  adr:
    default: "adr.tmpl"`
)

// type ADRConfig struct {
// 	Path       string
// 	IndexPage  string
// 	AddToIndex bool
// }
//
// // newADRConfig reads the configuration settings under "adr"
// func NewADRConfig() *ADRConfig {
// 	return &ADRConfig{
// 		Path:       viper.GetString("adr.path"),
// 		IndexPage:  viper.GetString("adr.index_page"),
// 		AddToIndex: viper.GetBool("adr.add_to_index"),
// 	}
// }
//
// // GetSettings returns the settings for ADR
// func (a *ADRConfig) GetSettings() *ADRConfig {
// 	return a
// 	// settings := make(map[string]any)
// 	// settings["path"] = a.path
// 	// settings["indexPage"] = a.indexPage
// 	// settings["addToIndex"] = a.addToIndex
// 	//
// 	// return settings
// }
