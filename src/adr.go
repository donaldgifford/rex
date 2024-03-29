/* package src

init.go

Methods and functions for initializing the rex tool.

Copyright © 2024 Donald Gifford <dgifford06@gmail.com>
*/

package src

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ConfigAdr struct {
	adrPath       string
	templatesPath string
	cwd           string
	force         bool
}

// create directory structure with default files

func (c *ConfigAdr) Init() {
	c.generateADR()
}

func (c *ConfigAdr) overwrite(b bool) {
	c.force = true
}

func (c *ConfigAdr) createDir(path string) {
	dirPath := fmt.Sprintf("%s/%s", c.cwd, path)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func (c *ConfigAdr) generateFile(path string, rexFileName string, rexFileContent string) {
	filePath := fmt.Sprintf("%s/%s", c.cwd, path)

	rexFile := createRexFile(rexFileName, rexFileContent)

	if c.force {
		dd := []byte(rexFile.content)

		f, err := os.Create(filePath + rexFile.fileName)
		c.check(err)
		defer f.Close()

		_, err = f.Write(dd)
		c.check(err)

	} else {
		dd := []byte(rexFile.content)
		fmt.Println(filePath + rexFile.fileName)
		file, err := os.Open(filePath + rexFile.fileName)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("File: %s does not exist\n", filePath+rexFile.fileName)
			fmt.Printf("Creating file: %s", filePath+rexFile.fileName)
			defer file.Close()
			_, err = file.Write(dd)
			c.check(err)
		} else {
			fmt.Printf("File: %s exists", filePath+rexFile.fileName)
		}
	}
}

func (c *ConfigAdr) generateADR() {
	c.createDir(c.templatesPath)
	c.createDir(c.adrPath)
	c.generateFile(c.templatesPath, configDefaultAdrTemplate, defaultAdrTemplateContent)
}

func (c *ConfigAdr) check(e error) {
	if e != nil {
		fmt.Printf("Error: %s", e.Error())
		panic(e)
	}
}

func (c *ConfigAdr) getDocsPaths() string {
	docs := strings.Split(c.adrPath, "adr/")
	return docs[0]
}

func (c *ConfigAdr) extraPages() {
	if viper.GetBool("extras") {
		docsPath := c.getDocsPaths()
		c.generateFile(docsPath, extraPagesInstallConfig, extraPagesContent)
		c.generateFile(docsPath, extraPagesUsageConfig, extraPagesContent)
	} else {
		fmt.Println("Extras not enabled.")
	}
}
