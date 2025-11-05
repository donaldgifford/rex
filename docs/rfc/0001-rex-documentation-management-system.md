# RFC 0001: Rex Documentation Management System

**Status**: Draft
**Author**: Claude Code
**Created**: 2025-11-05
**Last Updated**: 2025-11-05

## Executive Summary

Rex is currently a collection of bash scripts and makefiles that manage documentation files (RFCs, ADRs, Plans, Tasks) in a structured markdown format. While functional, this approach requires copying scripts and makefiles to each new repository, making setup non-trivial and maintenance difficult across multiple projects.

This RFC proposes refactoring Rex into a single Go CLI tool with embedded templates that can be installed once and used across all repositories. The CLI will standardize documentation workflows, eliminate the need for repository-specific tooling, and make it trivial to initialize new projects with proper documentation structure.

The core insight is simple: Rex is fundamentally about creating and managing markdown files in opinionated formats. By moving this logic into a distributable binary with Go's template embedding and using SQLite as a query cache, we gain portability, consistency, fast queries, and ease of use while maintaining the simplicity of markdown-based documentation.

**Key architectural decision:** Markdown files remain the source of truth (git-friendly), while a local SQLite database (`.rex.db`, gitignored) serves as a fast query cache. The database auto-rebuilds from markdown when stale or missing, eliminating the need to parse 100+ files for every `rex task list` command.

## Problem Statement

### Current State Analysis

Rex currently exists as:
- **Bash scripts** in `tools/docs/` for creating RFCs and ADRs
- **Bash scripts** in `tools/scripts/tasks/` for task management (create, complete, update, list, stats)
- **Makefiles** in `tools/makefiles/` for orchestration (docs.mk, go.mk, common.mk)
- **Template files** stored as `.md` files in each docs directory (docs/adr/template.md, docs/rfc/template.md, etc.)
- **Update scripts** that parse markdown and regenerate README tables

To use Rex in a new repository, you must:
1. Copy the entire `tools/` directory structure
2. Copy the root `Makefile` and configure includes
3. Create the `docs/` directory structure (adr/, rfc/, tasks/, plans/)
4. Copy template files to each docs directory
5. Ensure dependencies are installed (bash, yq, sed, awk, etc.)

### Pain Points

**Distribution and Setup:**
- Must copy ~14 bash scripts and 3 makefiles to each new repo
- No version management—updates require manual copying
- Setup is manual and error-prone
- No easy way to update Rex tooling across existing repos

**Dependency Management:**
- Scripts depend on external tools (yq, sed, awk) with varying versions
- mise.toml required for yq installation
- Different shells (bash) may behave differently across environments

**Consistency:**
- Each repo can diverge in tooling as scripts are modified locally
- No guarantee that RFC/ADR/Task workflows are consistent across projects
- Templates can be edited per-repo, breaking standardization

**Developer Experience:**
- Requires understanding of bash scripting to modify or debug
- Make targets are opaque—users must read Makefiles to understand behavior
- Error messages are inconsistent across scripts
- No built-in help or command discovery

### Evidence of Complexity

**Current tool distribution:**
```
tools/
├── docs/
│   ├── create-adr.sh          (68 lines)
│   ├── create-rfc.sh          (69 lines)
│   ├── update-adr-readme.sh   (111 lines)
│   ├── update-rfc-readme.sh   (95 lines)
│   └── labels.sh              (283 lines)
├── scripts/
│   └── tasks/
│       ├── task-create.sh           (295 lines)
│       ├── task-complete.sh         (156 lines)
│       ├── task-update.sh           (488 lines)
│       ├── task-list.sh             (145 lines)
│       ├── task-stats.sh            (260 lines)
│       └── migrate-to-frontmatter.sh (289 lines)
└── makefiles/
    ├── common.mk             (80 lines)
    ├── docs.mk               (48 lines)
    └── go.mk                 (55 lines)

Total: ~2,442 lines of bash/make to copy per repository
```

**Setup friction:**
- New developer setting up a repo: 15-20 minutes to copy and configure
- Updating Rex across 5 repos: 30+ minutes of manual copying and testing

## Proposed Solution

### Vision Statement

Refactor Rex into a single, distributable Go CLI binary that embeds all templates and logic for managing documentation files. Developers install Rex once (via `go install` or binary download) and use it across all repositories without copying scripts or makefiles.

