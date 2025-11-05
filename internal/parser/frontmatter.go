package parser

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ExtractFrontmatter extracts YAML frontmatter from a markdown file
// Returns the frontmatter content and the remaining markdown content
func ExtractFrontmatter(content []byte) ([]byte, []byte, error) {
	if !bytes.HasPrefix(content, []byte("---\n")) && !bytes.HasPrefix(content, []byte("---\r\n")) {
		return nil, nil, fmt.Errorf("missing frontmatter delimiter")
	}

	// Find the end of frontmatter
	lines := bytes.SplitN(content, []byte("\n"), -1)
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		line := bytes.TrimSpace(lines[i])
		if bytes.Equal(line, []byte("---")) {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil, nil, fmt.Errorf("unclosed frontmatter block")
	}

	// Extract frontmatter (skip first and last ---)
	frontmatter := bytes.Join(lines[1:endIdx], []byte("\n"))

	// Extract remaining content
	remaining := bytes.Join(lines[endIdx+1:], []byte("\n"))

	return frontmatter, remaining, nil
}

// ADRFrontmatter represents the YAML frontmatter for an ADR
type ADRFrontmatter struct {
	Number int    `yaml:"number"`
	Title  string `yaml:"title"`
	Status string `yaml:"status"`
	Date   string `yaml:"date"`
}

// Validate checks if the ADR frontmatter is valid
func (a *ADRFrontmatter) Validate() error {
	if a.Number <= 0 {
		return fmt.Errorf("number must be positive")
	}
	if strings.TrimSpace(a.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !isValidADRStatus(a.Status) {
		return fmt.Errorf("invalid status: %s (must be draft, accepted, deprecated, superseded)", a.Status)
	}
	if _, err := time.Parse("2006-01-02", a.Date); err != nil {
		return fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}
	return nil
}

func isValidADRStatus(status string) bool {
	validStatuses := []string{"draft", "accepted", "deprecated", "superseded"}
	for _, v := range validStatuses {
		if status == v {
			return true
		}
	}
	return false
}

// RFCFrontmatter represents the YAML frontmatter for an RFC
type RFCFrontmatter struct {
	Number      int    `yaml:"number"`
	Title       string `yaml:"title"`
	Status      string `yaml:"status"`
	Author      string `yaml:"author"`
	CreatedDate string `yaml:"created_date"`
	UpdatedDate string `yaml:"updated_date"`
}

// Validate checks if the RFC frontmatter is valid
func (r *RFCFrontmatter) Validate() error {
	if r.Number <= 0 {
		return fmt.Errorf("number must be positive")
	}
	if strings.TrimSpace(r.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !isValidRFCStatus(r.Status) {
		return fmt.Errorf("invalid status: %s (must be draft, review, approved, implemented, rejected)", r.Status)
	}
	if _, err := time.Parse("2006-01-02", r.CreatedDate); err != nil {
		return fmt.Errorf("invalid created_date format (expected YYYY-MM-DD): %w", err)
	}
	if _, err := time.Parse("2006-01-02", r.UpdatedDate); err != nil {
		return fmt.Errorf("invalid updated_date format (expected YYYY-MM-DD): %w", err)
	}
	return nil
}

func isValidRFCStatus(status string) bool {
	validStatuses := []string{"draft", "review", "approved", "implemented", "rejected"}
	for _, v := range validStatuses {
		if status == v {
			return true
		}
	}
	return false
}

// TaskFrontmatter represents the YAML frontmatter for a Task
type TaskFrontmatter struct {
	ID             string   `yaml:"id"`
	Title          string   `yaml:"title"`
	Type           string   `yaml:"type"`
	Status         string   `yaml:"status"`
	Priority       string   `yaml:"priority"`
	EstimatedHours float64  `yaml:"estimated_hours"`
	ActualHours    *float64 `yaml:"actual_hours"`
	StartedDate    *string  `yaml:"started_date"`
	CompletedDate  *string  `yaml:"completed_date"`
	BlockedBy      []string `yaml:"blocked_by"`
	Blocks         []string `yaml:"blocks"`
	RelatedTo      []string `yaml:"related_to"`
	Assignee       *string  `yaml:"assignee"`
	Tags           []string `yaml:"tags"`
	Phase          *int     `yaml:"phase"`
}

// Validate checks if the Task frontmatter is valid
func (t *TaskFrontmatter) Validate() error {
	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("id is required")
	}
	if !strings.HasPrefix(t.ID, "TASK-") {
		return fmt.Errorf("id must start with TASK-")
	}
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !isValidTaskType(t.Type) {
		return fmt.Errorf("invalid type: %s (must be core, plugin, ui, other)", t.Type)
	}
	if !isValidTaskStatus(t.Status) {
		return fmt.Errorf("invalid status: %s (must be planned, in_progress, blocked, completed, cancelled)", t.Status)
	}
	if !isValidTaskPriority(t.Priority) {
		return fmt.Errorf("invalid priority: %s (must be P0, P1, P2, P3)", t.Priority)
	}
	if t.EstimatedHours <= 0 {
		return fmt.Errorf("estimated_hours must be positive")
	}
	if t.ActualHours != nil && *t.ActualHours < 0 {
		return fmt.Errorf("actual_hours must be non-negative")
	}
	if t.StartedDate != nil {
		if _, err := time.Parse("2006-01-02", *t.StartedDate); err != nil {
			return fmt.Errorf("invalid started_date format (expected YYYY-MM-DD): %w", err)
		}
	}
	if t.CompletedDate != nil {
		if _, err := time.Parse("2006-01-02", *t.CompletedDate); err != nil {
			return fmt.Errorf("invalid completed_date format (expected YYYY-MM-DD): %w", err)
		}
	}
	if t.Phase != nil && *t.Phase < 0 {
		return fmt.Errorf("phase must be non-negative")
	}

	// Initialize empty arrays if nil
	if t.BlockedBy == nil {
		t.BlockedBy = []string{}
	}
	if t.Blocks == nil {
		t.Blocks = []string{}
	}
	if t.RelatedTo == nil {
		t.RelatedTo = []string{}
	}
	if t.Tags == nil {
		t.Tags = []string{}
	}

	return nil
}

func isValidTaskType(taskType string) bool {
	validTypes := []string{"core", "plugin", "ui", "other"}
	for _, v := range validTypes {
		if taskType == v {
			return true
		}
	}
	return false
}

func isValidTaskStatus(status string) bool {
	validStatuses := []string{"planned", "in_progress", "blocked", "completed", "cancelled"}
	for _, v := range validStatuses {
		if status == v {
			return true
		}
	}
	return false
}

func isValidTaskPriority(priority string) bool {
	validPriorities := []string{"P0", "P1", "P2", "P3"}
	for _, v := range validPriorities {
		if priority == v {
			return true
		}
	}
	return false
}

// PlanFrontmatter represents the YAML frontmatter for a Plan
type PlanFrontmatter struct {
	Number      int    `yaml:"number"`
	Title       string `yaml:"title"`
	Status      string `yaml:"status"`
	CreatedDate string `yaml:"created_date"`
}

// Validate checks if the Plan frontmatter is valid
func (p *PlanFrontmatter) Validate() error {
	if p.Number <= 0 {
		return fmt.Errorf("number must be positive")
	}
	if strings.TrimSpace(p.Title) == "" {
		return fmt.Errorf("title is required")
	}
	if !isValidPlanStatus(p.Status) {
		return fmt.Errorf("invalid status: %s (must be draft, active, completed, cancelled)", p.Status)
	}
	if _, err := time.Parse("2006-01-02", p.CreatedDate); err != nil {
		return fmt.Errorf("invalid created_date format (expected YYYY-MM-DD): %w", err)
	}
	return nil
}

func isValidPlanStatus(status string) bool {
	validStatuses := []string{"draft", "active", "completed", "cancelled"}
	for _, v := range validStatuses {
		if status == v {
			return true
		}
	}
	return false
}

// ParseADR parses ADR frontmatter from markdown content
func ParseADR(content []byte) (*ADRFrontmatter, []byte, error) {
	frontmatter, body, err := ExtractFrontmatter(content)
	if err != nil {
		return nil, nil, fmt.Errorf("extract frontmatter: %w", err)
	}

	var adr ADRFrontmatter
	if err := yaml.Unmarshal(frontmatter, &adr); err != nil {
		return nil, nil, fmt.Errorf("parse YAML: %w", err)
	}

	if err := adr.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate: %w", err)
	}

	return &adr, body, nil
}

