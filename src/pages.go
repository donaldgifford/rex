package src

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// github pages
type ConfigGithubPages struct {
	adrPath       string
	templatesPath string
	cwd           string
}

func (c *ConfigGithubPages) Init() {
	c.generateADR()
	c.generatePages()
}

func (c *ConfigGithubPages) createDir(path string) {
	dirPath := fmt.Sprintf("%s/%s", c.cwd, path)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func (c *ConfigGithubPages) check(e error) {
	if e != nil {
		panic(e)
	}
}

func (c *ConfigGithubPages) generateADR() {
	c.createDir(c.templatesPath)
	c.createDir(c.adrPath)
	c.generateFile(c.templatesPath, pagesDefaultADR)
}

func (c *ConfigGithubPages) getDocsPaths() string {
	docs := strings.Split(c.adrPath, "adr/")
	return docs[0]
}

func (c *ConfigGithubPages) generatePages() {
	docPath := c.getDocsPaths()
	layoutDir := c.adrPath + "_layouts/"
	assetsDir := c.adrPath + "_assets/"
	c.createDir(layoutDir)
	c.createDir(assetsDir)
	c.generateFile(docPath, pagesIndex)
	c.generateFile(docPath, pagesConfig)
	c.generateFile(layoutDir, pagesWebLayoutDefault)
	c.generateFile(layoutDir, pagesWebLayoutAdr)
}

func (c *ConfigGithubPages) generateFile(path string, rexFile RexFile) {
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