The CLI provides intuitive commands like `rex init`, `rex adr create`, `rex task list`, etc., with built-in help, validation, and error messages. Templates are embedded in the binary using Go's `embed` package, ensuring consistency and version control. The tool remains focused on its core mission: creating and managing opinionated markdown files for documentation.

### Architecture Philosophy

**Core Principles:**

1. **Single Binary Distribution** - One `rex` binary installed globally, no per-repo tooling
2. **Embedded Templates** - All templates bundled in the binary using Go embed, ensuring version consistency
3. **Opinionated Defaults** - Strong conventions for file structure, naming, and formats with escape hatches for customization
4. **Markdown First** - Files remain plain markdown with YAML frontmatter—no proprietary formats (source of truth)
5. **SQLite Cache** - Local `.rex.db` for fast queries; auto-rebuilt from markdown when stale/missing (gitignored)
6. **Cross-Repository Consistency** - Same version of Rex produces identical workflows across all projects

### High-Level Architecture

```
rex (Go CLI Binary)
├── cmd/
│   ├── root.go           - Root command and global flags
│   ├── init.go           - Initialize new repo with docs/ structure
│   ├── rebuild.go        - Rebuild .rex.db from markdown files
│   ├── adr/
│   │   ├── create.go     - Create new ADR (writes .md + updates DB)
│   │   ├── update.go     - Update ADR README
│   │   └── list.go       - List ADRs (queries DB)
│   ├── rfc/
│   │   ├── create.go     - Create new RFC (writes .md + updates DB)
│   │   ├── update.go     - Update RFC README
│   │   └── list.go       - List RFCs (queries DB)
│   ├── task/
│   │   ├── create.go     - Create new task (writes .md + updates DB)
│   │   ├── complete.go   - Mark task complete (updates .md + DB)
│   │   ├── list.go       - List/filter tasks (queries DB)
│   │   ├── stats.go      - Show task statistics (queries DB)
│   │   └── update.go     - Update task READMEs
│   └── plan/
│       ├── create.go     - Create new plan (writes .md + updates DB)
│       └── update.go     - Update plan README
├── internal/
│   ├── db/
│   │   ├── db.go         - SQLite connection management
│   │   ├── schema.go     - Database schema and migrations
│   │   ├── adr.go        - ADR table operations
│   │   ├── rfc.go        - RFC table operations
│   │   ├── task.go       - Task table operations
│   │   ├── plan.go       - Plan table operations
│   │   └── rebuild.go    - Scan docs/ and rebuild DB
│   ├── templates/
│   │   ├── embed.go      - Embed all .md templates
│   │   ├── adr.md
│   │   ├── rfc.md
│   │   ├── task.md
│   │   ├── plan.md
│   │   ├── adr-readme.md
│   │   ├── rfc-readme.md
│   │   └── task-readme.md
│   ├── parser/
│   │   ├── frontmatter.go  - Parse YAML frontmatter
│   │   └── markdown.go     - Parse markdown structure
│   ├── generator/
│   │   ├── file.go         - File creation logic
│   │   ├── slug.go         - Title -> filename slug
│   │   └── id.go           - Auto-increment IDs
│   └── config/
│       └── rex.go          - Config file (.rex.yaml) parsing
└── templates/              - Embedded directory
    ├── adr/
    │   ├── template.md
    │   └── README.md
    ├── rfc/
    │   ├── template.md
    │   └── README.md
    ├── tasks/
    │   ├── template.md
    │   └── README.md
    └── plans/
        ├── template.md
        └── README.md

Repository Structure:
├── .rex.db               - SQLite cache (gitignored)
├── .rex.yaml             - Optional config
└── docs/
    ├── adr/              - ADR markdown files (source of truth)
    ├── rfc/              - RFC markdown files (source of truth)
    ├── tasks/            - Task markdown files (source of truth)
    └── plans/            - Plan markdown files (source of truth)
```

**Key Components:**

- **CLI Layer (cmd/)**: Cobra-based commands matching current make targets
- **Database Layer (internal/db/)**: SQLite cache for fast queries; auto-syncs with markdown
- **Template Engine (internal/templates/)**: Go templates embedded via `//go:embed`
- **Parser (internal/parser/)**: Extract frontmatter, status, titles from markdown files
- **Generator (internal/generator/)**: Create markdown files from templates with substitutions
- **Config (internal/config/)**: Optional `.rex.yaml` for customization

**Data Flow:**

