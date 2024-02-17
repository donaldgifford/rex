package src

import "github.com/spf13/viper"

type RexFile struct {
	fileName string
	content  string
}

func createRexFile(configString string, content string) RexFile {
	f := viper.GetString(configString)
	return RexFile{
		fileName: f,
		content:  content,
	}
}

func createExtrasRexFile(name string, content string) RexFile {
	return RexFile{
		fileName: name,
		content:  content,
	}
}

var (
	configDefaultAdrTemplate  string = "templates.adr.default"
	defaultAdrTemplateContent string = `# {{ .ADR.Title }}

| Status | Author         | Date       |
| ------ | -------------- | ---------- |
| {{ .ADR.Status }} | {{ .ADR.Author }} | {{ .ADR.Date }} |

## Context and Problem Statement

## Decision Drivers


## Considererd Options


## Decision Outcome

  `
)
