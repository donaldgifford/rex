package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdrCreateCMD(t *testing.T) {
	tests := map[string]struct {
		file    string
		content string
		setArgs []string
		err     bool
	}{
		"adr": {
			file:    "tests/docs/adr/3-Test-ADR-Create.md",
			content: "# Test ADR Create\n\n| Status | Author         |  Created | Last Update | Current Version |\n| ------ | -------------- | -------- | ----------- | --------------- |\n| Draft | TESTER | 2024-12-27 | N/A | v0.0.1 |\n\n## Context and Problem Statement\n\n## Decision Drivers\n\n## Considered Options\n\n## Decision Outcome\n",
			// setArgs: []string{"--config=tests/.rex.yaml", "--title=Test ADR Create", "--author=TESTER"},
			setArgs: []string{"--config=tests/.rex.yaml", "adr", "create", "--title=Test ADR Create", "--author=TESTER"},
			err:     false,
		},
	}

	err := createConfigFile("tests/.rex.yaml")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(test.setArgs)
			// adrCreateCmd.SetOut(buf)
			// adrCreateCmd.SetArgs(test.setArgs)

			// err := adrCreateCmd.Execute()
			err := rootCmd.Execute()
			if err != nil {
				fmt.Println(err)
			}

			b, err := ReadTestFile(test.file)
			if err != nil {
				t.Errorf("error opening test file: %v, err: %v", test.file, err.Error())
			}
			assert.Equal(t, test.content, string(b), "")

			// if test.err {
			// 	assert.Error(t, err, fmt.Sprintf("Error: %v", err.Error()))
			// } else {
			// 	assert.Nil(t, err, "")
			// }
		})
	}

	err = removeTestConfigFile("tests/.rex.yaml")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