```
CREATE/UPDATE:
  User → CLI → Write .md file → Update SQLite → Done

LIST/QUERY:
  User → CLI → Query SQLite → Display results

REBUILD:
  User → CLI → Scan docs/ → Parse all .md → Rebuild SQLite → Done

AUTO-SYNC:
  CLI detects .rex.db missing/stale → Auto-rebuild from docs/ → Continue
```

### Key Architectural Decisions

- **SQLite as cache layer** using modernc.org/sqlite (pure Go, no CGO) with markdown as source of truth (see ADR-TBD-001)
- **Go embed package** for template bundling (see ADR-TBD-002)
- **Cobra CLI framework** for command structure and help system (see ADR-TBD-003)
- **YAML frontmatter parsing** using gopkg.in/yaml.v3 (see ADR-TBD-004)
- **Auto-rebuild strategy** for keeping SQLite cache in sync with markdown files (see ADR-TBD-005)
- **Semantic versioning** for rex binary releases (see ADR-TBD-006)

## Implementation Plan

### Phase 1: Foundation, Database, and Core Commands

**Scope:**
- Set up Go project structure with Cobra CLI
- Integrate SQLite (modernc.org/sqlite) with schema for ADRs, RFCs, Tasks, Plans
- Implement `rex init` to create docs/ directory structure
- Implement `rex rebuild` to scan docs/ and populate SQLite
- Implement auto-rebuild when .rex.db missing or stale
- Implement `rex adr create` replacing `create-adr.sh` (writes .md + updates DB)
- Implement `rex adr list` to query ADRs from DB
- Implement `rex adr update` replacing `update-adr-readme.sh`
- Embed ADR templates using Go embed
- Add tests for file creation, parsing, and DB operations

**Success Criteria:**
- `rex adr create "Title"` creates .md file and updates .rex.db
- `rex adr list` queries DB and displays results
- `rex adr list --status Draft` filters correctly
- `rex rebuild` scans docs/adr/ and populates DB
- Auto-rebuild triggers when .rex.db missing
- ADR template embedded in binary
- Zero external dependencies (no yq, sed, awk, no CGO required)

**Estimated Effort:** 3 weeks (increased for DB integration)

### Phase 2: RFC Support

**Scope:**
- Implement `rex rfc create` replacing `create-rfc.sh` (writes .md + updates DB)
- Implement `rex rfc list` to query RFCs from DB
- Implement `rex rfc update` replacing `update-rfc-readme.sh`
- Embed RFC templates
- Add RFC table to SQLite schema
- Add RFC parsing and README generation

**Success Criteria:**
- `rex rfc create "Title"` creates RFC with auto-incremented ID and updates DB
- `rex rfc list --status Draft` filters correctly
- `rex rfc update` regenerates docs/rfc/README.md
- `rex rebuild` includes RFCs in DB sync
- Feature parity with existing RFC bash scripts

**Estimated Effort:** 1 week

### Phase 3: Task Management

**Scope:**
- Implement `rex task create` with interactive prompts (writes .md + updates DB)
- Implement `rex task complete TASK-NNN` (updates .md + DB, moves file)
- Implement `rex task list` with filtering (TYPE, STATUS, PRIORITY) - queries DB
- Implement `rex task stats` for statistics dashboard - queries DB with aggregations
- Implement `rex task update` for README regeneration
- Embed task templates with YAML frontmatter
- Add task table to SQLite schema with all frontmatter fields

**Success Criteria:**
- Full feature parity with bash task management scripts
- Interactive prompts for metadata collection
- `rex task list --type core --status in_progress --priority P1` filters correctly
- `rex task stats` shows counts by type, status, priority (from DB)
- Task file movement (active/ to completed/) updates DB
- Fast queries even with 100+ tasks

**Estimated Effort:** 2 weeks

### Phase 4: Plan Support and Polish

**Scope:**
- Implement `rex plan create` and `rex plan update` (with DB support)
- Add `rex version` command
- Add `rex help` with examples
- Add `rex db info` to show DB stats (size, record counts, last rebuild)
- Configuration file support (.rex.yaml) for customization
- Error messages and validation
- Cross-platform binary releases (Linux, macOS, Windows)
- DB migration system for schema changes

**Success Criteria:**
- All four documentation types supported (ADR, RFC, Task, Plan)
- Helpful error messages and command documentation
- Binary releases for major platforms
- Migration guide for existing repos
- DB schema versioning works across rex versions

**Estimated Effort:** 1.5 weeks

### Phase 5: Migration and Documentation

