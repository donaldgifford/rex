package src

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// github pages
type ConfigGithubPages struct {
	adrPath       string
	templatesPath string
	cwd           string
	force         bool
}

func (c *ConfigGithubPages) Init() {
	c.generateADR()
	c.generatePages()
	c.extraPages()
}

func (c *ConfigGithubPages) overwrite(b bool) {
	c.force = true
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
	c.generateFile(c.templatesPath, pagesDefaultADRConfig, pagesDefaultADRContent)
}

func (c *ConfigGithubPages) getDocsPaths() string {
	docs := strings.Split(c.adrPath, "adr/")
	return docs[0]
}

func (c *ConfigGithubPages) generatePages() {
	docPath := c.getDocsPaths()
	layoutDir := docPath + "_layouts/"
	assetsDir := docPath + "_assets/"
	c.createDir(layoutDir)
	c.createDir(assetsDir)
	c.generateFile(docPath, pagesIndexConfig, pagesIndexContent)
	c.generateFile(docPath, pagesWebConfig, setConfig(c.cwd))
	c.generateFile(layoutDir, pagesWebLayoutDefault, pagesWebLayoutDefaultContent)
	c.generateFile(layoutDir, pagesWebLayoutAdr, pagesWebLayoutAdrContent)
}

func (c *ConfigGithubPages) extraPages() {
	if viper.GetBool("extras") {
		// fmt.Printf("Extra Pages: %v\n", viper.GetStringSlice("pages.extra_pages"))
		docsPath := c.getDocsPaths()
		c.generateFile(docsPath, extraPagesInstallConfig, extraPagesContent)
		c.generateFile(docsPath, extraPagesUsageConfig, extraPagesContent)
	} else {
		fmt.Println("Extras not enabled.")
	}
	// c.generateFile(docsPath, extraPagesConfig, extraPagesContent)
}

func (c *ConfigGithubPages) generateFile(path string, rexFileName string, rexFileContent string) {
	filePath := fmt.Sprintf("%s/%s", c.cwd, path)

	rexFile := createRexFile(rexFileName, rexFileContent)
	// fmt.Printf("rexFileName: %s\n", rexFile.fileName)
	// fmt.Printf("rexFileContent: %s\n", rexFile.content)

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
		file.Close()
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("File: %s does not exist\n", filePath+rexFile.fileName)
			fmt.Printf("Creating file: %s", filePath+rexFile.fileName)
			f, err := os.Create(filePath + rexFile.fileName)
			c.check(err)
			defer f.Close()

			_, err = f.Write(dd)
			c.check(err)
		} else {
			fmt.Printf("File: %s exists\n", filePath+rexFile.fileName)
		}
	}
	// filePath := fmt.Sprintf("%s/%s", c.cwd, path)
	//
	// rexFile := createRexFile(rexFileName, rexFileContent)
	//
	// dd := []byte(rexFile.content)
	// fmt.Println(filePath + rexFile.fileName)
	// f, err := os.Create(filePath + rexFile.fileName)
	// c.check(err)
	// defer f.Close()
	//
	// _, err = f.Write(dd)
	// c.check(err)
}
