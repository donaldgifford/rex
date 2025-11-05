# 0002. Use Go Embed for Templates

Date: 2025-11-05

## Status

Accepted

## Context

Rex needs to distribute templates for ADRs, RFCs, Tasks, and Plans to users without requiring:
- Copying template files to each repository
- Network access to fetch templates
- External template repositories
- Per-repository configuration for template locations

The goal is to make Rex a single, self-contained binary that can be installed once and used across all repositories with zero setup beyond `rex init`.

Template requirements:
- Markdown templates for each documentation type (ADR, RFC, Task, Plan)
- README templates for each documentation type
- Config file templates (.rex.yaml)
- Must support template variable substitution (title, date, number, etc.)
- Must be versioned with the rex binary
- Must work offline

## Decision

We will use Go's `embed` package (introduced in Go 1.16) to bundle all templates directly into the rex binary:

1. **Embed Directory**: Create a `templates/` directory in the rex repository containing all template files
2. **Go Embed Directive**: Use `//go:embed` directive to embed the entire templates directory
3. **Template Access**: Use `embed.FS` to access templates at runtime
4. **Go Templates**: Use Go's `text/template` package for variable substitution
5. **No External Files**: Templates are compiled into the binary - no external files needed

**Implementation**:
```go
package templates

import "embed"

//go:embed adr/*.md rfc/*.md tasks/*.md plans/*.md
var Templates embed.FS
```

**Usage**:
```go
// Load template from embedded FS
tmpl, err := template.ParseFS(Templates, "adr/template.md")

// Execute with data
err = tmpl.Execute(writer, data)
```

## Consequences

### Positive

- **Self-Contained**: Single binary contains all templates - no external files needed
- **Versioned**: Templates are versioned with the binary - upgrade rex, get new templates
- **Offline**: Works without network access
- **Fast**: No file I/O at runtime - templates loaded from memory
- **Consistent**: All users with same rex version get identical templates
- **Simple Distribution**: No need to package templates separately
- **Cross-Platform**: Works identically on Linux, macOS, Windows
- **No Configuration**: Users don't need to configure template paths

### Negative

- **Binary Size**: Adds ~10-20KB to binary size per template set (negligible)
- **Template Customization**: Requires building from source or using .rex.yaml overrides
- **Update Process**: Template changes require new rex binary release
- **Memory Usage**: All templates loaded into memory (small footprint ~100KB total)

### Neutral

- **Go 1.16+ Required**: Requires Go 1.16 or later to build (reasonable requirement in 2025)
- **Template Language**: Locked into Go templates (not Jinja2, Handlebars, etc.)
- **Compilation Time**: Slightly longer compile time to embed files (negligible)

## Alternatives Considered

### Alternative 1: External Template Files

**Pros**:
- Easy to customize templates
- No recompilation needed for changes
- Users can edit templates directly

**Cons**:
- Must copy templates to each repository
- Templates can diverge across repos
- Requires template distribution mechanism
- Version mismatch between binary and templates
- Offline usage requires pre-fetching

**Decision**: Rejected - defeats the goal of single binary distribution

### Alternative 2: Fetch Templates from GitHub

**Pros**:
- Always latest templates
- Easy to update centrally
- No binary size increase

**Cons**:
- Requires network access
- Slow (network latency)
- Fails offline
- Versioning issues (which version to fetch?)
- Rate limiting concerns
- Security concerns (MITM attacks)

**Decision**: Rejected - offline usage is critical for CLI tools

### Alternative 3: Generate Templates at Runtime

**Pros**:
- No file I/O
- Customizable
- Small binary

**Cons**:
- Hard to maintain (templates in Go strings)
- Difficult to preview/edit templates
- Poor developer experience
- No syntax highlighting for markdown

**Decision**: Rejected - maintainability is too important

### Alternative 4: Use go-bindata or pkger

**Pros**:
- Similar to embed
- Works with older Go versions

**Cons**:
- External dependency
- More complex build process
- Deprecated (go-bindata)
- go:embed is standard library

**Decision**: Rejected - use standard library when possible

### Alternative 5: Template Registry Service

**Pros**:
- Centralized template management
- Easy updates
- Shared across teams

**Cons**:
- Requires infrastructure
- Network dependency
- Complexity overhead
- Latency issues

**Decision**: Rejected - too much complexity for this use case

## References

- Go embed package: https://pkg.go.dev/embed
- Go text/template: https://pkg.go.dev/text/template
- RFC 0001: Rex Documentation Management System
- Go 1.16 Release Notes: https://golang.org/doc/go1.16#embed
