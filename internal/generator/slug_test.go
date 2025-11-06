package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "basic title",
			title: "Use SQLite as Cache Layer",
			want:  "use-sqlite-as-cache-layer",
		},
		{
			name:  "with special characters",
			title: "Don't Use Special Characters!",
			want:  "dont-use-special-characters",
		},
		{
			name:  "with multiple spaces",
			title: "Multiple   Spaces   Here",
			want:  "multiple-spaces-here",
		},
		{
			name:  "with hyphens",
			title: "Pre-existing-hyphens",
			want:  "pre-existing-hyphens",
		},
		{
			name:  "with numbers",
			title: "Phase 1 Implementation",
			want:  "phase-1-implementation",
		},
		{
			name:  "with leading/trailing spaces",
			title: "  Trim Me  ",
			want:  "trim-me",
		},
		{
			name:  "with parentheses",
			title: "Function (with params)",
			want:  "function-with-params",
		},
		{
			name:  "with colons",
			title: "RFC 0001: Title Here",
			want:  "rfc-0001-title-here",
		},
		{
			name:  "empty string",
			title: "",
			want:  "",
		},
		{
			name:  "only special characters",
			title: "!@#$%^&*()",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateADRFilename(t *testing.T) {
	tests := []struct {
		name   string
		number int
		title  string
		want   string
	}{
		{
			name:   "single digit",
			number: 1,
			title:  "Use SQLite as Cache Layer",
			want:   "0001-use-sqlite-as-cache-layer.md",
		},
		{
			name:   "double digit",
			number: 42,
			title:  "Test ADR",
			want:   "0042-test-adr.md",
		},
		{
			name:   "triple digit",
			number: 123,
			title:  "Another Test",
			want:   "0123-another-test.md",
		},
		{
			name:   "four digit",
			number: 1234,
			title:  "Many ADRs",
			want:   "1234-many-adrs.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateADRFilename(tt.number, tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateRFCFilename(t *testing.T) {
	tests := []struct {
		name   string
		number int
		title  string
		want   string
	}{
		{
			name:   "basic RFC",
			number: 1,
			title:  "Rex Documentation Management System",
			want:   "0001-rex-documentation-management-system.md",
		},
		{
			name:   "with special chars",
			number: 5,
			title:  "Don't Do This!",
			want:   "0005-dont-do-this.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRFCFilename(tt.number, tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateTaskFilename(t *testing.T) {
	tests := []struct {
		name  string
		id    string
		title string
		want  string
	}{
		{
			name:  "basic task",
			id:    "TASK-001",
			title: "Implement SQLite Cache",
			want:  "TASK-001-implement-sqlite-cache.md",
		},
		{
			name:  "with special chars",
			id:    "TASK-042",
			title: "Fix Bug: Memory Leak",
			want:  "TASK-042-fix-bug-memory-leak.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTaskFilename(tt.id, tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGeneratePlanFilename(t *testing.T) {
	tests := []struct {
		name   string
		number int
		title  string
		want   string
	}{
		{
			name:   "basic plan",
			number: 1,
			title:  "Phase 1 Implementation",
			want:   "0001-phase-1-implementation.md",
		},
		{
			name:   "larger number",
			number: 99,
			title:  "Big Plan",
			want:   "0099-big-plan.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GeneratePlanFilename(tt.number, tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}