**Scope:**
- Create migration guide for repos using bash scripts
- Update main README.md with installation instructions
- Add command reference documentation
- Create video/GIF demos of workflows
- Dogfood rex in this repository

**Success Criteria:**
- Clear migration path documented
- Installation via `go install` or binary download
- All bash scripts marked as deprecated
- Rex used for its own documentation

**Estimated Effort:** 1 week

## Database Schema

Rex uses SQLite (modernc.org/sqlite) as a query cache. The database is automatically generated and rebuilt from markdown files, so it can be safely deleted and regenerated at any time.

**Schema Overview:**

```sql
-- ADRs Table
CREATE TABLE adrs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    date TEXT NOT NULL,
    file_path TEXT NOT NULL UNIQUE,
    content TEXT,  -- Full markdown content for search
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_adrs_status ON adrs(status);
CREATE INDEX idx_adrs_number ON adrs(number);

-- RFCs Table
CREATE TABLE rfcs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    author TEXT,
    created_date TEXT NOT NULL,
    updated_date TEXT NOT NULL,
    file_path TEXT NOT NULL UNIQUE,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_rfcs_status ON rfcs(status);
CREATE INDEX idx_rfcs_number ON rfcs(number);

-- Tasks Table
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,  -- TASK-001 format
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    priority TEXT NOT NULL,
    estimated_hours REAL,
    actual_hours REAL,
    started_date TEXT,
    completed_date TEXT,
    assignee TEXT,
    phase INTEGER,
    file_path TEXT NOT NULL UNIQUE,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_type ON tasks(type);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_type_status ON tasks(type, status);

-- Task Relationships (for blocked_by, blocks, related_to)
CREATE TABLE task_relationships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id TEXT NOT NULL,
    related_task_id TEXT NOT NULL,
    relationship_type TEXT NOT NULL,  -- 'blocked_by', 'blocks', 'related_to'
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (related_task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_task_rel_task ON task_relationships(task_id);
CREATE INDEX idx_task_rel_type ON task_relationships(relationship_type);

-- Task Tags (many-to-many)
CREATE TABLE task_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id TEXT NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_task_tags_task ON task_tags(task_id);
CREATE INDEX idx_task_tags_tag ON task_tags(tag);

-- Plans Table
CREATE TABLE plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    created_date TEXT NOT NULL,
    file_path TEXT NOT NULL UNIQUE,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_plans_status ON plans(status);

-- Metadata Table (for DB versioning and rebuild tracking)
CREATE TABLE metadata (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Store DB schema version for migrations
INSERT INTO metadata (key, value) VALUES ('schema_version', '1');
INSERT INTO metadata (key, value) VALUES ('last_rebuild', datetime('now'));
```

**Query Examples:**

```sql
-- Fast task filtering
SELECT * FROM tasks
WHERE type = 'core'
  AND status = 'in_progress'
  AND priority IN ('P0', 'P1')
ORDER BY priority, created_at;

-- Task statistics
SELECT type, status, COUNT(*) as count, SUM(estimated_hours) as total_hours
FROM tasks
GROUP BY type, status;

-- Find blocked tasks with their blockers
SELECT t.id, t.title, GROUP_CONCAT(tr.related_task_id) as blocked_by
FROM tasks t
JOIN task_relationships tr ON t.id = tr.task_id
WHERE tr.relationship_type = 'blocked_by'
GROUP BY t.id;

-- Full-text search (if needed)
SELECT * FROM tasks WHERE content LIKE '%keyword%';
```

**Database Lifecycle:**

1. **Initial creation**: `rex init` or first command creates `.rex.db`
2. **Auto-rebuild**: Any command checks if DB is stale/missing and rebuilds
3. **Manual rebuild**: `rex rebuild` forces a full rebuild from docs/
4. **Stale detection**: Compare DB `last_rebuild` timestamp with file mtimes
5. **Migration**: On schema changes, rex auto-migrates or rebuilds

**Benefits:**

- **Fast queries**: No need to parse 100+ markdown files for `rex task list`
- **Complex filtering**: SQL WHERE clauses for multi-field filtering
- **Aggregations**: `rex task stats` runs SQL GROUP BY queries
- **Relationships**: Easy to query task dependencies via JOIN
- **Search**: Full-text search across all content if needed

## Migration Strategy

**For existing repositories using bash scripts:**

1. **Install rex binary**
   ```bash
   go install github.com/donaldgifford/rex@latest
   ```

