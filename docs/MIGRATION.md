# Migration Guide: Bash Scripts to Rex CLI

This guide helps you migrate from the old bash scripts and Makefile commands to the new `rex` CLI tool.

## Overview

Rex CLI replaces the bash scripts in `tools/` with a fast, SQLite-backed CLI tool written in Go. The new CLI provides:

- **Single Binary**: No dependencies on bash, yq, mise, or external tools
- **Fast Queries**: SQLite cache provides instant results (< 100ms even with 100+ documents)
- **Rich Filtering**: Filter tasks by type, status, priority, and tags
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Better UX**: Colored output, progress bars, and interactive prompts
- **Time Tracking**: Track estimated vs actual hours for tasks
- **Statistics**: Comprehensive stats and progress tracking

## Quick Command Reference

| Old Command | New Command | Notes |
|------------|-------------|-------|
| `make adr "Title"` | `rex adr create "Title"` | Creates ADR |
| `make adr-update` | `rex adr update` | Updates ADR README |
| `make rfc "Title"` | `rex rfc create "Title"` | Creates RFC |
| `make rfc-update` | `rex rfc update` | Updates RFC README |
| `make task` | `rex task create` | Interactive task creation |
| `make task-list` | `rex task list` | Lists all tasks |
| `make task-stats` | `rex task stats` | Task statistics |
| `make task-complete <id>` | `rex task complete <id>` | Completes task |
| `make task-update` | `rex task update` | Updates task READMEs |
| - | `rex plan create "Title"` | New: Create plans/roadmaps |
| - | `rex plan list` | New: List plans |
| - | `rex plan update` | New: Update plan README |
| - | `rex rebuild` | Manually rebuild SQLite cache |
| - | `rex version` | Show version information |
| - | `rex db info` | Show database statistics |

## Step-by-Step Migration

### 1. Install Rex CLI

```bash
# Option 1: Install from source (recommended during development)
cd /path/to/rex
go build -o rex main.go
sudo mv rex /usr/local/bin/

# Option 2: Install via go install (once released)
go install github.com/donaldgifford/rex@latest

# Verify installation
rex version
```

### 2. Build SQLite Cache

The first time you use rex, build the SQLite cache from your existing markdown files:

```bash
cd /path/to/your-project
rex rebuild
```

This command:
- Scans all markdown files in `docs/adr/`, `docs/rfc/`, `docs/tasks/`, `docs/plans/`
- Extracts YAML frontmatter
- Builds `.rex.db` SQLite cache
- Auto-ignores `.rex.db` in `.gitignore`

**Output:**
```
⏳ Rebuilding cache from markdown files...

✓ Rebuilt cache from markdown files

Document counts:
  ADRs  : 6
  RFCs  : 2
  Tasks : 23
  Plans : 3

Total documents: 34
Cache rebuilt in 45ms
```

### 3. Test New Commands

Verify rex can read your existing documents:

```bash
# List ADRs
rex adr list

# List tasks
rex task list

# View statistics
rex task stats

# Show database info
rex db info
```

### 4. Migrate Your Workflow

**Before (bash scripts):**
```bash
# Morning standup - check tasks
make task-list

# Create new task
make task

# Complete task
make task-complete TASK-001

# View stats
make task-stats
```

**After (rex CLI):**
```bash
# Morning standup - check tasks
rex task list --status in_progress

# Create new task
rex task create

# Complete task with time tracking
rex task complete TASK-001 --actual-hours 6

# View stats with progress bars
rex task stats
```

### 5. Update Scripts and CI/CD

If you have scripts or CI/CD pipelines that use the old commands, update them:

**Before:**
```bash
#!/bin/bash
make adr "Use PostgreSQL"
make task
make task-stats
```

**After:**
```bash
#!/bin/bash
rex adr create "Use PostgreSQL"
rex task create
rex task stats
```

### 6. Remove Old Bash Scripts (Optional)

Once you've verified rex works for your workflow, you can optionally remove the old bash scripts:

```bash
# Keep for reference initially
git mv tools/ tools.old/

# Or remove entirely after verification
rm -rf tools/
```

