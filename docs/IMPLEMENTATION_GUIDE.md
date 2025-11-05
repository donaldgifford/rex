# Rex Implementation Guide

This guide provides detailed implementation instructions for the Rex refactor from bash scripts to a Go CLI tool with SQLite caching.

## Prerequisites

**Required Reading:**
- RFC 0001: Rex Documentation Management System
- ADR 0001: Use SQLite as Cache Layer
- ADR 0002: Use Go Embed for Templates
- ADR 0003: Use Cobra CLI Framework
- ADR 0004: Use YAML Frontmatter for Metadata
- ADR 0005: Auto-Rebuild Cache on Stale Detection

**Development Environment:**
- Go 1.21+ (for embed support and modern features)
- SQLite knowledge (SQL queries, schema design)
- Git for version control
- Familiarity with Cobra CLI framework

## Project Structure

```
rex/
├── cmd/
│   ├── root.go              # Root command
│   ├── init.go              # Initialize repo
│   ├── rebuild.go           # Rebuild cache
│   ├── version.go           # Version command
│   ├── adr/
│   │   ├── adr.go          # ADR parent command
│   │   ├── create.go       # Create ADR
│   │   ├── list.go         # List ADRs
│   │   └── update.go       # Update README
│   ├── rfc/
│   │   ├── rfc.go          # RFC parent command
│   │   ├── create.go       # Create RFC
│   │   ├── list.go         # List RFCs
│   │   └── update.go       # Update README
│   └── task/
│       ├── task.go         # Task parent command
│       ├── create.go       # Create task
│       ├── complete.go     # Complete task
│       ├── list.go         # List tasks
│       ├── stats.go        # Task statistics
│       └── update.go       # Update README
├── internal/
│   ├── db/
│   │   ├── db.go           # SQLite connection management
│   │   ├── schema.go       # Schema definition and migrations
│   │   ├── adr.go          # ADR CRUD operations
│   │   ├── rfc.go          # RFC CRUD operations
│   │   ├── task.go         # Task CRUD operations
│   │   ├── plan.go         # Plan CRUD operations
│   │   └── rebuild.go      # Cache rebuild logic
│   ├── parser/
│   │   ├── frontmatter.go  # YAML frontmatter extraction
│   │   ├── markdown.go     # Markdown parsing utilities
│   │   └── adr.go          # ADR-specific parsing
│   ├── generator/
│   │   ├── file.go         # File creation from templates
│   │   ├── slug.go         # Title to filename slug
│   │   ├── id.go           # ID generation (TASK-001, ADR numbers)
│   │   └── readme.go       # README generation
│   ├── config/
│   │   └── config.go       # .rex.yaml parsing
│   └── templates/
│       ├── embed.go        # Embed directive
│       ├── adr.tmpl        # ADR template
│       ├── rfc.tmpl        # RFC template
│       ├── task.tmpl       # Task template
│       └── plan.tmpl       # Plan template
├── templates/               # Template files (embedded)
│   ├── adr/
│   │   ├── template.md
│   │   └── README.md
│   ├── rfc/
│   │   ├── template.md
│   │   └── README.md
│   └── tasks/
│       ├── template.md
│       └── README.md
├── main.go                  # Entry point
├── go.mod
├── go.sum
├── .goreleaser.yaml        # Release configuration
└── README.md
```

## Phase 1: Foundation & Database (3 weeks)

### Week 1: Project Setup & Database Layer

#### 1.1 Initialize Go Module and Dependencies

```bash
go mod init github.com/donaldgifford/rex
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get modernc.org/sqlite@latest
go get gopkg.in/yaml.v3@latest
```

**Files to create:**
- `main.go` - Entry point
- `cmd/root.go` - Root Cobra command
- `go.mod` - Module definition

#### 1.2 Database Schema Implementation

**File**: `internal/db/schema.go`

```go
package db

const Schema = `
-- ADRs Table
CREATE TABLE IF NOT EXISTS adrs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    date TEXT NOT NULL,
    file_path TEXT NOT NULL UNIQUE,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_adrs_status ON adrs(status);
CREATE INDEX IF NOT EXISTS idx_adrs_number ON adrs(number);

