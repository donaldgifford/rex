---
id: {{.ID}}
title: {{.Title}}
type: {{.Type}}
status: planned
priority: {{.Priority}}
estimated_hours: {{.EstimatedHours}}
actual_hours: null
started_date: null
completed_date: null
blocked_by: []
blocks: []
related_to: []
assignee: {{if .Assignee}}{{.Assignee}}{{else}}null{{end}}
tags: [{{range $i, $tag := .Tags}}{{if $i}}, {{end}}{{$tag}}{{end}}]
phase: {{if .Phase}}{{.Phase}}{{else}}null{{end}}
---

# {{.Title}}

**Status**: 📋 Planned | **Priority**: {{.Priority}} | **Estimated**: {{.EstimatedHours}} hours

## Description

Clear description of what needs to be done and why.

## Context

Background information, related decisions, or constraints.

## Requirements

- Requirement 1
- Requirement 2
- Requirement 3

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Implementation Notes

Implementation approach, key decisions, or technical details.

## Testing

How to verify the task is complete.

## References

- Related ADRs, RFCs, docs, or external resources