// ParseRFC parses RFC frontmatter from markdown content
func ParseRFC(content []byte) (*RFCFrontmatter, []byte, error) {
	frontmatter, body, err := ExtractFrontmatter(content)
	if err != nil {
		return nil, nil, fmt.Errorf("extract frontmatter: %w", err)
	}

	var rfc RFCFrontmatter
	if err := yaml.Unmarshal(frontmatter, &rfc); err != nil {
		return nil, nil, fmt.Errorf("parse YAML: %w", err)
	}

	if err := rfc.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate: %w", err)
	}

	return &rfc, body, nil
}

// ParseTask parses Task frontmatter from markdown content
func ParseTask(content []byte) (*TaskFrontmatter, []byte, error) {
	frontmatter, body, err := ExtractFrontmatter(content)
	if err != nil {
		return nil, nil, fmt.Errorf("extract frontmatter: %w", err)
	}

	var task TaskFrontmatter
	if err := yaml.Unmarshal(frontmatter, &task); err != nil {
		return nil, nil, fmt.Errorf("parse YAML: %w", err)
	}

	if err := task.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate: %w", err)
	}

	return &task, body, nil
}

// ParsePlan parses Plan frontmatter from markdown content
func ParsePlan(content []byte) (*PlanFrontmatter, []byte, error) {
	frontmatter, body, err := ExtractFrontmatter(content)
	if err != nil {
		return nil, nil, fmt.Errorf("extract frontmatter: %w", err)
	}

	var plan PlanFrontmatter
	if err := yaml.Unmarshal(frontmatter, &plan); err != nil {
		return nil, nil, fmt.Errorf("parse YAML: %w", err)
	}

	if err := plan.Validate(); err != nil {
		return nil, nil, fmt.Errorf("validate: %w", err)
	}

	return &plan, body, nil
}
