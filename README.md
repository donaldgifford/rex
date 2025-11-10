# Rex

> Fast, SQLite-backed documentation management CLI

Rex is a modern documentation management tool that helps you organize Architecture Decision Records (ADRs), RFCs, Tasks, and Plans in your project. All documents are stored as git-friendly markdown files with YAML frontmatter, and cached in SQLite for lightning-fast queries.

[![Go Version](https://img.shields.io/github/go-mod/go-version/donaldgifford/rex)](https://golang.org/)
[![License](https://img.shields.io/github/license/donaldgifford/rex)](LICENSE)
[![Release](https://img.shields.io/github/v/release/donaldgifford/rex)](https://github.com/donaldgifford/rex/releases)

## ✨ Features

- 📝 **Manage 4 Document Types**: ADRs, RFCs, Tasks, Plans
- ⚡ **Fast**: SQLite cache provides < 100ms queries even with 100+ documents
- 🔍 **Rich Filtering**: Filter by type, status, priority, tags, and more
- 📊 **Statistics**: Comprehensive stats and progress tracking
- 🔗 **Relationships**: Track task dependencies (blocked_by, blocks, related_to)
- ⏱️ **Time Tracking**: Estimate vs actual hours for tasks
- 🚀 **Zero Dependencies**: Pure Go, no external runtime requirements
- 🌍 **Cross-Platform**: Linux, macOS, Windows binaries
- 📦 **Git-Friendly**: Plain markdown files with clear diffs
- 🔄 **Auto-Rebuild**: Cache automatically syncs with file changes

## 📦 Installation

### Option 1: Go Install (Recommended)

```bash
go install github.com/donaldgifford/rex@latest
```

### Option 2: Download Binary

Download the latest release for your platform from the [Releases page](https://github.com/donaldgifford/rex/releases).

**Linux/macOS:**
```bash
# Download and extract
tar -xzf rex_*_linux_amd64.tar.gz
sudo mv rex /usr/local/bin/

# Verify installation
rex version
```

**Windows:**
Download the `.zip` file, extract, and add the directory to your `PATH`.

### Option 3: Build from Source

```bash
git clone https://github.com/donaldgifford/rex.git
cd rex
go build -o rex main.go
```

## 🚀 Quick Start

### Initialize Your Project

Rex works with any project that has a `docs/` directory:

```bash
mkdir -p docs/{adr,rfc,tasks,plans}
rex rebuild
```

### Create Your First Documents

```bash
# Create an Architecture Decision Record
rex adr create "Use SQLite as Cache Layer"

# Create an RFC
rex rfc create "Add Plugin System"

# Create a task (interactive prompts)
rex task create

# Create a plan
rex plan create "Q4 2025 Roadmap"
```

### List and Filter

```bash
# List all documents
rex adr list
rex rfc list
rex task list
rex plan list

# Filter by status
rex adr list --status accepted
rex rfc list --status draft
rex task list --status in_progress

# Filter tasks by multiple criteria
rex task list --type core --priority P1 --tags database
```

### View Statistics

```bash
# Task statistics with progress bars
rex task stats

# Database cache information
rex db info
```

## 📚 Documentation Types

### ADRs (Architecture Decision Records)

Document significant architectural decisions with context, consequences, and alternatives.

**Statuses:** `draft`, `accepted`, `deprecated`, `superseded`

```bash
# Create, list, update
rex adr create "Use PostgreSQL for Analytics"
rex adr list --status accepted
rex adr update  # Generate README
```

### RFCs (Requests for Comments)

Propose features or changes that need discussion and approval.

**Statuses:** `draft`, `review`, `approved`, `implemented`, `rejected`

```bash
# Create with author
rex rfc create "Add GraphQL API"
rex rfc list --status review
rex rfc update  # Generate README
```

### Tasks

Manage project tasks with relationships, time tracking, and rich metadata.

**Types:** `core`, `plugin`, `ui`, `other`
**Statuses:** `planned`, `in_progress`, `blocked`, `completed`, `cancelled`
**Priorities:** `P0` (blocker), `P1` (high), `P2` (medium), `P3` (low)

```bash
# Create task interactively
rex task create

# List and filter
rex task list --type core --status in_progress
rex task list --priority P0,P1
rex task list --tags database,cache

# Complete a task (moves to completed/)
rex task complete TASK-001

# View statistics
rex task stats
```

**Task Relationships:**
- `blocked_by`: Tasks blocking this one
- `blocks`: Tasks this one blocks
- `related_to`: Related tasks

### Plans

Document roadmaps, sprint plans, release plans, and quarterly objectives.

**Statuses:** `draft`, `active`, `completed`, `cancelled`

```bash
# Create, list, update
rex plan create "Q1 2025 Roadmap"
rex plan list --status active
rex plan update  # Generate README
```

## 💡 Common Workflows

### Daily Development Workflow

```bash
# Morning: Check your tasks
rex task list --status in_progress

# Create a new task for today's work
rex task create

# Filter high-priority items
rex task list --priority P0,P1

# Complete a task when done
rex task complete TASK-123

# Evening: View your progress
rex task stats
```

### Sprint Planning

```bash
# Create sprint plan
rex plan create "Sprint 42 - Authentication"

# Create tasks for the sprint
rex task create  # Repeat for each task

# List all planned work
rex task list --status planned --phase 42

# Track progress
rex task stats
```

### Architectural Decision Tracking

```bash
# Propose a decision
rex adr create "Use Redis for Caching"

# Update status after approval (edit markdown file)
# Then regenerate README
rex adr update

# List all accepted decisions
rex adr list --status accepted
```

## 🗄️ How It Works

### Markdown Files (Source of Truth)

All documents are plain markdown files with YAML frontmatter:

```markdown
---
id: TASK-001
title: Implement SQLite cache
type: core
status: in_progress
priority: P1
estimated_hours: 8
tags: [database, cache]
---

# Implement SQLite cache

Description of the task...
```

### SQLite Cache (Fast Queries)

Rex maintains a `.rex.db` SQLite cache for instant queries:

- **Auto-Sync**: Automatically rebuilds when files change
- **Fast**: Queries return in < 100ms even with 100+ documents
- **Relationships**: SQL JOINs for task dependencies
- **Aggregations**: GROUP BY for statistics

The cache is automatically managed - you rarely need to think about it.

### Directory Structure

```
your-project/
├── docs/
│   ├── adr/                    # Architecture Decision Records
│   │   ├── 0001-use-sqlite.md
│   │   ├── 0002-go-embed.md
│   │   └── README.md           # Auto-generated
│   ├── rfc/                    # RFCs
│   │   ├── 0001-plugin-system.md
│   │   └── README.md           # Auto-generated
│   ├── tasks/                  # Tasks
│   │   ├── core/
│   │   │   ├── active/
│   │   │   │   └── TASK-001-implement-cache.md
│   │   │   └── completed/
│   │   │       └── TASK-002-add-tests.md
│   │   ├── plugin/
│   │   ├── ui/
│   │   └── other/
│   └── plans/                  # Plans/Roadmaps
│       ├── 0001-q4-roadmap.md
│       └── README.md           # Auto-generated
└── .rex.db                     # SQLite cache (gitignored)
```

## 🔧 Configuration

Rex uses sensible defaults and works out-of-the-box. Optional configuration via `$HOME/.rex.yaml`:

```yaml
# Coming soon - currently all settings use defaults
docs_dir: docs
```

## 🛠️ Advanced Usage

### Manual Cache Rebuild

```bash
rex rebuild
```

### Cache Information

```bash
rex db info
```

**Output:**
```
Database Cache Information
==========================

Path: /path/to/.rex.db
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
```

### Task Completion with Time Tracking

```bash
# Complete with inline hours
rex task complete TASK-001 --actual-hours 10

# Or be prompted
rex task complete TASK-001
> Actual hours spent: 10

✓ Completed task TASK-001: Implement cache
  Actual hours: 10.0 (estimated: 8.0)
  ⚠️  Over estimate by 2.0 hours
```

## 📖 Command Reference

### Core Commands

| Command | Description |
|---------|-------------|
| `rex adr create <title>` | Create a new ADR |
| `rex adr list [--status]` | List ADRs with optional filtering |
| `rex adr update` | Regenerate ADR README |
| `rex rfc create <title>` | Create a new RFC |
| `rex rfc list [--status]` | List RFCs with optional filtering |
| `rex rfc update` | Regenerate RFC README |
| `rex task create` | Create a new task (interactive) |
| `rex task list [filters]` | List tasks with rich filtering |
| `rex task complete <id>` | Mark task complete and move to completed/ |
| `rex task stats` | View task statistics |
| `rex plan create <title>` | Create a new plan |
| `rex plan list [--status]` | List plans with optional filtering |
| `rex plan update` | Regenerate plan README |
| `rex rebuild` | Manually rebuild SQLite cache |
| `rex version` | Show version information |
| `rex db info` | Show database statistics |

### Task Filtering Options

```bash
--type       Filter by type (core, plugin, ui, other)
--status     Filter by status (planned, in_progress, blocked, completed, cancelled)
--priority   Filter by priority (P0, P1, P2, P3)
--tags       Filter by tags (comma-separated)
```

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Rex is inspired by:
- [adr-tools](https://github.com/npryce/adr-tools) - The original ADR tool
- [Linear](https://linear.app) - Task management excellence
- [SQLite](https://sqlite.org) - The best embedded database

## 🔗 Links

- [Documentation](docs/)
- [Releases](https://github.com/donaldgifford/rex/releases)
- [Issues](https://github.com/donaldgifford/rex/issues)
- [Discussions](https://github.com/donaldgifford/rex/discussions)

---

**Made with ❤️ by [Donald Gifford](https://github.com/donaldgifford)**
