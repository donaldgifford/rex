/* package src

init.go

Methods and functions for initializing the rex tool.

Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/

package src

import (
	"fmt"
	"log"
	"os"
)

type ConfigAdr struct {
	adrPath       string
	templatesPath string
	cwd           string
}

// create directory structure with default files

func (c *ConfigAdr) Init() {
	c.generateADR()
}

func (c *ConfigAdr) createDir(path string) {
	dirPath := fmt.Sprintf("%s/%s", c.cwd, path)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func (c *ConfigAdr) generateFile(path string, rexFile RexFile) {
	filePath := fmt.Sprintf("%s/%s", c.cwd, path)

	fmt.Println(rexFile.name)

	dd := []byte(rexFile.content)
	fmt.Println(filePath + rexFile.name)
	f, err := os.Create(filePath + rexFile.name)
	c.check(err)
	defer f.Close()

	_, err = f.Write(dd)
	c.check(err)
}

func (c *ConfigAdr) generateADR() {
	c.createDir(c.templatesPath)
	c.createDir(c.adrPath)
	c.generateFile(c.templatesPath, defaultAdrTemplate)
}

func (c *ConfigAdr) check(e error) {
	if e != nil {
		panic(e)
	}
}
