package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderADR(t *testing.T) {
	data := ADRData{
		Number: 1,
		Title:  "Use SQLite as Cache Layer",
		Date:   "2025-11-05",
	}

	result, err := RenderADR(data)
	require.NoError(t, err)

	// Check frontmatter
	assert.Contains(t, result, "number: 1")
	assert.Contains(t, result, "title: Use SQLite as Cache Layer")
	assert.Contains(t, result, "status: draft")
	assert.Contains(t, result, "date: 2025-11-05")

	// Check content
	assert.Contains(t, result, "# 1. Use SQLite as Cache Layer")
	assert.Contains(t, result, "**Status**: Draft")
	assert.Contains(t, result, "## Context")
	assert.Contains(t, result, "## Decision")
	assert.Contains(t, result, "## Consequences")
	assert.Contains(t, result, "## Alternatives Considered")
}

func TestRenderRFC(t *testing.T) {
	data := RFCData{
		Number:      1,
		Title:       "Rex Documentation Management System",
		Author:      "Donald Gifford",
		CreatedDate: "2025-11-01",
		UpdatedDate: "2025-11-05",
	}

	result, err := RenderRFC(data)
	require.NoError(t, err)

	// Check frontmatter
	assert.Contains(t, result, "number: 1")
	assert.Contains(t, result, "title: Rex Documentation Management System")
	assert.Contains(t, result, "status: draft")
	assert.Contains(t, result, "author: Donald Gifford")
	assert.Contains(t, result, "created_date: 2025-11-01")
	assert.Contains(t, result, "updated_date: 2025-11-05")

	// Check content
	assert.Contains(t, result, "# RFC 1: Rex Documentation Management System")
	assert.Contains(t, result, "**Author**: Donald Gifford")
	assert.Contains(t, result, "## Executive Summary")
	assert.Contains(t, result, "## Problem Statement")
	assert.Contains(t, result, "## Proposed Solution")
}

func TestRenderTask(t *testing.T) {
	assignee := "test-user"
	phase := 1
	data := TaskData{
		ID:             "TASK-001",
		Title:          "Implement SQLite Cache",
		Type:           "core",
		Priority:       "P1",
		EstimatedHours: 8,
		Assignee:       &assignee,
		Tags:           []string{"database", "cache"},
		Phase:          &phase,
	}

	result, err := RenderTask(data)
	require.NoError(t, err)

	// Check frontmatter
	assert.Contains(t, result, "id: TASK-001")
	assert.Contains(t, result, "title: Implement SQLite Cache")
	assert.Contains(t, result, "type: core")
	assert.Contains(t, result, "status: planned")
	assert.Contains(t, result, "priority: P1")
	assert.Contains(t, result, "estimated_hours: 8")
	assert.Contains(t, result, "assignee: test-user")
	assert.Contains(t, result, "tags: [database, cache]")
	assert.Contains(t, result, "phase: 1")

	// Check content
	assert.Contains(t, result, "# Implement SQLite Cache")
	assert.Contains(t, result, "**Priority**: P1")
	assert.Contains(t, result, "**Estimated**: 8 hours")
	assert.Contains(t, result, "## Description")
	assert.Contains(t, result, "## Acceptance Criteria")
}

func TestRenderTaskWithNullFields(t *testing.T) {
	data := TaskData{
		ID:             "TASK-002",
		Title:          "Test Task",
		Type:           "ui",
		Priority:       "P2",
		EstimatedHours: 4,
		Assignee:       nil,
		Tags:           []string{},
		Phase:          nil,
	}

	result, err := RenderTask(data)
	require.NoError(t, err)

	// Check that null fields are rendered as "null"
	assert.Contains(t, result, "assignee: null")
	assert.Contains(t, result, "tags: []")
	assert.Contains(t, result, "phase: null")
}

func TestRenderPlan(t *testing.T) {
	data := PlanData{
		Number:      1,
		Title:       "Phase 1 Implementation",
		CreatedDate: "2025-11-01",
	}

	result, err := RenderPlan(data)
	require.NoError(t, err)

	// Check frontmatter
	assert.Contains(t, result, "number: 1")
	assert.Contains(t, result, "title: Phase 1 Implementation")
	assert.Contains(t, result, "status: draft")
	assert.Contains(t, result, "created_date: 2025-11-01")

	// Check content
	assert.Contains(t, result, "# Plan 1: Phase 1 Implementation")
	assert.Contains(t, result, "**Created**: 2025-11-01")
	assert.Contains(t, result, "## Overview")
	assert.Contains(t, result, "## Goals")
	assert.Contains(t, result, "## Timeline")
}

func TestRenderADRValidFrontmatter(t *testing.T) {
	// Test that rendered template produces valid frontmatter that can be parsed
	data := ADRData{
		Number: 42,
		Title:  "Test ADR",
		Date:   "2025-11-05",
	}

	result, err := RenderADR(data)
	require.NoError(t, err)

	// Ensure frontmatter delimiters are present
	lines := strings.Split(result, "\n")
	assert.Equal(t, "---", lines[0], "First line should be frontmatter delimiter")

	// Find closing delimiter
	closingIdx := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			closingIdx = i
			break
		}
	}
	assert.Greater(t, closingIdx, 0, "Should have closing frontmatter delimiter")
}