2. **Build local database from existing markdown files**
   ```bash
   rex rebuild  # Scans docs/ and populates .rex.db
   ```

3. **Verify everything works**
   ```bash
   rex adr list  # Should show existing ADRs
   rex task list  # Should show existing tasks
   ```

4. **Remove old tooling**
   ```bash
   rm -rf tools/makefiles tools/docs tools/scripts
   git rm Makefile  # If only used for rex
   ```

5. **Add .rex.db to .gitignore**
   ```bash
   echo ".rex.db" >> .gitignore
   git add .gitignore
   ```

6. **Update workflows**
   - Replace `make adr "Title"` with `rex adr create "Title"`
   - Replace `make rfc "Title"` with `rex rfc create "Title"`
   - Replace `make task` with `rex task create`
   - Replace `make task-list TYPE=core` with `rex task list --type core`
   - Replace `make task-stats` with `rex task stats`

7. **Commit changes**
   ```bash
   git add -A
   git commit -m "Migrate to rex CLI tool"
   ```

**Backward Compatibility:**
- Rex will read existing ADR/RFC/Task files without modification
- Existing file formats remain unchanged (plain markdown)
- README generation produces same output format

**Deprecation Timeline:**
- Phase 1-4: Rex CLI developed alongside bash scripts
- Phase 5: Bash scripts marked deprecated but still functional
- 3 months after v1.0 release: Bash scripts removed from main branch

## Risks and Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Feature gaps compared to bash scripts | High | Medium | Maintain feature parity checklist; test each command against bash equivalent |
| Breaking changes to existing workflows | High | Low | Keep file formats identical; provide migration script to validate compatibility |
| Go binary distribution overhead | Medium | Low | Provide multiple installation methods (go install, homebrew, direct download); binaries are small (~10MB) |
| Template customization limited | Medium | Medium | Support .rex.yaml config for template overrides; embed default templates but allow local overrides |
| Cross-platform compatibility issues | Medium | Low | Test on Linux, macOS, Windows in CI; Go handles most platform differences |
| Users unfamiliar with CLI tools | Low | Medium | Provide comprehensive help text; create video demos; interactive prompts guide users |

## Success Criteria

**Technical Success:**
- Single rex binary < 15MB in size
- All templates embedded in binary
- Zero runtime dependencies (no yq, sed, awk)
- Feature parity with bash scripts (ADR, RFC, Task, Plan)
- Cross-platform support (Linux, macOS, Windows)
- Automated releases via GoReleaser

**Operational Success:**
- < 2 minutes to install and initialize new repo
- < 10 seconds to create any documentation type
- README updates in < 1 second
- 90% reduction in tooling LOC (2442 lines -> ~250 lines)

**Developer Experience:**
- Built-in help for all commands (`rex help`, `rex adr --help`)
- Clear error messages with suggestions
- Interactive prompts for complex commands (task create)
- No need to read source code to understand behavior
- Versioned releases with changelog

## Alternatives Considered

### Alternative 1: Keep Bash Scripts, Add Distribution

**Pros:**
- No rewrite required
- Current users already familiar with behavior
- Bash is universal on Unix systems

**Cons:**
- Still requires copying scripts to each repo
- No solution for Windows users
- Dependency management remains a problem
- Harder to version and distribute updates

**Decision:** Rejected—doesn't solve the core distribution and maintenance problems

### Alternative 2: Python CLI with pip Distribution

**Pros:**
- Rich ecosystem for CLI tools (click, typer)
- pip provides distribution mechanism
- Template engines well-established (Jinja2)

**Cons:**
- Adds Python as a dependency
- Virtual environments complicate installation
- Slower startup than Go binary
- Python version compatibility issues

**Decision:** Rejected—Go binary is more portable and has zero runtime dependencies

### Alternative 3: NPM Package with Node.js

**Pros:**
- npm provides good distribution
- JavaScript/TypeScript familiar to many developers
- Rich ecosystem for CLI tools

**Cons:**
- Requires Node.js runtime
- node_modules overhead
- Slower than Go binary
- Adds JavaScript build complexity

**Decision:** Rejected—Go binary is simpler and more performant

### Alternative 4: GitHub Actions Workflow

**Pros:**
- No local installation required
- Runs in CI/CD pipelines
- Centrally managed

**Cons:**
- Requires GitHub
- No offline usage
- Slower feedback loop (push, wait for CI)
- Not usable during local development