**Note:** You can keep both systems running in parallel during migration. Rex and the bash scripts operate independently.

## Handling Existing Files

### Files Without YAML Frontmatter

If you have existing documents without YAML frontmatter, rex will skip them during cache rebuild. You have two options:

#### Option 1: Add Frontmatter Manually

Edit each file and add YAML frontmatter:

```markdown
---
id: TASK-001
title: Your Task Title
type: core
status: in_progress
priority: P1
estimated_hours: 8
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

# Your Task Title

Content here...
```

#### Option 2: Use Migration Script (If Available)

If you have the old migration script:

```bash
# This was in tools/scripts/tasks/migrate-to-frontmatter.sh
./tools/scripts/tasks/migrate-to-frontmatter.sh
```

Then rebuild the cache:

```bash
rex rebuild
```

### Files in Wrong Locations

Rex expects specific directory structures:

```
docs/
├── adr/
│   └── 0001-your-adr.md
├── rfc/
│   └── 0001-your-rfc.md
├── tasks/
│   ├── core/
│   │   ├── active/
│   │   │   └── TASK-001-your-task.md
│   │   └── completed/
│   │       └── TASK-002-done-task.md
│   ├── plugin/
│   ├── ui/
│   └── other/
└── plans/
    └── 0001-your-plan.md
```

If your files are in different locations, either:
1. Move them to the expected structure
2. Update file paths in the database after running `rex rebuild`

## What's New in Rex CLI

### 1. Rich Task Filtering

```bash
# Filter by multiple criteria
rex task list --type core --status in_progress --priority P1

# Filter by tags
rex task list --tags database,cache

# Combine filters
rex task list --type core,plugin --status planned,in_progress
```

### 2. Time Tracking

```bash
# Complete task with actual hours
rex task complete TASK-001 --actual-hours 10

# View time tracking stats
rex task stats

# Output includes:
# Time Tracking:
#   Total Estimated: 120.0 hours
#   Total Actual   : 95.5 hours
#   ⚠️  Over estimate: 5 tasks
```

### 3. Task Relationships

Rex fully supports task relationships from YAML frontmatter:

```yaml
blocked_by: [TASK-005, TASK-010]
blocks: [TASK-020]
related_to: [TASK-003]
```

These are stored in the database and can be queried.

### 4. Progress Statistics

```bash
rex task stats

# Output includes progress bars:
# Status Breakdown:
#   planned     : ████░░░░░░░░░░░░░░░░ 20% (5)
#   in_progress : ████████░░░░░░░░░░░░ 40% (10)
#   completed   : ████████████░░░░░░░░ 60% (15)
```

### 5. Database Cache Info

```bash
rex db info

# Output:
# Database Cache Information
# ==========================
# Path: /path/to/.rex.db
# Size: 128.00 KB
# Modified: 2025-11-10 13:00:18
# Last Rebuild: 2025-11-10 13:00:18
# Document Counts:
#   ADRs  : 6
#   RFCs  : 2
#   Tasks : 23
#   Plans : 3
```

### 6. Plan Management

New document type for roadmaps and sprint plans:

```bash
# Create plan
rex plan create "Q4 2025 Roadmap"

# List active plans
rex plan list --status active

# Update plan README
rex plan update
```

### 7. Auto-Rebuild Cache

Rex automatically detects when markdown files are newer than the cache:

```
⚠️  Cache is stale (files modified since last rebuild)
⏳ Rebuilding cache from markdown files...
✓ Rebuilt cache
```

You rarely need to manually run `rex rebuild`.

## Feature Comparison

| Feature | Bash Scripts | Rex CLI |
|---------|-------------|---------|
| **Speed** | Slow (parses all files) | Fast (SQLite cache < 100ms) |
| **Dependencies** | bash, yq, mise | None (single binary) |
| **Filtering** | Limited | Rich (type, status, priority, tags) |
| **Time Tracking** | Manual editing | Built-in prompts & stats |
| **Statistics** | Basic counts | Progress bars, breakdowns |
| **Cross-Platform** | Unix only | Linux, macOS, Windows |
| **Auto-Complete** | No | Yes (with shell completion) |
| **Version Info** | No | Yes (`rex version`) |
| **Database Info** | No | Yes (`rex db info`) |

