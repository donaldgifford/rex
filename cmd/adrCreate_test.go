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
package cmd

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func parseContentWithDate(content string) string {
	d := time.Now()
	formattedDate := d.Format("2006-01-02")
	return fmt.Sprintf(content, formattedDate)
}

func TestAdrCreateCMD(t *testing.T) {
	tests := map[string]struct {
		file    string
		content string
		setArgs []string
		err     bool
	}{
		"adr": {
			file:    "tests/docs/adr/3-Test-ADR-Create.md",
			content: parseContentWithDate("# Test ADR Create\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| Draft | TESTER | %s | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome\n"),
			setArgs: []string{"--config=tests/.rex.yaml", "adr", "create", "--title=Test ADR Create", "--author=TESTER"},
			err:     false,
		},
	}

	// err := createConfigFile("tests/.rex.yaml")
	// if err != nil {
	// 	log.Print(err)
	// 	os.Exit(1)
	// }

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(test.setArgs)

			err := rootCmd.Execute()
			if err != nil {
				fmt.Println(err)
			}

			b, err := ReadTestFile(test.file)
			if err != nil {
				t.Errorf("error opening test file: %v, err: %v", test.file, err.Error())
			}
			assert.Equal(t, test.content, string(b), "")
		})
	}

	// err = removeTestConfigFile("tests/.rex.yaml")
	// if err != nil {
	// 	log.Print(err)
	// 	os.Exit(1)
	// }
}
