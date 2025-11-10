# Contributing to Rex

Thank you for your interest in contributing to Rex! This guide will help you get started.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Commit Messages](#commit-messages)
- [Documentation](#documentation)
- [Getting Help](#getting-help)

---

## Getting Started

### Prerequisites

- **Go 1.23.2 or later** - [Install Go](https://golang.org/doc/install)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **Make** (optional) - For convenience commands

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/rex.git
cd rex
```

3. Add upstream remote:

```bash
git remote add upstream https://github.com/donaldgifford/rex.git
```

4. Create a feature branch:

```bash
git checkout -b feature/your-feature-name
```

---

## Development Setup

### Install Dependencies

```bash
# Download Go dependencies
go mod download

# Verify everything works
go build -o rex main.go
./rex version
```

### Build and Run

```bash
# Build the binary
go build -o rex main.go

# Run directly
./rex --help

# Install locally for testing
go install

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Development Workflow

```bash
# Make changes to code
$EDITOR cmd/task_list.go

# Build and test locally
go build -o rex main.go
./rex task list

# Run tests
go test ./internal/db/...

# Check for issues
go vet ./...

# Format code
go fmt ./...
```

---

## Project Structure

```
rex/
├── main.go                    # Entry point
├── cmd/                       # CLI commands (Cobra)
│   ├── root.go               # Root command
│   ├── adr.go                # ADR commands
│   ├── adr_create.go
│   ├── adr_list.go
│   ├── adr_update.go
│   ├── rfc.go                # RFC commands
│   ├── task.go               # Task commands
│   ├── plan.go               # Plan commands
│   ├── version.go            # Version command
│   └── db.go                 # Database commands
├── internal/                  # Internal packages
│   ├── db/                   # Database layer
│   │   ├── db.go            # Database connection
│   │   ├── schema.go        # Schema definition
│   │   ├── rebuild.go       # Cache rebuild logic
│   │   ├── adr.go           # ADR operations
│   │   ├── rfc.go           # RFC operations
│   │   ├── task.go          # Task operations
│   │   └── plan.go          # Plan operations
│   └── parser/               # Markdown/YAML parsing
│       └── parser.go
├── templates/                 # Document templates (embedded)
│   ├── adr.md
│   ├── rfc.md
│   ├── task.md
│   └── plan.md
├── docs/                      # Documentation
│   ├── adr/                  # ADRs
│   ├── rfc/                  # RFCs
│   ├── tasks/                # Tasks
│   ├── plans/                # Plans
│   ├── MIGRATION.md          # Migration guide
│   ├── COMMANDS.md           # Command reference
│   ├── IMPLEMENTATION_GUIDE.md
│   └── MVP_GUIDE.md
├── .goreleaser.yaml           # Release configuration
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── README.md                  # Project README
└── CONTRIBUTING.md            # This file
```

### Key Packages

**cmd/** - CLI command implementations
- Uses [Cobra](https://github.com/spf13/cobra) for CLI framework
- Each command in its own file
- Follow existing patterns for consistency

**internal/db/** - Database operations
- Uses [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) (pure Go, no CGO)
- CRUD operations for each document type
- Transaction management
- Cache rebuild logic

**internal/parser/** - Parsing logic
- YAML frontmatter parsing ([gopkg.in/yaml.v3](https://gopkg.in/yaml.v3))
- Markdown content extraction
- File I/O utilities

**templates/** - Document templates
- Embedded using `//go:embed`
- Markdown with YAML frontmatter
- One template per document type

---

## Making Changes

### Types of Contributions

We welcome:

- **Bug fixes** - Fix issues in existing functionality
- **Features** - Add new functionality (check existing issues first)
- **Documentation** - Improve or add documentation
- **Tests** - Add or improve test coverage
- **Performance** - Optimize existing code
- **Refactoring** - Improve code quality

### Before You Start

1. **Check existing issues** - Someone might already be working on it
2. **Create an issue** - Discuss major changes before implementing
3. **Follow conventions** - Match existing code style and patterns
4. **Write tests** - Add tests for new functionality
5. **Update docs** - Update documentation for user-facing changes

### Development Guidelines

#### 1. Keep Changes Focused

- One feature/fix per pull request
- Small, reviewable changes
- Avoid unrelated refactoring

#### 2. Write Tests

```go
// Example test structure
func TestCreateTask(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()

    // Execute
    task := &Task{
        ID:    "TASK-001",
        Title: "Test Task",
        Type:  "core",
    }
    err := db.CreateTask(task)

    // Assert
    if err != nil {
        t.Fatalf("CreateTask failed: %v", err)
    }

    // Verify
    retrieved, err := db.GetTask("TASK-001")
    if err != nil {
        t.Fatalf("GetTask failed: %v", err)
    }
    if retrieved.Title != "Test Task" {
        t.Errorf("Expected title 'Test Task', got '%s'", retrieved.Title)
    }
}
```

#### 3. Handle Errors Properly

```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create task: %w", err)
}

// Good: Return errors to caller
func CreateTask(task *Task) error {
    if err := validate(task); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // ...
}

// Bad: Ignore errors
db.Exec(query) // Missing error check
```

#### 4. Follow Go Conventions

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use meaningful variable names
- Keep functions small and focused
- Document exported functions

```go
// Good: Well-documented exported function
// CreateTask creates a new task in the database.
// It validates the task fields and returns an error if validation fails.
func (db *Database) CreateTask(task *Task) error {
    // Implementation
}

// Good: Clear variable names
taskID := "TASK-001"
estimatedHours := 8.0

// Bad: Unclear variable names
id := "TASK-001"
h := 8.0
```

#### 5. Maintain Backward Compatibility

- Don't break existing APIs
- Deprecate before removing
- Add migration paths for breaking changes

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/db/

# Run with coverage
go test -cover ./...

# Run with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...

# Run specific test
go test -run TestCreateTask ./internal/db/

# Verbose output
go test -v ./...
```

### Writing Tests

**Unit Tests:**
```go
// internal/db/task_test.go
func TestCreateTask(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    task := &Task{
        ID:    "TASK-001",
        Title: "Test",
        Type:  "core",
    }

    err := db.CreateTask(task)
    if err != nil {
        t.Fatalf("CreateTask failed: %v", err)
    }
}
```

**Integration Tests:**
```go
// Test full workflow
func TestTaskWorkflow(t *testing.T) {
    // Create task
    // List tasks
    // Complete task
    // Verify completion
}
```

**Table-Driven Tests:**
```go
func TestValidatePriority(t *testing.T) {
    tests := []struct {
        name     string
        priority string
        wantErr  bool
    }{
        {"valid P0", "P0", false},
        {"valid P1", "P1", false},
        {"invalid", "P5", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validatePriority(tt.priority)
            if (err != nil) != tt.wantErr {
                t.Errorf("validatePriority() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Test Coverage Goals

- **Minimum**: 70% coverage
- **Target**: 80%+ coverage
- **Critical paths**: 90%+ coverage (database operations, parsing)

---

## Submitting Changes

### Pull Request Process

1. **Update your branch**:

```bash
git fetch upstream
git rebase upstream/main
```

2. **Run tests**:

```bash
go test ./...
go vet ./...
go fmt ./...
```

3. **Commit your changes**:

```bash
git add .
git commit -m "feat: add task filtering by assignee"
```

4. **Push to your fork**:

```bash
git push origin feature/your-feature-name
```

5. **Create pull request** on GitHub

### Pull Request Checklist

Before submitting, ensure:

- [ ] Code builds successfully (`go build`)
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No vet warnings (`go vet ./...`)
- [ ] Added tests for new functionality
- [ ] Updated documentation if needed
- [ ] Commit messages follow conventions
- [ ] PR description explains the change
- [ ] Linked to related issue (if applicable)

### PR Description Template

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Related Issue
Fixes #123

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] Code formatted
```

---

## Code Style

### Go Style

Follow [Uber's Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md):

**Formatting:**
```go
// Use gofmt (automatic with go fmt)
if err != nil {
    return err
}

// Not
if err!=nil{return err}
```

**Naming:**
```go
// Exported: PascalCase
func CreateTask() {}
type TaskDatabase struct {}

// Unexported: camelCase
func validateTask() {}
var taskCache map[string]Task
```

**Error Messages:**
```go
// Good: lowercase, no punctuation
return fmt.Errorf("failed to create task")

// Bad: uppercase, with punctuation
return fmt.Errorf("Failed to create task.")
```

**Imports:**
```go
// Group: stdlib, external, internal
import (
    "fmt"
    "os"

    "github.com/spf13/cobra"

    "github.com/donaldgifford/rex/internal/db"
)
```

### Comments

```go
// Good: Complete sentence, explains why
// CreateTask inserts a new task into the database.
// It validates the task before insertion to ensure data integrity.
func CreateTask(task *Task) error {
    // Validate required fields first
    if task.ID == "" {
        return fmt.Errorf("task ID is required")
    }
    // ...
}

// Bad: Incomplete, explains what (already obvious)
// create task
func CreateTask(task *Task) error {
```

### SQL Style

```go
// Good: Multi-line, readable
const query = `
    SELECT id, title, type, status
    FROM tasks
    WHERE type = ?
    ORDER BY priority DESC, id ASC
`

// Use placeholders for safety
rows, err := db.Query(query, taskType)
```

---

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

### Examples

```bash
# Good commit messages
feat(task): add filtering by assignee
fix(db): prevent database lock on rebuild
docs(readme): update installation instructions
test(task): add tests for task completion
refactor(parser): simplify YAML parsing logic

# With body
feat(task): add time tracking for tasks

Add actual_hours field to track time spent on tasks.
This allows users to compare estimated vs actual hours
and improve future estimates.

Closes #42
```

### Commit Message Guidelines

- **First line**: 50 characters or less
- **Body**: Wrap at 72 characters
- **Use imperative mood**: "add feature" not "added feature"
- **Reference issues**: "Closes #123" or "Fixes #456"
- **Explain why**, not what (code shows what)

---

## Documentation

### Types of Documentation

1. **Code Comments** - Explain complex logic
2. **Package Documentation** - Document packages and exported functions
3. **User Documentation** - README, guides, command reference
4. **ADRs** - Document architectural decisions

### Updating Documentation

When adding features:

- Update [README.md](README.md) if user-facing
- Update [COMMANDS.md](docs/COMMANDS.md) for new commands
- Add examples to demonstrate usage
- Update help text in Cobra commands

### Writing Documentation

- **Be clear and concise**
- **Provide examples**
- **Explain why**, not just how
- **Keep it up-to-date**

---

## Getting Help

### Questions?

- **Open an issue** - For bug reports or feature requests
- **Discussions** - For general questions or ideas
- **Documentation** - Check existing docs first

### Useful Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Effective Go](https://golang.org/doc/effective_go)

### Community

- **Be respectful** - Follow the [Code of Conduct](CODE_OF_CONDUCT.md)
- **Be patient** - Maintainers are volunteers
- **Be helpful** - Help others when you can

---

## Release Process

Releases are automated using [GoReleaser](https://goreleaser.com/).

### Creating a Release

1. **Update version** in code (if needed)
2. **Tag the release**:

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

3. **GoReleaser runs automatically** via GitHub Actions
4. **Binaries are published** to GitHub Releases

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0) - Breaking changes
- **MINOR** (v1.1.0) - New features, backward compatible
- **PATCH** (v1.0.1) - Bug fixes

---

## Code of Conduct

Be respectful, inclusive, and constructive:

- **Be welcoming** - Welcome newcomers
- **Be respectful** - Respect differing opinions
- **Be professional** - Keep discussions on topic
- **Be constructive** - Provide actionable feedback

Unacceptable behavior:
- Harassment or discrimination
- Trolling or insulting comments
- Personal attacks
- Publishing private information

---

## License

By contributing to Rex, you agree that your contributions will be licensed under the MIT License.

---

## Thank You!

Thank you for contributing to Rex! Every contribution, no matter how small, helps make Rex better for everyone.

**Happy coding! 🦖**
