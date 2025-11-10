# Rex Command Reference

Complete reference for all `rex` commands, flags, and usage examples.

## Table of Contents

- [Global Flags](#global-flags)
- [ADR Commands](#adr-commands)
- [RFC Commands](#rfc-commands)
- [Task Commands](#task-commands)
- [Plan Commands](#plan-commands)
- [Database Commands](#database-commands)
- [Utility Commands](#utility-commands)

---

## Global Flags

These flags work with any rex command:

```bash
--config string    Config file (default: $HOME/.rex.yaml)
--help, -h         Show help for any command
```

**Examples:**

```bash
# Use custom config file
rex --config /path/to/config.yaml task list

# Show help for a command
rex task --help
rex task create --help
```

---

## ADR Commands

Architecture Decision Records (ADRs) document significant architectural decisions.

### rex adr create

Create a new ADR from template.

**Usage:**
```bash
rex adr create <title>
rex adr create "Your Decision Title"
```

**Arguments:**
- `<title>` - Title of the ADR (required, use quotes for multi-word titles)

**Behavior:**
- Auto-generates next ADR number (e.g., 0001, 0002, etc.)
- Creates file: `docs/adr/NNNN-title-slug.md`
- Sets status to "draft"
- Opens file for editing after creation

**Template Structure:**
```markdown
---
id: 0001
title: Your Decision Title
status: draft
date: 2025-11-10
---

# Your Decision Title

## Status

Draft

## Context

Background and problem statement.

## Decision

What was decided and why.

## Consequences

Positive and negative impacts.
```

**Examples:**

```bash
# Create ADR
rex adr create "Use SQLite as Cache Layer"
# Output: ✓ Created ADR 0001: Use SQLite as Cache Layer
#         File: docs/adr/0001-use-sqlite-as-cache-layer.md

# Create another ADR
rex adr create "Adopt Go Embed for Templates"
# Output: ✓ Created ADR 0002: Adopt Go Embed for Templates
```

**Related Commands:**
- `rex adr list` - View all ADRs
- `rex adr update` - Regenerate README

---

### rex adr list

List all ADRs with optional filtering.

**Usage:**
```bash
rex adr list [flags]
```

**Flags:**
```bash
--status string    Filter by status (draft, accepted, deprecated, superseded)
```

**Output Format:**
```
ADRs
════

ID    Title                              Status     Date
──────────────────────────────────────────────────────────────
0001  Use SQLite as Cache Layer          accepted   2025-11-10
0002  Adopt Go Embed for Templates       draft      2025-11-10
0003  Use Cobra for CLI Framework        accepted   2025-11-09
```

**Examples:**

```bash
# List all ADRs
rex adr list

# List only accepted ADRs
rex adr list --status accepted

# List draft ADRs
rex adr list --status draft

# List deprecated ADRs
rex adr list --status deprecated
```

**Status Values:**
- `draft` - Proposed, not yet decided
- `accepted` - Approved and implemented
- `deprecated` - No longer applicable
- `superseded` - Replaced by another ADR

---

### rex adr update

Regenerate `docs/adr/README.md` with current ADR list.

**Usage:**
```bash
rex adr update
```

**Behavior:**
- Scans all ADR files in `docs/adr/`
- Generates table of ADRs grouped by status
- Preserves content outside auto-generated markers
- Updates file in-place

**Generated README Structure:**
```markdown
# Architecture Decision Records

<!-- BEGIN AUTO-GENERATED -->

## Accepted

- [ADR 0001: Use SQLite](0001-use-sqlite.md) - 2025-11-10
- [ADR 0003: Use Cobra CLI](0003-use-cobra.md) - 2025-11-09

## Draft

- [ADR 0002: Go Embed Templates](0002-go-embed.md) - 2025-11-10

<!-- END AUTO-GENERATED -->
```

**Examples:**

```bash
# Update ADR README
rex adr update
# Output: ✓ Updated docs/adr/README.md
```

**When to Use:**
- After creating new ADRs
- After changing ADR status
- After modifying ADR titles

---

## RFC Commands

Requests for Comments (RFCs) propose features or changes that need discussion.

### rex rfc create

Create a new RFC from template.

**Usage:**
```bash
rex rfc create <title>
rex rfc create "Your Proposal Title"
```

**Arguments:**
- `<title>` - Title of the RFC (required, use quotes for multi-word titles)

**Behavior:**
- Auto-generates next RFC number (e.g., 0001, 0002, etc.)
- Creates file: `docs/rfc/NNNN-title-slug.md`
- Sets status to "draft"
- Prompts for author name

**Template Structure:**
```markdown
---
id: 0001
title: Your Proposal Title
status: draft
author: Your Name
date: 2025-11-10
---

# RFC 0001: Your Proposal Title

## Summary

One-paragraph explanation of the proposal.

## Motivation

Why are we doing this?

## Detailed Design

How does it work?

## Alternatives

What other approaches were considered?
```

**Examples:**

```bash
# Create RFC
rex rfc create "Add Plugin System"
> Author: John Doe
# Output: ✓ Created RFC 0001: Add Plugin System

# Create RFC (author prompted)
rex rfc create "Improve Task Statistics"
> Author: Jane Smith
# Output: ✓ Created RFC 0002: Improve Task Statistics
```

---

### rex rfc list

List all RFCs with optional filtering.

**Usage:**
```bash
rex rfc list [flags]
```

**Flags:**
```bash
--status string    Filter by status (draft, review, approved, implemented, rejected)
```

**Output Format:**
```
RFCs
════

ID    Title                    Author        Status        Date
────────────────────────────────────────────────────────────────────
0001  Add Plugin System        John Doe      approved      2025-11-10
0002  Improve Statistics       Jane Smith    draft         2025-11-09
```

**Examples:**

```bash
# List all RFCs
rex rfc list

# List RFCs in review
rex rfc list --status review

# List approved RFCs
rex rfc list --status approved

# List implemented RFCs
rex rfc list --status implemented
```

**Status Values:**
- `draft` - Initial proposal
- `review` - Under review/discussion
- `approved` - Approved for implementation
- `implemented` - Completed
- `rejected` - Not accepted

---

### rex rfc update

Regenerate `docs/rfc/README.md` with current RFC list.

**Usage:**
```bash
rex rfc update
```

**Behavior:**
- Scans all RFC files in `docs/rfc/`
- Generates table of RFCs grouped by status
- Preserves content outside auto-generated markers

**Examples:**

```bash
# Update RFC README
rex rfc update
# Output: ✓ Updated docs/rfc/README.md
```

---

## Task Commands

Task management with rich filtering, relationships, and time tracking.

### rex task create

Create a new task interactively.

**Usage:**
```bash
rex task create
```

**Interactive Prompts:**

```
Title: Implement IAM Collector
Type (core/plugin/ui/other) [core]: core
Priority (P0/P1/P2/P3) [P2]: P1
Estimated hours [8]: 12
Tags (comma-separated, optional): aws,security,collector
Assignee (optional): john@example.com
Phase (optional): 3

Relationships (optional):
Blocked by (TASK-IDs, comma-separated):
Blocks (TASK-IDs, comma-separated): TASK-010
Related to (TASK-IDs, comma-separated): TASK-005
```

**Behavior:**
- Auto-generates next TASK-ID (e.g., TASK-001, TASK-002)
- Creates file in: `docs/tasks/{type}/active/TASK-NNN-title.md`
- Sets initial status to "planned"
- Validates all input

**Task Types:**
- `core` - Core functionality
- `plugin` - Plugin-related
- `ui` - User interface
- `other` - Other/miscellaneous

**Priority Levels:**
- `P0` - Blocker (must fix immediately)
- `P1` - High priority
- `P2` - Medium priority (default)
- `P3` - Low priority

**Examples:**

```bash
# Create task (interactive)
rex task create
# Follow prompts...
# Output: ✓ Created task TASK-001: Implement IAM Collector
#         File: docs/tasks/core/active/TASK-001-implement-iam-collector.md
```

---

### rex task list

List tasks with rich filtering options.

**Usage:**
```bash
rex task list [flags]
```

**Flags:**
```bash
--type string        Filter by type (core, plugin, ui, other)
--status string      Filter by status (planned, in_progress, blocked, completed, cancelled)
--priority string    Filter by priority (P0, P1, P2, P3)
--tags string        Filter by tags (comma-separated)
--assignee string    Filter by assignee
--phase int          Filter by phase number
```

**Output Format:**
```
Tasks
═════

Statistics: 23 total | 5 planned | 10 in_progress | 2 blocked | 6 completed

ID        Title                           Type    Status         Priority  Est.  Tags
─────────────────────────────────────────────────────────────────────────────────────────────
TASK-001  Implement IAM Collector         core    in_progress    P1        12h   aws, security
TASK-002  Add Task Statistics             core    completed      P1        8h    cli, stats
TASK-003  Create Plugin System            plugin  planned        P2        20h   plugin
```

**Examples:**

```bash
# List all tasks
rex task list

# List tasks in progress
rex task list --status in_progress

# List high-priority core tasks
rex task list --type core --priority P1

# List P0 and P1 tasks (comma-separated)
rex task list --priority P0,P1

# List tasks by tag
rex task list --tags database

# List multiple tags (AND logic)
rex task list --tags database,cache

# Combine multiple filters
rex task list --type core --status in_progress --priority P1 --tags aws

# List blocked tasks
rex task list --status blocked

# List tasks by assignee
rex task list --assignee john@example.com

# List tasks in phase 3
rex task list --phase 3
```

**Task Statuses:**
- `planned` - Not started yet
- `in_progress` - Currently being worked on
- `blocked` - Blocked by another task
- `completed` - Finished
- `cancelled` - Cancelled/abandoned

---

### rex task complete

Mark a task as completed, track actual hours, and move to completed directory.

**Usage:**
```bash
rex task complete <task-id> [flags]
rex task complete TASK-001
rex task complete TASK-001 --actual-hours 10
```

**Arguments:**
- `<task-id>` - Task ID to complete (e.g., TASK-001)

**Flags:**
```bash
--actual-hours float    Actual hours spent (prompted if not provided)
```

**Behavior:**
- Updates task status to "completed"
- Sets completed_date to current date
- Records actual_hours (prompts if not provided)
- Moves file from `active/` to `completed/` directory
- Updates database cache
- Shows comparison to estimated hours

**Examples:**

```bash
# Complete task (prompted for hours)
rex task complete TASK-001
> Actual hours spent: 10
# Output: ✓ Completed task TASK-001: Implement IAM Collector
#         Actual hours: 10.0 (estimated: 12.0)
#         ✓ Under estimate by 2.0 hours
#         Moved: docs/tasks/core/completed/TASK-001-implement-iam-collector.md

# Complete with inline hours
rex task complete TASK-002 --actual-hours 9
# Output: ✓ Completed task TASK-002: Add Task Statistics
#         Actual hours: 9.0 (estimated: 8.0)
#         ⚠️  Over estimate by 1.0 hours

# Complete task that was over estimate
rex task complete TASK-003 --actual-hours 25
# Output: ✓ Completed task TASK-003: Create Plugin System
#         Actual hours: 25.0 (estimated: 20.0)
#         ⚠️  Over estimate by 5.0 hours
```

**File Movement:**
- Before: `docs/tasks/core/active/TASK-001-name.md`
- After: `docs/tasks/core/completed/TASK-001-name.md`

---

### rex task stats

Show comprehensive task statistics with progress bars.

**Usage:**
```bash
rex task stats
```

**Output:**

```
Task Statistics
═══════════════

Total Tasks: 23

Status Breakdown:
  planned     : ████░░░░░░░░░░░░░░░░ 20% (5)
  in_progress : ████████░░░░░░░░░░░░ 40% (10)
  blocked     : ██░░░░░░░░░░░░░░░░░░  8% (2)
  completed   : ████████████░░░░░░░░ 60% (6)
  cancelled   : ░░░░░░░░░░░░░░░░░░░░  0% (0)

Type Breakdown:
  core        : ████████████████░░░░ 80% (18)
  plugin      : ████░░░░░░░░░░░░░░░░ 20% (5)
  ui          : ░░░░░░░░░░░░░░░░░░░░  0% (0)
  other       : ░░░░░░░░░░░░░░░░░░░░  0% (0)

Priority Breakdown:
  P0 (blocker): ██░░░░░░░░░░░░░░░░░░  8% (2)
  P1 (high)   : ████████░░░░░░░░░░░░ 40% (10)
  P2 (medium) : ██████░░░░░░░░░░░░░░ 30% (7)
  P3 (low)    : ████░░░░░░░░░░░░░░░░ 22% (4)

Time Tracking:
  Total Estimated: 120.0 hours
  Total Actual   : 95.5 hours (6 completed tasks)
  Average per task: 15.9 hours actual

  ✓ Under estimate: 3 tasks
  ⚠️  Over estimate : 2 tasks
  ✓ On target     : 1 task
```

**Examples:**

```bash
# Show all statistics
rex task stats
```

**What's Included:**
- Total task count
- Status breakdown with progress bars
- Type breakdown
- Priority breakdown
- Time tracking summary
- Estimate accuracy

---

## Plan Commands

Plans document roadmaps, sprint plans, and quarterly objectives.

### rex plan create

Create a new plan from template.

**Usage:**
```bash
rex plan create <title>
rex plan create "Q4 2025 Roadmap"
```

**Arguments:**
- `<title>` - Title of the plan (required, use quotes for multi-word titles)

**Behavior:**
- Auto-generates next plan number (e.g., 0001, 0002)
- Creates file: `docs/plans/NNNN-title-slug.md`
- Sets status to "draft"

**Template Structure:**
```markdown
---
id: 0001
title: Q4 2025 Roadmap
status: draft
date: 2025-11-10
---

# Q4 2025 Roadmap

## Overview

High-level summary of the plan.

## Goals

Key objectives for this period.

## Timeline

Breakdown by week/month/quarter.

## Tasks

Related tasks and deliverables.
```

**Examples:**

```bash
# Create roadmap
rex plan create "Q4 2025 Roadmap"
# Output: ✓ Created plan 0001: Q4 2025 Roadmap

# Create sprint plan
rex plan create "Sprint 42 - Authentication"
# Output: ✓ Created plan 0002: Sprint 42 - Authentication
```

---

### rex plan list

List all plans with optional filtering.

**Usage:**
```bash
rex plan list [flags]
```

**Flags:**
```bash
--status string    Filter by status (draft, active, completed, cancelled)
```

**Output Format:**
```
Plans
═════

ID    Title                    Status      Date
─────────────────────────────────────────────────────────
0001  Q4 2025 Roadmap          active      2025-11-10
0002  Sprint 42                completed   2025-11-05
```

**Examples:**

```bash
# List all plans
rex plan list

# List active plans
rex plan list --status active

# List completed plans
rex plan list --status completed
```

**Status Values:**
- `draft` - Planning stage
- `active` - Currently executing
- `completed` - Finished
- `cancelled` - Abandoned

---

### rex plan update

Regenerate `docs/plans/README.md` with current plan list.

**Usage:**
```bash
rex plan update
```

**Behavior:**
- Scans all plan files in `docs/plans/`
- Generates table of plans grouped by status
- Preserves content outside auto-generated markers

**Examples:**

```bash
# Update plan README
rex plan update
# Output: ✓ Updated docs/plans/README.md
```

---

## Database Commands

Manage the SQLite cache database.

### rex rebuild

Manually rebuild the SQLite cache from markdown files.

**Usage:**
```bash
rex rebuild
```

**Behavior:**
- Scans all markdown files in:
  - `docs/adr/`
  - `docs/rfc/`
  - `docs/tasks/`
  - `docs/plans/`
- Parses YAML frontmatter
- Rebuilds `.rex.db` SQLite cache
- Shows progress and statistics

**When to Use:**
- After manually editing many markdown files
- After git pull with many document changes
- If cache appears stale or corrupted
- After bulk file operations

**Output:**

```
⏳ Rebuilding cache from markdown files...

Parsing files:
  ADRs  : 6 files
  RFCs  : 2 files
  Tasks : 23 files
  Plans : 3 files

✓ Rebuilt cache from markdown files

Document counts:
  ADRs  : 6
  RFCs  : 2
  Tasks : 23
  Plans : 3

Total documents: 34
Cache rebuilt in 45ms
```

**Examples:**

```bash
# Rebuild cache
rex rebuild

# After git pull
git pull origin main
rex rebuild
```

**Note:** Rex automatically detects stale cache and rebuilds when needed. Manual rebuild is rarely necessary.

---

### rex db info

Show comprehensive database cache information.

**Usage:**
```bash
rex db info
```

**Output:**

```
Database Cache Information
==========================

Path: /Users/you/project/.rex.db
Size: 128.00 KB
Modified: 2025-11-10 13:00:18

Last Rebuild: 2025-11-10 13:00:18
Time Since: 5m ago

Schema Version: 1

Document Counts:
  ADRs  : 6
  RFCs  : 2
  Tasks : 23
  Plans : 3

Total Documents: 34

Task Breakdown:
  By Status:
    planned     : 5
    in_progress : 10
    blocked     : 2
    completed   : 6

  By Type:
    core        : 18
    plugin      : 5

  By Priority:
    P0          : 2
    P1          : 10
    P2          : 7
    P3          : 4
```

**Examples:**

```bash
# Show database info
rex db info
```

**What's Included:**
- Database file path and size
- Last modification time
- Last rebuild timestamp
- Schema version
- Document counts by type
- Task breakdowns by status, type, priority

---

## Utility Commands

### rex version

Show version information.

**Usage:**
```bash
rex version
```

**Output:**

```
Rex Version: v1.0.0
Commit: abc1234
Built: 2025-11-10T10:30:00Z
Go Version: go1.23.2
Platform: darwin/amd64
```

**Examples:**

```bash
# Show version
rex version

# Check if rex is installed
rex version
```

---

### rex help

Show help for rex or any command.

**Usage:**
```bash
rex --help
rex help
rex <command> --help
rex help <command>
```

**Examples:**

```bash
# Show main help
rex --help
rex help

# Show command help
rex task --help
rex help task

# Show subcommand help
rex task create --help
rex help task create
```

---

## Common Workflows

### Morning Standup

```bash
# Check tasks in progress
rex task list --status in_progress

# Check blocked tasks
rex task list --status blocked

# View high-priority items
rex task list --priority P0,P1
```

### Starting a New Task

```bash
# Create task
rex task create
# Follow interactive prompts...

# Verify it was created
rex task list --status planned
```

### Completing Tasks

```bash
# Complete with time tracking
rex task complete TASK-001 --actual-hours 8

# View updated stats
rex task stats
```

### Sprint Planning

```bash
# Create sprint plan
rex plan create "Sprint 42"

# Create tasks for sprint
rex task create  # Repeat for each task

# List planned work
rex task list --status planned --phase 42

# Track progress
rex task stats
```

### Making Architectural Decisions

```bash
# Create ADR
rex adr create "Use PostgreSQL for Analytics"

# Edit ADR (change status to accepted)
$EDITOR docs/adr/0001-use-postgresql.md

# Update README
rex adr update

# List accepted decisions
rex adr list --status accepted
```

### Proposing a New Feature

```bash
# Create RFC
rex rfc create "Add GraphQL API"

# Share for review (change status to review)
$EDITOR docs/rfc/0001-add-graphql-api.md

# List RFCs in review
rex rfc list --status review

# After approval, update README
rex rfc update
```

---

## Tips and Best Practices

### 1. Use Quotes for Multi-Word Titles

```bash
# Correct
rex adr create "Use SQLite as Cache"

# Incorrect (only captures "Use")
rex adr create Use SQLite as Cache
```

### 2. Combine Filters for Precision

```bash
# Find exactly what you need
rex task list --type core --status in_progress --priority P1 --tags database
```

### 3. Track Time Accurately

```bash
# Always provide actual hours when completing
rex task complete TASK-001 --actual-hours 8.5
```

### 4. Use Tags Effectively

```bash
# Create tasks with tags
rex task create
> Tags: aws,security,iam,collector

# Search by tags later
rex task list --tags aws
rex task list --tags security,iam
```

### 5. Review Stats Regularly

```bash
# Daily standup
rex task stats

# Weekly review
rex task list --status completed
rex db info
```

### 6. Keep Cache Fresh

```bash
# After git operations
git pull
rex rebuild  # Usually automatic, but can be manual

# Check cache status
rex db info
```

---

## Exit Codes

Rex uses standard exit codes:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments or flags |
| 3 | File not found |
| 4 | Database error |
| 5 | Parsing error |

**Examples:**

```bash
# Check exit code
rex task list
echo $?  # 0 on success

# Handle errors in scripts
if ! rex task complete TASK-001; then
    echo "Failed to complete task"
    exit 1
fi
```

---

## Environment Variables

Rex respects these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `REX_CONFIG` | Path to config file | `$HOME/.rex.yaml` |
| `EDITOR` | Editor for file editing | `vim` |
| `PAGER` | Pager for long output | `less` |

**Examples:**

```bash
# Use custom config
export REX_CONFIG=/path/to/.rex.yaml
rex task list

# Use specific editor
export EDITOR=nano
rex adr create "New Decision"
```

---

## Further Reading

- [README.md](../README.md) - Quick start and features
- [MIGRATION.md](MIGRATION.md) - Migrate from bash scripts
- [IMPLEMENTATION_GUIDE.md](IMPLEMENTATION_GUIDE.md) - Technical details
- [MVP_GUIDE.md](MVP_GUIDE.md) - MVP scope and phases

---

**Questions?** Open an issue: [rex/issues](https://github.com/donaldgifford/rex/issues)
