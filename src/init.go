/* package src

init.go

Methods and functions for initializing the rex tool.

Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/

package src

// func Init() {
// }
//
// func CreateConfig() {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	createDefaultConfigFile(cwd)
// }
//
// func createDefaultConfigFile(cwd string) {
// 	rexConfigFile := ".rex.yaml"
//
// 	rexConfig := `adr:
//   path: "docs/adr/"
//   index_page: "README.md"
//   add_to_index: true # on rex create, a new record will be added to the index page
// templates:
//   path: "templates/"
//   adr:
//     default: "adr.tmpl"
// enable_github_pages: false`
// 	rc := []byte(rexConfig)
// 	fileName := cwd + "/" + rexConfigFile
// 	fmt.Println(fileName)
// 	f, err := os.Create(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
//
// 	_, err = f.Write(rc)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