-- RFCs Table
CREATE TABLE IF NOT EXISTS rfcs (
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

CREATE INDEX IF NOT EXISTS idx_rfcs_status ON rfcs(status);
CREATE INDEX IF NOT EXISTS idx_rfcs_number ON rfcs(number);

-- Tasks Table
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
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

CREATE INDEX IF NOT EXISTS idx_tasks_type ON tasks(type);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX IF NOT EXISTS idx_tasks_type_status ON tasks(type, status);

-- Task Relationships
CREATE TABLE IF NOT EXISTS task_relationships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id TEXT NOT NULL,
    related_task_id TEXT NOT NULL,
    relationship_type TEXT NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (related_task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_task_rel_task ON task_relationships(task_id);
CREATE INDEX IF NOT EXISTS idx_task_rel_type ON task_relationships(relationship_type);

-- Task Tags
CREATE TABLE IF NOT EXISTS task_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id TEXT NOT NULL,
    tag TEXT NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_task_tags_task ON task_tags(task_id);
CREATE INDEX IF NOT EXISTS idx_task_tags_tag ON task_tags(tag);

-- Metadata
CREATE TABLE IF NOT EXISTS metadata (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`
```

#### 1.3 Database Connection Management

**File**: `internal/db/db.go`

```go
package db

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"

    _ "modernc.org/sqlite"
)

type Database struct {
    conn *sql.DB
    path string
}

func Open(path string) (*Database, error) {
    // Ensure directory exists
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("create db directory: %w", err)
    }

    // Open SQLite connection
    conn, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }

    // Set pragmas for performance
    pragmas := []string{
        "PRAGMA journal_mode=WAL",
        "PRAGMA synchronous=NORMAL",
        "PRAGMA foreign_keys=ON",
    }

    for _, pragma := range pragmas {
        if _, err := conn.Exec(pragma); err != nil {
            return nil, fmt.Errorf("set pragma: %w", err)
        }
    }

    db := &Database{
        conn: conn,
        path: path,
    }

    // Initialize schema
    if err := db.initSchema(); err != nil {
        return nil, err
    }

    return db, nil
}

func (db *Database) initSchema() error {
    _, err := db.conn.Exec(Schema)
    if err != nil {
        return fmt.Errorf("initialize schema: %w", err)
    }

    // Set schema version
    _, err = db.conn.Exec(
        "INSERT OR REPLACE INTO metadata (key, value) VALUES (?, ?)",
        "schema_version", "1",
    )
    return err
}

func (db *Database) Close() error {
    return db.conn.Close()
}

func (db *Database) GetLastRebuild() (string, error) {
    var timestamp string
    err := db.conn.QueryRow(
        "SELECT value FROM metadata WHERE key = ?",
        "last_rebuild",
    ).Scan(&timestamp)

    if err == sql.ErrNoRows {
        return "", nil
    }

    return timestamp, err
}

func (db *Database) SetLastRebuild() error {
    _, err := db.conn.Exec(
        "INSERT OR REPLACE INTO metadata (key, value, updated_at) VALUES (?, datetime('now'), datetime('now'))",
        "last_rebuild",
    )
    return err
}
```

**Testing checklist:**
- [ ] Database opens successfully
- [ ] Schema creates all tables
- [ ] Indexes are created
- [ ] Foreign keys are enforced
- [ ] Metadata table initialized with schema_version

### Week 2: Parser & Template System

#### 2.1 YAML Frontmatter Parser

**File**: `internal/parser/frontmatter.go`

```go
package parser

import (
    "bytes"
    "fmt"

    "gopkg.in/yaml.v3"
)

const delimiter = "---"

// ExtractFrontmatter extracts YAML frontmatter from markdown
func ExtractFrontmatter(content []byte, v interface{}) ([]byte, error) {
    // Find first delimiter
    parts := bytes.SplitN(content, []byte(delimiter), 3)

    if len(parts) < 3 {
        return nil, fmt.Errorf("no frontmatter found")
    }

    // parts[0] is empty or whitespace before first ---
    // parts[1] is the frontmatter
    // parts[2] is the markdown content

    frontmatter := bytes.TrimSpace(parts[1])
    markdown := bytes.TrimSpace(parts[2])

    if err := yaml.Unmarshal(frontmatter, v); err != nil {
        return nil, fmt.Errorf("parse frontmatter: %w", err)
    }

    return markdown, nil
}

// TaskFrontmatter represents task metadata
type TaskFrontmatter struct {
    ID             string   `yaml:"id"`
    Title          string   `yaml:"title"`
    Type           string   `yaml:"type"`
    Status         string   `yaml:"status"`
    Priority       string   `yaml:"priority"`
    EstimatedHours float64  `yaml:"estimated_hours"`
    ActualHours    *float64 `yaml:"actual_hours"`
    StartedDate    string   `yaml:"started_date"`
    CompletedDate  string   `yaml:"completed_date"`
    BlockedBy      []string `yaml:"blocked_by"`
    Blocks         []string `yaml:"blocks"`
    RelatedTo      []string `yaml:"related_to"`
    Assignee       string   `yaml:"assignee"`
    Tags           []string `yaml:"tags"`
    Phase          *int     `yaml:"phase"`
}

// Validate checks required fields
func (t *TaskFrontmatter) Validate() error {
    if t.ID == "" {
        return fmt.Errorf("id is required")
    }
    if t.Title == "" {
        return fmt.Errorf("title is required")
    }
    // ... more validation
    return nil
}
```

**Testing checklist:**
- [ ] Extracts frontmatter correctly
- [ ] Handles missing frontmatter
- [ ] Validates required fields
- [ ] Parses arrays correctly (blocked_by, tags)
- [ ] Handles null values

#### 2.2 Template Embedding

**File**: `internal/templates/embed.go`

```go
package templates

import (
    "embed"
    "io/fs"
)

//go:embed templates/**/*.md templates/**/*.tmpl
var FS embed.FS

// Get returns template content by path
func Get(path string) ([]byte, error) {
    return fs.ReadFile(FS, path)
}

// List returns all template paths
func List() ([]string, error) {
    var paths []string
    err := fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            paths = append(paths, path)
        }
        return nil
    })
    return paths, err
}
```

**Testing checklist:**
- [ ] Templates are embedded at compile time
- [ ] Can read template content
- [ ] List all templates
- [ ] Templates available in binary

### Week 3: ADR Commands & Auto-Rebuild

#### 3.1 ADR Create Command

**File**: `cmd/adr/create.go`

```go
package adr

import (
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/spf13/cobra"
    "github.com/donaldgifford/rex/internal/db"
    "github.com/donaldgifford/rex/internal/generator"
)

var createCmd = &cobra.Command{
    Use:   "create [title]",
    Short: "Create a new ADR",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        title := args[0]

        // Open database
        database, err := db.Open(".rex.db")
        if err != nil {
            return err
        }
        defer database.Close()

        // Generate ADR number
        number, err := database.GetNextADRNumber()
        if err != nil {
            return err
        }

        // Generate filename
        slug := generator.Slugify(title)
        filename := fmt.Sprintf("%04d-%s.md", number, slug)
        filepath := filepath.Join("docs", "adr", filename)

        // Generate from template
        data := map[string]interface{}{
            "Number": fmt.Sprintf("%04d", number),
            "Title":  title,
            "Date":   time.Now().Format("2006-01-02"),
        }

        if err := generator.CreateFromTemplate("adr", filepath, data); err != nil {
            return err
        }

        // Insert into database
        if err := database.InsertADR(number, title, "Proposed", filepath); err != nil {
            return err
        }

        fmt.Printf("✓ Created ADR: %s\n", filepath)
        fmt.Printf("✓ Updated: .rex.db\n")

        return nil
    },
}
```

**Testing checklist:**
- [ ] Creates ADR file with correct number
- [ ] Updates SQLite cache
- [ ] Generates proper slug from title
- [ ] Handles spaces and special characters
- [ ] Auto-increments ADR numbers

#### 3.2 Auto-Rebuild Implementation

**File**: `internal/db/rebuild.go`

```go
package db

import (
    "fmt"
    "os"
    "path/filepath"
    "time"

    "github.com/donaldgifford/rex/internal/parser"
)

func (db *Database) NeedsRebuild() (bool, error) {
    // Check if DB exists
    if _, err := os.Stat(db.path); os.IsNotExist(err) {
        return true, nil
    }

    // Get last rebuild timestamp
    lastRebuild, err := db.GetLastRebuild()
    if err != nil || lastRebuild == "" {
        return true, nil
    }

    rebuildTime, err := time.Parse(time.RFC3339, lastRebuild)
    if err != nil {
        return true, nil
    }

    // Check if any markdown file is newer
    hasNewer, err := hasNewerFiles("docs", rebuildTime)
    if err != nil {
        return false, err
    }

    return hasNewer, nil
}

func hasNewerFiles(dir string, since time.Time) (bool, error) {
    var newer bool

    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && filepath.Ext(path) == ".md" {
            if info.ModTime().After(since) {
                newer = true
                return filepath.SkipDir
            }
        }

        return nil
    })

    return newer, err
}

func (db *Database) Rebuild() error {
    fmt.Println("⚡ Rebuilding cache...")

    // Clear existing data
    if err := db.clearTables(); err != nil {
        return err
    }

    // Scan and parse all docs
    if err := db.scanADRs(); err != nil {
        return fmt.Errorf("scan ADRs: %w", err)
    }

    if err := db.scanRFCs(); err != nil {
        return fmt.Errorf("scan RFCs: %w", err)
    }

    if err := db.scanTasks(); err != nil {
        return fmt.Errorf("scan tasks: %w", err)
    }

    // Update last rebuild timestamp
    if err := db.SetLastRebuild(); err != nil {
        return err
    }

    fmt.Println("✓ Cache rebuilt successfully")
    return nil
}
```

**Testing checklist:**
- [ ] Detects missing .rex.db
- [ ] Detects stale cache via mtime
- [ ] Rebuilds cache from markdown files
- [ ] Updates last_rebuild timestamp
- [ ] Shows progress indicator

## Phase 2: RFC Support (1 week)

### RFC Commands Implementation

Similar to ADR commands but with RFC-specific schema:
- RFC numbering (0001, 0002, etc.)
- RFC status (Draft, Proposed, Accepted, Rejected, Superseded)
- Author field
- Created/Updated dates

**Files to create:**
- `cmd/rfc/create.go`
- `cmd/rfc/list.go`
- `cmd/rfc/update.go`
- `internal/db/rfc.go`

**Testing checklist:**
- [ ] RFC creation with auto-incrementing IDs
- [ ] RFC listing with filters
- [ ] README generation
- [ ] Database cache updates

## Phase 3: Task Management (2 weeks)

### Task Commands Implementation

Complex due to YAML frontmatter and relationships:

**Files to create:**
- `cmd/task/create.go` - Interactive prompts
- `cmd/task/complete.go` - Mark complete, move file
- `cmd/task/list.go` - Filter by type/status/priority
- `cmd/task/stats.go` - Aggregations and statistics
- `internal/db/task.go` - Task CRUD with relationships

**Key features:**
- Interactive prompts for all metadata fields
- Relationship management (blocked_by, blocks, related_to)
- Tag support (many-to-many)
- File movement (active/ to completed/)
- Statistics with SQL aggregations

**Testing checklist:**
- [ ] Task creation with all frontmatter fields
- [ ] Filtering by multiple criteria
- [ ] Statistics calculations
- [ ] File movement on completion
- [ ] Relationship queries work correctly

## Testing Strategy

### Unit Tests

```go
// internal/parser/frontmatter_test.go
func TestExtractFrontmatter(t *testing.T) {
    content := []byte(`---
id: TASK-001
title: Test Task
type: core
---

# Test Task Content`)

    var fm parser.TaskFrontmatter
    markdown, err := parser.ExtractFrontmatter(content, &fm)

    assert.NoError(t, err)
    assert.Equal(t, "TASK-001", fm.ID)
    assert.Contains(t, string(markdown), "Test Task Content")
}
```

### Integration Tests

```go
// Test full workflow
func TestCreateADRWorkflow(t *testing.T) {
    // Setup temp directory
    tmpDir := t.TempDir()
    os.Chdir(tmpDir)

    // Initialize
    os.MkdirAll("docs/adr", 0755)

    // Create ADR
    cmd := exec.Command("rex", "adr", "create", "Test ADR")
    err := cmd.Run()
    assert.NoError(t, err)

    // Verify file exists
    _, err = os.Stat("docs/adr/0001-test-adr.md")
    assert.NoError(t, err)

    // Verify in database
    db, _ := db.Open(".rex.db")
    count, _ := db.CountADRs()
    assert.Equal(t, 1, count)
}
```

## Performance Benchmarks

Target performance metrics:

- **ADR Creation**: < 100ms
- **Task List (100 tasks)**: < 100ms
- **Task Stats (100 tasks)**: < 200ms
- **Cache Rebuild (100 files)**: < 2 seconds
- **Binary Size**: < 15MB

Benchmark code:

```go
func BenchmarkTaskList(b *testing.B) {
    db := setupTestDB(100) // 100 tasks

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        tasks, _ := db.ListTasks(db.TaskFilters{})
    }
}
```

## Release Process

1. **Update version** in `cmd/version.go`
2. **Update CHANGELOG.md**
3. **Tag release**: `git tag v0.1.0`
4. **Push**: `git push origin v0.1.0`
5. **GoReleaser** builds binaries automatically
6. **Test binaries** on Linux, macOS, Windows
7. **Publish GitHub release** with binaries

## Migration from Bash Scripts

Create migration script to help users:

```bash
#!/bin/bash
# migrate-to-rex.sh

echo "Migrating to rex CLI..."

# Install rex
go install github.com/donaldgifford/rex@latest

# Rebuild cache from existing markdown
rex rebuild

# Verify
rex adr list
rex task list

# Remove old tooling
echo "Remove old tools? (y/n)"
read answer
if [ "$answer" = "y" ]; then
    rm -rf tools/makefiles tools/docs tools/scripts
    echo "✓ Old tooling removed"
fi

echo "✓ Migration complete!"
echo "Update your workflows:"
echo "  make adr \"Title\" → rex adr create \"Title\""
echo "  make task → rex task create"
```

## Troubleshooting

### Database Locked

```bash
# Another process has the database open
# Solution: Wait or kill other rex processes
pkill rex
```

### Cache Out of Sync

```bash
# Force rebuild
rex rebuild
```

### Template Not Found

```bash
# Templates should be embedded
# Rebuild rex binary
go build -o rex main.go
```

## References

- RFC 0001: Rex Documentation Management System
- All ADRs in docs/adr/
- Cobra Documentation: https://cobra.dev/
- modernc.org/sqlite: https://pkg.go.dev/modernc.org/sqlite
- Go embed: https://pkg.go.dev/embed
