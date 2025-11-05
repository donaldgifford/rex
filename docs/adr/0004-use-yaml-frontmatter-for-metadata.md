# 0004. Use YAML Frontmatter for Metadata

Date: 2025-11-05

## Status

Accepted

## Context

Rex manages structured documentation (ADRs, RFCs, Tasks, Plans) that needs both human-readable content and machine-parseable metadata. Tasks in particular have complex metadata:
- id, title, type, status, priority
- estimated_hours, actual_hours, dates
- blocked_by, blocks, related_to (relationships)
- assignee, tags, phase

We need a format that:
- Separates metadata from content clearly
- Is human-readable and editable
- Can be parsed reliably by Go code
- Works well with markdown
- Is git-friendly (text diffs)
- Is familiar to developers

The metadata must be parsed to populate the SQLite cache for fast queries.

## Decision

We will use YAML frontmatter (delimited by `---`) for all documentation metadata:

1. **Frontmatter Format**: YAML block at the top of markdown files, delimited by `---`
2. **Parser**: Use gopkg.in/yaml.v3 for YAML parsing
3. **Extraction**: Split file content on first `---` pair to extract frontmatter
4. **Validation**: Validate required fields after parsing
5. **Defaults**: Provide sensible defaults for optional fields

**Example (Task)**:
```markdown
---
id: TASK-001
title: Implement SQLite Cache
type: core
status: in_progress
priority: P1
estimated_hours: 8
actual_hours: null
started_date: 2025-11-05
completed_date: null
blocked_by: []
blocks: [TASK-002]
related_to: []
assignee: null
tags: [database, cache]
phase: 1
---

# Implement SQLite Cache

Task content here...
```

**Example (ADR)**:
```markdown
# 0001. Use SQLite as Cache Layer

Date: 2025-11-05

## Status

Accepted

...
```

Note: ADRs use inline metadata (Date, Status sections) rather than frontmatter for simplicity and tradition. Only Tasks use YAML frontmatter due to complexity.

## Consequences

### Positive

- **Structured**: Clear separation of metadata and content
- **Human-Readable**: YAML is easy to read and edit
- **Machine-Parseable**: Reliable parsing with gopkg.in/yaml.v3
- **Git-Friendly**: Text format, line-by-line diffs work well
- **Familiar**: YAML frontmatter used by Jekyll, Hugo, many static site generators
- **Extensible**: Easy to add new fields without breaking format
- **Type-Safe**: Go structs map directly to YAML fields
- **Validation**: Can validate required fields after unmarshaling
- **Nested Data**: Supports arrays and complex relationships
- **Comments**: YAML comments available if needed

### Negative

- **YAML Complexity**: YAML has edge cases (indentation, special characters)
- **Parsing Overhead**: Must parse YAML on every file read (~1ms per file)
- **User Errors**: Users can create invalid YAML (e.g., wrong indentation)
- **Format Strictness**: Must enforce frontmatter delimiter (`---`) exactly
- **Two Syntaxes**: Markdown + YAML = two languages in one file
- **Whitespace Sensitivity**: YAML is sensitive to indentation

### Neutral

- **Array Syntax**: YAML arrays can be `[a, b]` or `- a\n- b` (allow both)
- **Null Values**: Use `null` for empty optional fields (standard YAML)
- **Date Format**: Use ISO 8601 (YYYY-MM-DD) for consistency

## Alternatives Considered

### Alternative 1: TOML Frontmatter

**Pros**:
- Simpler syntax than YAML
- Less whitespace-sensitive
- Gaining popularity

**Cons**:
- Less familiar than YAML
- Hugo/Jekyll use YAML (consistency with other tools)
- Slightly less human-readable for nested structures

**Decision**: Rejected - YAML is more familiar to developers

### Alternative 2: JSON Frontmatter

**Pros**:
- Standard format
- Easy to parse
- Strict syntax

**Cons**:
- Not human-friendly (no comments, trailing commas forbidden)
- Verbose for simple data
- Harder to write by hand

**Decision**: Rejected - too verbose and unfriendly for humans

### Alternative 3: Inline Key-Value (Like ADRs)

**Example**:
```markdown
# Task: TASK-001

Type: core
Status: in_progress
Priority: P1
Estimated Hours: 8
```

**Pros**:
- Very readable
- Familiar (like ADR format)
- Simple to parse with regex

**Cons**:
- No support for arrays or relationships
- Harder to parse reliably
- Doesn't scale to complex metadata
- No standard format

**Decision**: Rejected - insufficient for complex task metadata

### Alternative 4: Separate Metadata File

**Example**: `TASK-001.md` + `TASK-001.meta.yaml`

**Pros**:
- Complete separation of concerns
- Easy to ignore content when parsing metadata
- Could use different formats

**Cons**:
- Two files per document (confusing)
- Must keep files in sync
- Harder to edit (two files to manage)
- More complex file operations

**Decision**: Rejected - single file is simpler

### Alternative 5: XML Frontmatter

**Pros**:
- Structured and validated
- Industry standard

**Cons**:
- Extremely verbose
- Not human-friendly
- Overkill for this use case

**Decision**: Rejected - XML is too heavy for this use case

## References

- gopkg.in/yaml.v3: https://github.com/go-yaml/yaml
- YAML Specification: https://yaml.org/spec/
- Jekyll Frontmatter: https://jekyllrb.com/docs/front-matter/
- Hugo Frontmatter: https://gohugo.io/content-management/front-matter/
- RFC 0001: Rex Documentation Management System
- ADR 0001: Use SQLite as Cache Layer