**Decision:** Rejected—developers need local tooling for interactive workflows

## Open Questions

1. Should rex support custom templates via local .rex/templates/ directory or only through .rex.yaml overrides?
2. Should rex generate HTML output (like the old ADR feature) or remain markdown-only?
3. Should rex support exporting documentation to other formats (PDF, Confluence, etc.)?
4. Should rex integrate with git hooks to auto-update READMEs and .rex.db on commit?
5. Should rex support task time tracking with timers (`rex task start`, `rex task stop`)?
6. Should rex support SQLite FTS5 (full-text search) for content search across all docs?
7. Should .rex.db location be configurable via .rex.yaml or always in repo root?

## References

- Current bash scripts: `tools/docs/`, `tools/scripts/tasks/`
- Current makefiles: `tools/makefiles/`
- Cobra CLI framework: https://github.com/spf13/cobra
- Go embed package: https://pkg.go.dev/embed
- SQLite (pure Go): https://gitlab.com/cznic/sqlite (modernc.org/sqlite)
- YAML parser: https://github.com/go-yaml/yaml
- GoReleaser: https://goreleaser.com/

## Appendix: Example Workflows

### Workflow 1: Creating a New ADR

**Current:**
```bash
cd /path/to/repo
make adr "Use PostgreSQL for Data Storage"
# Created: docs/adr/0015-use-postgresql-for-data-storage.md
make adr-update
git add docs/adr
git commit -m "Add ADR for PostgreSQL decision"
```

**Proposed:**
```bash
cd /path/to/repo
rex adr create "Use PostgreSQL for Data Storage"
# Created: docs/adr/0015-use-postgresql-for-data-storage.md
# Updated: docs/adr/README.md
# Updated: .rex.db (cache)
git add docs/adr
git commit -m "Add ADR for PostgreSQL decision"

# Query ADRs from DB
rex adr list --status Proposed
# 0015 | Use PostgreSQL for Data Storage | Proposed | 2025-11-05
```

**Improvement:** Single command auto-updates README and DB; no separate update step required. Fast queries from SQLite cache.

### Workflow 2: Initializing a New Repository

**Current:**
```bash
cd /path/to/new-repo
# Manually copy tools/ directory from another repo
cp -r /path/to/old-repo/tools .
# Manually copy Makefile
cp /path/to/old-repo/Makefile .
# Edit Makefile to remove non-rex targets
# Create docs directory structure
mkdir -p docs/{adr,rfc,tasks,plans}
# Copy templates
cp /path/to/old-repo/docs/adr/template.md docs/adr/
cp /path/to/old-repo/docs/rfc/template.md docs/rfc/
# ... etc for all templates
# Create README files
# Test that everything works
```

**Proposed:**
```bash
cd /path/to/new-repo
rex init
# Created: docs/adr/
# Created: docs/rfc/
# Created: docs/tasks/
# Created: docs/plans/
# Created: .rex.yaml (optional config)
# Created: .rex.db (SQLite cache)
# ✓ Repository initialized for documentation management

# Add .rex.db to .gitignore
echo ".rex.db" >> .gitignore
```

**Improvement:** Single command replaces 10+ minutes of manual setup. SQLite cache created automatically.

### Workflow 3: Task Management

**Current:**
```bash
make task
# Interactive prompts...
# Created: docs/tasks/core/active/TASK-042-implement-api.md
make task-update
# Updated READMEs
make task-list TYPE=core STATUS=in_progress
# Shows filtered tasks
make task-complete TASK=TASK-042
# Marks complete, moves file
```

**Proposed:**
```bash
rex task create
# Interactive prompts...
# Created: docs/tasks/core/active/TASK-042-implement-api.md
# Updated: docs/tasks/README.md
# Updated: .rex.db

# Fast queries from SQLite
rex task list --type core --status in_progress --priority P1
# TASK-042 | Implement API | core | in_progress | P1 | 8h

# Aggregated statistics
rex task stats
# Type: core (15 tasks, 120 hours estimated)
# Type: plugin (8 tasks, 64 hours estimated)
# Status: in_progress (5 tasks)
# Status: planned (18 tasks)

# Complete task
rex task complete TASK-042
# Updated: docs/tasks/core/completed/TASK-042-implement-api.md
# Updated: .rex.db
```

**Improvement:** READMEs auto-update; cleaner flag syntax; fast DB queries even with 100+ tasks; powerful statistics.

---

**Approval Signatures:**

_To be completed after review_
