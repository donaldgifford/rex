# Task Management

Rex provides a powerful task management system built on markdown files with YAML frontmatter and SQLite caching for fast queries.

## Core Principles

- **Markdown First**: Tasks are plain markdown files (source of truth)
- **YAML Frontmatter**: Structured metadata for filtering and relationships
- **SQLite Cache**: Fast queries without parsing all files
- **Git-Friendly**: Text files with clear diffs
- **Offline-First**: No external dependencies or APIs
- **Auto-Sync**: Cache automatically rebuilds when stale

## Directory Structure

```
docs/tasks/
   core/              # Backend/API tasks
      active/
      completed/
   plugins/           # Plugin development tasks
      active/
      completed/
   ui/                # Frontend tasks
      active/
      completed/
   other/             # Infrastructure, docs, operations
       active/
       completed/
```

**Note**: The `rex` tool will support generating these directories with `rex init`.

## YAML Frontmatter Fields

Each task file contains 15 structured fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier (TASK-001) |
| `title` | string | Yes | Task title |
| `type` | string | Yes | Domain: core, plugin, ui, other |
| `status` | string | Yes | planned, in_progress, blocked, completed, cancelled |
| `priority` | string | Yes | P0 (blocker), P1 (high), P2 (medium), P3 (low) |
| `estimated_hours` | number | Yes | Estimated time to complete |
| `actual_hours` | number/null | No | Actual time spent (set on completion) |
| `started_date` | string/null | No | When work began (YYYY-MM-DD) |
| `completed_date` | string/null | No | When work finished (YYYY-MM-DD) |
| `blocked_by` | array | Yes | Task IDs blocking this task |
| `blocks` | array | Yes | Task IDs this task blocks |
| `related_to` | array | Yes | Related task IDs |
| `assignee` | string/null | No | Assigned person |
| `tags` | array | Yes | Searchable tags |
| `phase` | number/null | No | Project phase |

**Example**:
```yaml
---
id: TASK-008
title: Implement SQLite Cache Layer
type: core
status: in_progress
priority: P1
estimated_hours: 8
actual_hours: null
started_date: 2025-11-05
completed_date: null
blocked_by: []
blocks: [TASK-010]
related_to: [TASK-007]
assignee: null
tags: [database, cache, sqlite]
phase: 1
---
```

## Rex CLI Commands

### Create a New Task

```bash
rex task create
```

**Interactive prompts for**:
- Title (required)
- Type (core, plugin, ui, other)
- Priority (P0-P3)
- Estimated hours
- Assignee (optional)
- Phase (optional)
- Tags (comma-separated)
- Relationships (blocked_by, blocks, related_to)

**Actions**:
- Auto-generates unique TASK-NNN ID
- Creates filename slug from title
- Places in correct domain/active/ directory
- Writes markdown file
- Updates SQLite cache
- Updates README files

### List Tasks (Fast SQLite Queries)

```bash
# All tasks
rex task list

# Filter by type
rex task list --type core

# Filter by status
rex task list --status in_progress

# Filter by priority
rex task list --priority P0,P1

# Combined filters
rex task list --type core --status in_progress --priority P1

# Filter by tags
rex task list --tags database,cache
```

**Performance**: SQLite queries return results instantly, even with 100+ tasks.

### View Task Statistics

```bash
rex task stats
```

**Shows**:
- Breakdown by domain (core, plugin, ui, other)
- Breakdown by status (planned, in_progress, blocked, completed)
- Breakdown by priority (P0, P1, P2, P3)
- Time tracking summary (estimated vs actual)
- Lists blocked tasks with dependencies
- Progress visualization

**Example output**:
```
Task Statistics
================

By Type:
  core    : 15 tasks (120 hours estimated)
  plugin  : 8 tasks (64 hours estimated)
  ui      : 5 tasks (40 hours estimated)
  other   : 2 tasks (8 hours estimated)

By Status:
  planned     : 18 tasks
  in_progress : 5 tasks
  blocked     : 2 tasks
  completed   : 5 tasks

By Priority:
  P0 (blocker): 1 task
  P1 (high)   : 8 tasks
  P2 (medium) : 15 tasks
  P3 (low)    : 6 tasks

Time Tracking:
  Estimated: 232 hours
  Actual: 45 hours
  Remaining: 187 hours
```

### Complete a Task

```bash
rex task complete TASK-008
```

**Actions**:
- Prompts for actual hours spent
- Updates frontmatter (status, completed_date, actual_hours)
- Moves file from {type}/active/ to {type}/completed/
- Updates SQLite cache
- Updates README files

### Update Task READMEs

```bash
rex task update
```

**Actions**:
- Queries SQLite for all tasks
- Generates root docs/tasks/README.md
- Generates domain READMEs (core/, plugin/, ui/, other/)
- Shows statistics and active/completed tables

**Note**: This happens automatically on create/complete, but can be run manually if needed.

## Task Lifecycle

