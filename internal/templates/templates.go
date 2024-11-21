package templates

var (
	configDefaultAdrTemplate  string = "templates.adr.default"
	defaultAdrTemplateContent string = `# {{ .ADR.Title }}

| Status | Author         | Date       |
| ------ | -------------- | ---------- |
| {{ .ADR.Status }} | {{ .ADR.Author }} | {{ .ADR.Date }} |

## Context and Problem Statement

## Decision Drivers


## Considererd Options


## Decision Outcome`
)

type ADR struct {
	Name      string
	IsDefault bool
	Content   string
}

func (a *ADR) Create()   {}
func (a *ADR) Generate() {}

type ITemplate interface {
	Create()
	Generate()
}
