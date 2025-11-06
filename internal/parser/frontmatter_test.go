package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractFrontmatter(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		wantFrontmatter string
		wantBody       string
		wantErr        bool
	}{
		{
			name: "valid frontmatter",
			content: `---
title: Test
status: draft
---

# Body content`,
			wantFrontmatter: "title: Test\nstatus: draft",
			wantBody:        "\n# Body content",
			wantErr:         false,
		},
		{
			name: "missing opening delimiter",
			content: `title: Test
---
Body`,
			wantErr: true,
		},
		{
			name: "missing closing delimiter",
			content: `---
title: Test
Body`,
			wantErr: true,
		},
		{
			name: "empty frontmatter",
			content: `---
---
Body`,
			wantFrontmatter: "",
			wantBody:        "Body",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frontmatter, body, err := ExtractFrontmatter([]byte(tt.content))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantFrontmatter, string(frontmatter))
			assert.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestParseADR(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *ADRFrontmatter
		wantErr bool
	}{
		{
			name: "valid ADR",
			content: `---
number: 1
title: Use SQLite as Cache Layer
status: accepted
date: 2025-11-05
---

# Body`,
			want: &ADRFrontmatter{
				Number: 1,
				Title:  "Use SQLite as Cache Layer",
				Status: "accepted",
				Date:   "2025-11-05",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			content: `---
number: 1
title: Test
status: invalid
date: 2025-11-05
---

Body`,
			wantErr: true,
		},
		{
			name: "invalid date format",
			content: `---
number: 1
title: Test
status: draft
date: 11/05/2025
---

Body`,
			wantErr: true,
		},
		{
			name: "missing title",
			content: `---
number: 1
title: ""
status: draft
date: 2025-11-05
---

Body`,
			wantErr: true,
		},
		{
			name: "zero number",
			content: `---
number: 0
title: Test
status: draft
date: 2025-11-05
---

Body`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adr, body, err := ParseADR([]byte(tt.content))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, adr)
			assert.NotEmpty(t, body)
		})
	}
}

func TestParseRFC(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *RFCFrontmatter
		wantErr bool
	}{
		{
			name: "valid RFC",
			content: `---
number: 1
title: Rex Documentation Management System
status: approved
author: Donald Gifford
created_date: 2025-11-01
updated_date: 2025-11-05
---

# Body`,
			want: &RFCFrontmatter{
				Number:      1,
				Title:       "Rex Documentation Management System",
				Status:      "approved",
				Author:      "Donald Gifford",
				CreatedDate: "2025-11-01",
				UpdatedDate: "2025-11-05",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			content: `---
number: 1
title: Test
status: invalid
author: Test
created_date: 2025-11-01
updated_date: 2025-11-05
---

Body`,
			wantErr: true,
		},
		{
			name: "invalid created_date",
			content: `---
number: 1
title: Test
status: draft
author: Test
created_date: invalid-date
updated_date: 2025-11-05
---

Body`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rfc, body, err := ParseRFC([]byte(tt.content))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, rfc)
			assert.NotEmpty(t, body)
		})
	}
}