## Deprecation Timeline

The bash scripts in `tools/` are **immediately deprecated** but will remain in the repository for reference during migration.

**Recommended Timeline:**
- **Week 1-2**: Run rex and bash scripts in parallel, verify behavior
- **Week 3-4**: Switch entirely to rex CLI
- **Month 2**: Remove `tools/` directory from repository

## Troubleshooting

### Cache is Empty After Rebuild

**Problem:** `rex task list` shows no results after `rex rebuild`

**Solutions:**
1. Check if task files have YAML frontmatter:
   ```bash
   head -20 docs/tasks/core/active/TASK-001-*.md
   ```

2. Verify file locations match expected structure:
   ```bash
   find docs/tasks -name "*.md"
   ```

3. Check for YAML parsing errors during rebuild:
   ```bash
   rex rebuild 2>&1 | grep -i error
   ```

### Database is Locked

**Problem:** `Error: database is locked`

**Solution:**
- Rex uses WAL mode which supports concurrent reads
- Only one write operation can occur at a time
- If you see this error, another rex process is running
- Wait a few seconds and retry

### Files Not Showing Up

**Problem:** Created task file manually but `rex task list` doesn't show it

**Solutions:**
1. Ensure YAML frontmatter is valid:
   ```yaml
   ---
   id: TASK-001
   title: Task Title
   type: core
   status: planned
   ---
   ```

2. Rebuild cache:
   ```bash
   rex rebuild
   ```

3. Check database info to verify document count:
   ```bash
   rex db info
   ```

### Migration from Old Numbering

**Problem:** Old tasks used different ID format

**Solution:**
- Rex expects `TASK-NNN` format (e.g., `TASK-001`, `TASK-042`)
- If your old tasks used a different format, rename files:
  ```bash
  # Example: Rename TASK_001 to TASK-001
  for f in docs/tasks/*/active/TASK_*.md; do
    mv "$f" "${f//TASK_/TASK-}"
  done
  ```

Then rebuild:
```bash
rex rebuild
```

## Getting Help

If you encounter issues during migration:

1. **Check Documentation:**
   - Main README: [README.md](../README.md)
   - Implementation Guide: [IMPLEMENTATION_GUIDE.md](IMPLEMENTATION_GUIDE.md)
   - MVP Guide: [MVP_GUIDE.md](MVP_GUIDE.md)

2. **View Command Help:**
   ```bash
   rex --help
   rex task --help
   rex task create --help
   ```

3. **Inspect Database:**
   ```bash
   rex db info
   sqlite3 .rex.db "SELECT * FROM tasks LIMIT 5;"
   ```

4. **Report Issues:**
   - GitHub Issues: [rex/issues](https://github.com/donaldgifford/rex/issues)

## Next Steps

After migrating to rex CLI:

1. **Learn Advanced Features:**
   ```bash
   # Rich filtering
   rex task list --type core --priority P0,P1 --status in_progress

   # Time tracking
   rex task complete TASK-001 --actual-hours 8

   # Statistics
   rex task stats
   ```

2. **Set Up Aliases (Optional):**
   ```bash
   # Add to ~/.bashrc or ~/.zshrc
   alias rt='rex task'
   alias rtl='rex task list'
   alias rts='rex task stats'
   alias rtc='rex task create'
   ```

3. **Integrate with Editor:**
   ```bash
   # Vim: Add to ~/.vimrc
   command! RexTaskList !rex task list

   # VSCode: Add to tasks.json
   {
     "label": "Rex: List Tasks",
     "type": "shell",
     "command": "rex task list"
   }
   ```

4. **Share with Team:**
   - Update team documentation
   - Add rex installation to onboarding
   - Update CI/CD pipelines

---

**Welcome to Rex CLI! 🦖**

Questions? See [README.md](../README.md) or open an issue on GitHub.
