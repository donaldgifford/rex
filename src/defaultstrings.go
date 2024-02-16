package src

import "github.com/spf13/viper"

type RexFile struct {
	name    string
	content string
}

var defaultAdrTemplate = RexFile{viper.GetString("template.adr.default"), `# {{ .Title }}

| Status | Author         | Date       |
| ------ | -------------- | ---------- |
| {{ .Status }} | {{ .Author }} | {{ .Date }} |

## Context and Problem Statement

## Decision Drivers


## Considererd Options


## Decision Outcome
  `}