```
📋 Planned → 🔨 In Progress → ✅ Completed
                    ↓
                🚫 Blocked
                    ↓
            (dependency resolved)
                    ↓
                🔨 In Progress
```

### Status Values

| Status | Emoji | Meaning |
|--------|-------|---------|
| `planned` | 📋 | Task defined but not started |
| `in_progress` | 🔨 | Actively being worked on |
| `blocked` | 🚫 | Waiting on dependencies (tracked via `blocked_by`) |
| `completed` | ✅ | Finished successfully, moved to {type}/completed/ |
| `cancelled` | ❌ | Decided not to implement |

### Priority Levels

| Priority | Meaning |
|----------|---------|
| `P0` | Blocker - stops all work |
| `P1` | High - important, work on next |
| `P2` | Medium - normal priority |
| `P3` | Low - nice to have |

## Task File Format

**File naming**: `TASK-NNN-slug-of-title.md`

**Location**: `docs/tasks/{type}/active/TASK-NNN-slug-of-title.md`

**Example**: `docs/tasks/core/active/TASK-008-implement-sqlite-cache.md`

```markdown
---
id: TASK-008
title: Implement SQLite Cache Layer
type: core
status: in_progress
priority: P1
estimated_hours: 8
actual_hours: null
started_date: 2025-11-05
completed_date: null
blocked_by: []
blocks: [TASK-010]
related_to: [TASK-007]
assignee: null
tags: [database, cache, sqlite]
phase: 1
---

# Implement SQLite Cache Layer

**Status**: 🔨 In Progress | **Priority**: P1 (High) | **Estimated**: 8 hours

## Description

Implement SQLite as a cache layer for rex documentation management...

## Context

Rex needs fast queries for filtering tasks without parsing all markdown files...

## Requirements

- [ ] Create SQLite schema for tasks, ADRs, RFCs, plans
- [ ] Implement auto-rebuild on stale detection
- [ ] Add relationship tables for blocked_by/blocks
- [ ] Add indexes for fast filtering

## Acceptance Criteria

- [ ] `rex task list --type core --status in_progress` returns results in <100ms
- [ ] Cache auto-rebuilds when markdown files change
- [ ] `rex rebuild` forces full cache rebuild
- [ ] All task relationships queryable via SQL JOINs

## Implementation Notes

Using modernc.org/sqlite (pure Go, no CGO) for portability...

## Testing

- Unit tests for cache operations
- Integration tests for auto-rebuild
- Performance tests with 100+ tasks

## References

- ADR 0001: Use SQLite as Cache Layer
- ADR 0005: Auto-Rebuild Cache on Stale Detection
- RFC 0001: Rex Documentation Management System
```

## SQLite Cache

### How It Works

1. **Markdown is Source of Truth**: All tasks are markdown files in docs/tasks/
2. **SQLite is Cache**: `.rex.db` stores parsed task data for fast queries
3. **Auto-Sync**: Rex detects stale cache and rebuilds automatically
4. **Manual Rebuild**: `rex rebuild` forces full rebuild from markdown

### Cache Benefits

- **Fast Queries**: No need to parse 100+ markdown files
- **Complex Filtering**: SQL WHERE clauses for multi-field filtering
- **Aggregations**: SQL GROUP BY for statistics
- **Relationships**: SQL JOINs for task dependencies

### Cache Location

```
.rex.db    # SQLite database (gitignored)
```

**Important**: Add `.rex.db` to your `.gitignore` - it's a cache, not source data.

### Stale Detection

Rex automatically rebuilds cache when:
- `.rex.db` doesn't exist
- Any markdown file is newer than last cache rebuild
- User runs `rex rebuild` manually

**Performance**: Rebuild takes 1-3 seconds for 100+ tasks, but only happens after git operations or manual edits.

## Migration from Bash Scripts

If migrating from the old bash script system:

1. **Install rex**: `go install github.com/donaldgifford/rex@latest`
2. **Rebuild cache**: `rex rebuild` (scans all existing tasks)
3. **Verify**: `rex task list` should show all tasks
4. **Remove old tools**: `rm -rf tools/makefiles tools/docs tools/scripts`
5. **Update workflows**: Replace `make task` with `rex task create`, etc.

## Why Markdown Files? (Not GitHub Issues)

**Rationale**:
- **Offline-first**: Work without internet, no API rate limits
- **Git-tracked**: Full history, easy rollback, branch-safe
- **Simple**: No external dependencies, no UI bugs
- **Flexible**: Custom fields, any structure
- **Searchable**: Use grep, ripgrep, fzf
- **Versionable**: Tasks change with code in same commit

**Rejected alternatives**: GitHub Issues, Jira, Trello/Linear (see TASK_CLAUDE_OLD.md for details)

## References

- ADR 0001: Use SQLite as Cache Layer
- ADR 0004: Use YAML Frontmatter for Metadata
- ADR 0005: Auto-Rebuild Cache on Stale Detection
- RFC 0001: Rex Documentation Management System
- Old system documentation: TASK_CLAUDE_OLD.md
