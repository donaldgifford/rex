package templates

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed files/adr/*.md files/rfc/*.md files/task/*.md files/plan/*.md
var templatesFS embed.FS

// ADRData holds data for rendering ADR templates
type ADRData struct {
	Number int
	Title  string
	Date   string
}

// RFCData holds data for rendering RFC templates
type RFCData struct {
	Number      int
	Title       string
	Author      string
	CreatedDate string
	UpdatedDate string
}

// TaskData holds data for rendering Task templates
type TaskData struct {
	ID             string
	Title          string
	Type           string
	Priority       string
	EstimatedHours float64
	Assignee       *string
	Tags           []string
	Phase          *int
}

// PlanData holds data for rendering Plan templates
type PlanData struct {
	Number      int
	Title       string
	CreatedDate string
}

// RenderADR renders the ADR template with the given data
func RenderADR(data ADRData) (string, error) {
	content, err := templatesFS.ReadFile("files/adr/template.md")
	if err != nil {
		return "", fmt.Errorf("read ADR template: %w", err)
	}

	tmpl, err := template.New("adr").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("parse ADR template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute ADR template: %w", err)
	}

	return buf.String(), nil
}

// RenderRFC renders the RFC template with the given data
func RenderRFC(data RFCData) (string, error) {
	content, err := templatesFS.ReadFile("files/rfc/template.md")
	if err != nil {
		return "", fmt.Errorf("read RFC template: %w", err)
	}

	tmpl, err := template.New("rfc").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("parse RFC template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute RFC template: %w", err)
	}

	return buf.String(), nil
}

// RenderTask renders the Task template with the given data
func RenderTask(data TaskData) (string, error) {
	content, err := templatesFS.ReadFile("files/task/template.md")
	if err != nil {
		return "", fmt.Errorf("read Task template: %w", err)
	}

	tmpl, err := template.New("task").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("parse Task template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute Task template: %w", err)
	}

	return buf.String(), nil
}

// RenderPlan renders the Plan template with the given data
func RenderPlan(data PlanData) (string, error) {
	content, err := templatesFS.ReadFile("files/plan/template.md")
	if err != nil {
		return "", fmt.Errorf("read Plan template: %w", err)
	}

	tmpl, err := template.New("plan").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("parse Plan template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute Plan template: %w", err)
	}

	return buf.String(), nil
}