func TestParseTask(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *TaskFrontmatter
		wantErr bool
	}{
		{
			name: "valid task with all fields",
			content: `---
id: TASK-001
title: Implement SQLite Cache
type: core
status: in_progress
priority: P1
estimated_hours: 8
actual_hours: 6
started_date: 2025-11-05
completed_date: null
blocked_by: []
blocks: [TASK-002]
related_to: [TASK-003]
assignee: null
tags: [database, cache]
phase: 1
---

# Body`,
			want: &TaskFrontmatter{
				ID:             "TASK-001",
				Title:          "Implement SQLite Cache",
				Type:           "core",
				Status:         "in_progress",
				Priority:       "P1",
				EstimatedHours: 8,
				ActualHours:    ptrFloat64(6),
				StartedDate:    ptrString("2025-11-05"),
				CompletedDate:  nil,
				BlockedBy:      []string{},
				Blocks:         []string{"TASK-002"},
				RelatedTo:      []string{"TASK-003"},
				Assignee:       nil,
				Tags:           []string{"database", "cache"},
				Phase:          ptrInt(1),
			},
			wantErr: false,
		},
		{
			name: "minimal valid task",
			content: `---
id: TASK-002
title: Test Task
type: ui
status: planned
priority: P2
estimated_hours: 4
actual_hours: null
started_date: null
completed_date: null
blocked_by: []
blocks: []
related_to: []
assignee: null
tags: []
phase: null
---

# Body`,
			want: &TaskFrontmatter{
				ID:             "TASK-002",
				Title:          "Test Task",
				Type:           "ui",
				Status:         "planned",
				Priority:       "P2",
				EstimatedHours: 4,
				ActualHours:    nil,
				StartedDate:    nil,
				CompletedDate:  nil,
				BlockedBy:      []string{},
				Blocks:         []string{},
				RelatedTo:      []string{},
				Assignee:       nil,
				Tags:           []string{},
				Phase:          nil,
			},
			wantErr: false,
		},
		{
			name: "invalid ID format",
			content: `---
id: 001
title: Test
type: core
status: planned
priority: P1
estimated_hours: 4
actual_hours: null
started_date: null
completed_date: null
blocked_by: []
blocks: []
related_to: []
assignee: null
tags: []
phase: null
---

Body`,
			wantErr: true,
		},
		{
			name: "invalid type",
			content: `---
id: TASK-001
title: Test
type: invalid
status: planned
priority: P1
estimated_hours: 4
actual_hours: null
started_date: null
completed_date: null
blocked_by: []
blocks: []
related_to: []
assignee: null
tags: []
phase: null
---

Body`,
			wantErr: true,
		},
		{
			name: "invalid status",
			content: `---
id: TASK-001
title: Test
type: core
status: invalid
priority: P1
estimated_hours: 4
actual_hours: null
started_date: null
completed_date: null
blocked_by: []
blocks: []
related_to: []
assignee: null
tags: []
phase: null
---

Body`,
			wantErr: true,
		},
		{
			name: "invalid priority",
			content: `---
id: TASK-001
title: Test
type: core
status: planned
priority: P99
estimated_hours: 4
actual_hours: null
started_date: null
completed_date: null
blocked_by: []
blocks: []
related_to: []
assignee: null
tags: []
phase: null
---

Body`,
			wantErr: true,
		},
		{
			name: "negative estimated_hours",
			content: `---
id: TASK-001
title: Test
type: core
status: planned
priority: P1
estimated_hours: -1
actual_hours: null
started_date: null
completed_date: null
blocked_by: []
blocks: []
related_to: []
assignee: null
tags: []
phase: null
---

Body`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, body, err := ParseTask([]byte(tt.content))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, task)
			assert.NotEmpty(t, body)
		})
	}
}

func TestParsePlan(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *PlanFrontmatter
		wantErr bool
	}{
		{
			name: "valid plan",
			content: `---
number: 1
title: Phase 1 Implementation
status: active
created_date: 2025-11-01
---

# Body`,
			want: &PlanFrontmatter{
				Number:      1,
				Title:       "Phase 1 Implementation",
				Status:      "active",
				CreatedDate: "2025-11-01",
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			content: `---
number: 1
title: Test
status: invalid
created_date: 2025-11-01
---

Body`,
			wantErr: true,
		},
		{
			name: "invalid date",
			content: `---
number: 1
title: Test
status: active
created_date: invalid
---

Body`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, body, err := ParsePlan([]byte(tt.content))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, plan)
			assert.NotEmpty(t, body)
		})
	}
}

func TestADRFrontmatterValidate(t *testing.T) {
	tests := []struct {
		name    string
		adr     ADRFrontmatter
		wantErr bool
	}{
		{
			name: "valid",
			adr: ADRFrontmatter{
				Number: 1,
				Title:  "Test",
				Status: "draft",
				Date:   "2025-11-05",
			},
			wantErr: false,
		},
		{
			name: "all valid statuses",
			adr: ADRFrontmatter{
				Number: 1,
				Title:  "Test",
				Status: "accepted",
				Date:   "2025-11-05",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.adr.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskFrontmatterValidate(t *testing.T) {
	tests := []struct {
		name    string
		task    TaskFrontmatter
		wantErr bool
	}{
		{
			name: "valid with all fields",
			task: TaskFrontmatter{
				ID:             "TASK-001",
				Title:          "Test",
				Type:           "core",
				Status:         "planned",
				Priority:       "P1",
				EstimatedHours: 4,
			},
			wantErr: false,
		},
		{
			name: "initializes empty arrays",
			task: TaskFrontmatter{
				ID:             "TASK-001",
				Title:          "Test",
				Type:           "core",
				Status:         "planned",
				Priority:       "P1",
				EstimatedHours: 4,
				BlockedBy:      nil,
				Blocks:         nil,
				RelatedTo:      nil,
				Tags:           nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Check that nil arrays are initialized
				assert.NotNil(t, tt.task.BlockedBy)
				assert.NotNil(t, tt.task.Blocks)
				assert.NotNil(t, tt.task.RelatedTo)
				assert.NotNil(t, tt.task.Tags)
			}
		})
	}
}

// Helper functions for pointer values
func ptrFloat64(f float64) *float64 {
	return &f
}

func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}
