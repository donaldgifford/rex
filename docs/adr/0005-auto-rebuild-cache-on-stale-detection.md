# 0005. Auto-Rebuild Cache on Stale Detection

Date: 2025-11-05

## Status

Accepted

## Context

Rex uses a SQLite cache (`.rex.db`) for fast queries, with markdown files as the source of truth. The cache can become stale when:
- User manually edits markdown files outside rex
- User pulls changes from git
- User switches branches
- User runs `git rebase`, `git merge`, etc.
- `.rex.db` is deleted or corrupted

Without proper handling, stale cache leads to:
- Queries returning outdated/incorrect data
- Missing documents not showing in `rex task list`
- User confusion when CLI output doesn't match files

We need a strategy that:
- Detects when cache is stale
- Automatically rebuilds when needed
- Provides manual rebuild command
- Minimizes rebuild overhead
- Gives user visibility into rebuild process

## Decision

We will implement automatic cache rebuild with intelligent stale detection:

1. **Stale Detection**: Check cache staleness on every rex command that queries data
2. **Detection Logic**:
   - If `.rex.db` doesn't exist → rebuild
   - If any markdown file mtime > last_rebuild timestamp → rebuild
   - If docs/ directory structure changed → rebuild
3. **Auto-Rebuild**: Automatically rebuild cache when stale (with progress indicator)
4. **Manual Rebuild**: Provide `rex rebuild` command for forced rebuild
5. **Rebuild Optimization**: Only scan changed files during auto-rebuild when possible
6. **User Feedback**: Show rebuild progress for operations taking >1 second

**Implementation**:

```go
func ensureCacheValid(db *sql.DB) error {
    // Check if DB exists
    if !dbExists() {
        return rebuildCache(db)
    }

    // Get last rebuild timestamp
    lastRebuild := getLastRebuild(db)

    // Check if any files are newer
    if hasNewerFiles("docs/", lastRebuild) {
        fmt.Println("⚡ Cache stale, rebuilding...")
        return rebuildCache(db)
    }

    return nil
}

// Called before any query operation
func (db *Database) List(filters Filters) ([]Task, error) {
    if err := ensureCacheValid(db.conn); err != nil {
        return nil, err
    }
    // Continue with query...
}
```

**Rebuild Triggers**:
- `.rex.db` missing
- File mtime check (any .md file newer than last_rebuild)
- Manual `rex rebuild` command
- Cache corruption detected

## Consequences

### Positive

- **Always Fresh**: Queries always return current data
- **No Manual Intervention**: Users don't need to remember to rebuild
- **Git-Friendly**: Works seamlessly with git operations
- **Fast Common Case**: No rebuild needed when cache is valid (99% of cases)
- **Transparent**: Rebuild happens automatically, user sees progress
- **Safe**: Corruption automatically fixed by rebuild
- **Developer-Friendly**: Can edit markdown files freely
- **Simple Mental Model**: "Cache always reflects current files"

### Negative

- **First Command Latency**: First command after git pull is slower (1-3 seconds)
- **Filesystem Overhead**: Must stat() all markdown files to check mtimes
- **False Positives**: Touch a file without changes → triggers rebuild
- **Concurrent Access**: Multiple rex processes could rebuild simultaneously
- **Progress Interruption**: User sees rebuild message during workflow

### Neutral

- **Rebuild Frequency**: Depends on user workflow (mostly after git operations)
- **Performance vs. Freshness**: Trade slower first command for guaranteed freshness
- **Cache File**: .rex.db should be gitignored to avoid conflicts

## Alternatives Considered

### Alternative 1: Manual Rebuild Only

**Pros**:
- No automatic overhead
- User has full control
- Predictable performance

**Cons**:
- Users forget to rebuild
- Stale data causes confusion
- Poor user experience
- Requires documentation/training

**Decision**: Rejected - manual rebuild is error-prone

### Alternative 2: Watch Files for Changes

**Pros**:
- Real-time updates
- No stale detection needed
- Could use fsnotify

**Cons**:
- Background daemon required
- More complex
- Resource overhead (file watches)
- Cross-platform challenges

**Decision**: Rejected - too complex for benefit

### Alternative 3: Rebuild on Every Command

**Pros**:
- Always fresh
- Simple logic
- No stale detection needed

**Cons**:
- Slow (1-3 seconds per command)
- Terrible user experience
- Defeats purpose of cache

**Decision**: Rejected - cache exists to be fast

### Alternative 4: TTL-Based Cache Invalidation

**Example**: Rebuild every 5 minutes

**Pros**:
- Simple to implement
- Predictable rebuild schedule

**Cons**:
- Arbitrary TTL value
- Could miss updates within TTL
- Could rebuild unnecessarily
- Not responsive to actual changes

**Decision**: Rejected - mtime detection is more precise

### Alternative 5: Hash-Based Change Detection

**Example**: Store file content hashes, compare on each command

**Pros**:
- Detects actual content changes
- No false positives from touch

**Cons**:
- Must read and hash every file
- Slower than mtime check
- More complex
- Overkill for this use case

**Decision**: Rejected - mtime is good enough and faster

### Alternative 6: Git Hook Integration

**Example**: Rebuild on post-merge, post-checkout hooks

**Pros**:
- Proactive rebuild
- Fast queries (no stale check)

**Cons**:
- Requires git hook installation
- Doesn't handle manual edits
- Not all operations trigger hooks
- Complex setup

**Decision**: Rejected - doesn't cover all cases, adds setup complexity

## References

- ADR 0001: Use SQLite as Cache Layer
- RFC 0001: Rex Documentation Management System
- SQLite Metadata Table: Stores last_rebuild timestamp
- Go os.Stat: https://pkg.go.dev/os#Stat
- Filesystem mtime: Modified time for change detection
