# 0001. Use SQLite as Cache Layer

Date: 2025-11-05

## Status

Accepted

## Context

Rex needs to efficiently query and filter documentation files (ADRs, RFCs, Tasks, Plans) without parsing all markdown files on every command execution. Users may have 100+ tasks or dozens of ADRs, making repeated file parsing slow and inefficient.

Key requirements:
- Fast queries for filtering (e.g., `rex task list --type core --status in_progress --priority P1`)
- Complex aggregations (e.g., `rex task stats` showing counts by type/status)
- Task relationship queries (blocked_by, blocks, related_to)
- Maintain markdown files as source of truth for git-friendly diffs
- No external dependencies or CGO for cross-platform compatibility

## Decision

We will use SQLite (via modernc.org/sqlite) as a local query cache, with markdown files remaining the source of truth:

1. **Markdown as Source of Truth**: All documentation files remain plain markdown with YAML frontmatter, stored in `docs/` directories
2. **SQLite as Cache**: A local `.rex.db` file caches parsed documentation for fast queries
3. **Auto-Rebuild**: The CLI automatically detects when the cache is stale/missing and rebuilds from markdown
4. **Manual Rebuild**: Users can run `rex rebuild` to force a full cache rebuild
5. **Gitignored DB**: The `.rex.db` file is gitignored - only markdown files are version controlled
6. **Pure Go SQLite**: Use modernc.org/sqlite (pure Go implementation, no CGO)

**Data Flow**:
- **Create/Update**: CLI writes markdown file → Updates SQLite cache → Done
- **Query/List**: CLI queries SQLite cache → Returns results instantly
- **Rebuild**: CLI scans `docs/` → Parses markdown → Populates SQLite

## Consequences

### Positive

- **Fast Queries**: No need to parse 100+ markdown files for every list/filter command
- **Complex Filtering**: SQL WHERE clauses enable multi-field filtering efficiently
- **Aggregations**: SQL GROUP BY for statistics (task counts, time estimates)
- **Relationships**: SQL JOINs for task dependencies and relationships
- **Scalability**: Handles hundreds of documents without performance degradation
- **Git-Friendly**: Markdown files remain the source, so git diffs work normally
- **No CGO**: modernc.org/sqlite is pure Go, enabling easy cross-compilation
- **Resilient**: Cache can be deleted/corrupted - just rebuild from markdown
- **Full-Text Search**: Can add SQLite FTS5 in future for content search

### Negative

- **Cache Staleness**: DB can become out of sync if markdown files edited outside rex
- **Disk Space**: Adds ~1-5MB `.rex.db` file per repository (small but non-zero)
- **Complexity**: Two representations of data (markdown + SQLite) to maintain
- **Learning Curve**: Developers need to understand cache rebuild mechanism
- **Migration Burden**: Schema changes require migration logic or rebuild

### Neutral

- **Auto-Rebuild Overhead**: First command after manual markdown edits may be slower (1-2 seconds to rebuild)
- **Query vs. Parse Tradeoff**: Gains in query speed, small overhead in create/update
- **Stale Detection**: Need to compare file mtimes vs. DB last_rebuild timestamp

## Alternatives Considered

### Alternative 1: Parse Markdown on Every Command

**Pros**:
- No cache to maintain
- Always in sync with files
- Simpler architecture

**Cons**:
- Slow for 100+ files (5-10 seconds)
- Cannot do complex SQL queries
- Poor user experience

**Decision**: Rejected - performance is too important for CLI tools

### Alternative 2: In-Memory Cache Only

**Pros**:
- Fast queries
- No disk storage
- No staleness issues

**Cons**:
- Must rebuild cache on every rex invocation
- Slow startup time (2-3 seconds)
- Cannot persist query results

**Decision**: Rejected - startup time matters for CLI UX

### Alternative 3: Index Files (One File per Type)

**Pros**:
- Git-tracked index
- No binary format
- Simple to implement

**Cons**:
- Requires manual updates
- No query flexibility
- Duplication of data
- Merge conflicts on index files

**Decision**: Rejected - maintenance burden too high

### Alternative 4: BoltDB/BadgerDB Key-Value Store

**Pros**:
- Pure Go
- Simple key-value model
- Embedded database

**Cons**:
- No SQL queries (need custom filtering)
- Manual indexing required
- Less flexible than SQL
- Requires more code for complex queries

**Decision**: Rejected - SQL provides better query flexibility

### Alternative 5: JSON Cache File

**Pros**:
- Human-readable
- Simple to implement
- Git-diffable if needed

**Cons**:
- No query optimization
- Must load entire file into memory
- No indexing support
- Slow for large datasets

**Decision**: Rejected - doesn't scale beyond 50-100 items

## References

- RFC 0001: Rex Documentation Management System
- modernc.org/sqlite: https://gitlab.com/cznic/sqlite
- SQLite Documentation: https://www.sqlite.org/docs.html
- Database schema defined in RFC 0001, Section "Database Schema"
