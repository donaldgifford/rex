/*
Copyright © 2024-2025 Donald Gifford <dgifford06@gmail.com>

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
package cmd

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func ReadTestFile(file string) ([]byte, error) {
	t, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// path: "tests/docs/adr/"

func createConfigFile(file string, adrPath string, templatesEnabled bool, templatePath string) error {
	rexConfig := fmt.Sprintf(`adr:
  path: %s
  index_page: "README.md"
  add_to_index: true # on rex create, a new record will be added to the index page
templates:
  enabled: %v
  path: %s
  adr:
    default: "adr.tmpl"
    index: "index.tmpl"
enable_github_pages: false
pages:
  index: "index.md"
  web:
    config: "_config.yml"
    layout:
      adr: "adr.html"
      default: "default.html"
extras: true
extra_pages:
  install: install.md
  usage: usage.md"`, adrPath, templatesEnabled, templatePath)
	rc := []byte(rexConfig)
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(rc)
	if err != nil {
		return err
	}
	return nil
}

func createTestIndexTemplate(name string) error {
	idxTemplate := `# {{ .Content.Title }}

## ADRs

| ID | Title | Link |
| -- | ----- | ---- |
{{- range .Content.Adrs }}
| {{ .Id }} | {{ .Title }} | link |
{{- end }}`
	idx := []byte(idxTemplate)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(idx)
	if err != nil {
		return err
	}

	return nil
}

func createTestADRFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func createTestFolder(name string) error {
	err := os.MkdirAll(name, 0755)
	if err != nil {
		return err
	}
	return nil
}

func removeTestFolder(name string) error {
	err := os.RemoveAll(name)
	if err != nil {
		return err
	}
	return nil
}

func removeTestConfigFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	adrDocsPath := "tests/docs/adr/"

	err := createTestFolder(adrDocsPath)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createConfigFile("tests/.rex.yaml", "tests/docs/adr/", false, "tests/docs/templates/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createConfigFile("tests/.dirs-rex.yaml", "tests/dirs/docs/adr/", false, "tests/dirs/docs/templates/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createConfigFile("tests/.dirs-enabled-rex.yaml", "tests/dirs/docs/adr/", true, "tests/dirs/docs/templates/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestFolder("tests/dirs/docs/templates/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestIndexTemplate("tests/dirs/docs/templates/index.tmpl")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createConfigFile("tests/.dirs-error-rex.yaml", "tests/docs/adr/1-test1.md", false, "tests/dirs/docs/templates/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createConfigFile("tests/.templates-rex.yaml", "tests/templates/docs/adr/", false, "tests/templates/docs/templates/")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// setup some test files
	err = createTestADRFile(fmt.Sprintf("%s%s", adrDocsPath, "1-test1.md"))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestADRFile(fmt.Sprintf("%s%s", adrDocsPath, "2-test2.md"))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	err = createTestADRFile(fmt.Sprintf("%s%s", adrDocsPath, "README.md"))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	code := m.Run()

	err = removeTestFolder("tests")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	os.Exit(code)
}
