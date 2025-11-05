# Rex MVP Guide

This guide defines the Minimum Viable Product (MVP) for the Rex refactor and provides a clear path to a usable v1.0 release.

## MVP Definition

**Goal**: Replace bash scripts with a Go CLI tool that provides feature parity for core documentation management (ADRs, RFCs, Tasks).

**Timeline**: 8.5 weeks (Phases 1-5 from RFC 0001)

**Success Criteria**:
- ✅ All bash script functionality replaced with rex CLI commands
- ✅ SQLite cache provides fast queries (< 100ms for list operations)
- ✅ Zero external dependencies at runtime (no yq, sed, awk)
- ✅ Cross-platform binaries (Linux, macOS, Windows)
- ✅ Comprehensive documentation and migration guide

## MVP Scope (In)

### Core Features

**Documentation Management**:
- ✅ ADR creation, listing, README generation
- ✅ RFC creation, listing, README generation
- ✅ Task creation, completion, listing, statistics
- ✅ Plan creation, listing, README generation

**Performance**:
- ✅ SQLite cache for fast queries
- ✅ Auto-rebuild on stale detection
- ✅ Manual rebuild command

**Developer Experience**:
- ✅ Interactive prompts for task creation
- ✅ Rich filtering (type, status, priority, tags)
- ✅ Statistics and aggregations
- ✅ Clear error messages
- ✅ Built-in help for all commands

**Distribution**:
- ✅ Single binary distribution
- ✅ `go install` support
- ✅ Direct binary downloads
- ✅ Cross-platform support

## MVP Scope (Out)

These features are explicitly **not** in the MVP:

### Post-MVP Features (Future Phases)

**Template System** (Phase 6-8):
- ❌ Repo initialization with GitHub Actions, Docker, Makefiles
- ❌ Remote template fetching (go-getter)
- ❌ Custom template support
- ❌ Template variables and customization

**Advanced Features**:
- ❌ HTML generation from markdown
- ❌ Export to PDF/Confluence
- ❌ Git hook integration
- ❌ Task time tracking (start/stop timers)
- ❌ Full-text search (FTS5)
- ❌ Task dependencies visualization
- ❌ Interactive TUI mode

**Nice-to-Have**:
- ❌ Shell completion scripts
- ❌ Config file validation
- ❌ Dry-run mode
- ❌ Undo/rollback functionality

## Phase Breakdown

### Phase 1: Foundation & Database (3 weeks)

**Week 1: Core Setup**
```bash
# Deliverables
✓ Go module initialized with dependencies
✓ SQLite schema implemented and tested
✓ Database connection management
✓ Schema migrations
```

**Week 2: Parsing & Templates**
```bash
# Deliverables
✓ YAML frontmatter parser
✓ Markdown extraction utilities
✓ Go embed for templates
✓ Template rendering system
```

**Week 3: ADR Commands**
```bash
# Deliverables
✓ rex adr create
✓ rex adr list
✓ rex adr update (README generation)
✓ Auto-rebuild implementation
✓ rex rebuild command
```

**Milestone 1 Demo**:
```bash
rex init
rex adr create "Use SQLite as Cache"
rex adr list
rex rebuild
```

### Phase 2: RFC Support (1 week)

**Deliverables**:
```bash
✓ rex rfc create
✓ rex rfc list
✓ rex rfc update
✓ RFC database schema
✓ RFC parsing logic
```

**Milestone 2 Demo**:
```bash
rex rfc create "Rex Documentation System"
rex rfc list --status Draft
rex rfc update
```

### Phase 3: Task Management (2 weeks)

**Week 1: Task CRUD**
```bash
# Deliverables
✓ rex task create (interactive prompts)
✓ rex task list (with filtering)
✓ Task database schema with relationships
✓ Task parsing with frontmatter
```

**Week 2: Task Operations**
```bash
# Deliverables
✓ rex task complete
✓ rex task stats
✓ rex task update
✓ File movement (active/ to completed/)
✓ Relationship handling
```

**Milestone 3 Demo**:
```bash
rex task create
rex task list --type core --status in_progress --priority P1
rex task stats
rex task complete TASK-001
```

### Phase 4: Plan Support & Polish (1.5 weeks)

**Deliverables**:
```bash
✓ rex plan create
✓ rex plan list
✓ rex plan update
✓ rex version command
✓ rex help with examples
✓ rex db info (show cache stats)
✓ Error messages and validation
✓ Binary releases (GoReleaser)
✓ DB migration system
```

**Milestone 4 Demo**:
```bash
rex plan create "Q4 Roadmap"
rex plan list
rex version
rex db info
```

### Phase 5: Documentation & Migration (1 week)

**Deliverables**:
```bash
✓ Migration guide for bash script users
✓ Updated README.md
✓ Command reference documentation
✓ Workflow examples and demos
✓ Bash script deprecation notices
✓ FAQ and troubleshooting guide
✓ Video/GIF demos
```

**Milestone 5 (MVP Complete)**:
- All bash scripts replaced
- Documentation complete
- Binaries released for all platforms
- Migration path clear
- Ready for v1.0 release

## Critical Path

These items **must** be completed in order:

1. **Database Schema** → Everything depends on this
2. **Parser** → Required for reading existing files
3. **Auto-rebuild** → Core feature, needed early
4. **ADR Commands** → Simplest type, proves architecture
5. **RFC Commands** → Similar to ADR, builds confidence
6. **Task Commands** → Most complex, saved for when patterns are established

## Development Workflow

### Daily Workflow

```bash
# 1. Pick a task from the current phase
# 2. Create feature branch
git checkout -b feature/adr-create

# 3. Implement with tests
go test ./...

# 4. Test locally
go build -o rex main.go
./rex adr create "Test ADR"

# 5. Commit and push
git add .
git commit -m "Implement rex adr create command"
git push origin feature/adr-create

# 6. Create PR and merge
```

### Testing Checklist

Before marking a phase complete:

- [ ] All unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed
- [ ] Documentation updated
- [ ] No known bugs
- [ ] Performance benchmarks met

### Definition of Done

A feature is "done" when:

1. ✅ Code implemented and reviewed
2. ✅ Unit tests written and passing
3. ✅ Integration tests passing
4. ✅ Documentation updated
5. ✅ Manual testing completed
6. ✅ No regressions in existing features
7. ✅ Performance meets targets

## MVP Risks & Mitigations

### High Risk

**Risk**: SQLite performance doesn't meet targets with 100+ tasks

**Mitigation**:
- Benchmark early (Week 1)
- Optimize indexes
- Consider caching query results
- Acceptable: 200ms instead of 100ms

**Risk**: Go embed increases binary size too much

**Mitigation**:
- Compress templates
- Only embed essential templates
- Target: < 15MB binary

### Medium Risk

**Risk**: Complex task relationships cause bugs

**Mitigation**:
- Thorough testing of relationships
- Validate relationships on creation
- Provide clear error messages

**Risk**: Cross-platform issues

**Mitigation**:
- Test on all platforms early
- Use GoReleaser for consistent builds
- Avoid platform-specific code

### Low Risk

**Risk**: Migration from bash scripts is difficult

**Mitigation**:
- Provide migration script
- Comprehensive migration guide
- `rex rebuild` handles existing files

## Performance Targets

### Must Meet (MVP)

- ADR/RFC creation: < 200ms
- Task creation: < 300ms (interactive prompts)
- List operations: < 150ms (100 items)
- Stats operations: < 300ms (100 items)
- Cache rebuild: < 3 seconds (100 files)

### Stretch Goals (Post-MVP)

- ADR/RFC creation: < 100ms
- Task creation: < 200ms
- List operations: < 100ms
- Stats operations: < 200ms
- Cache rebuild: < 2 seconds

## Binary Size Targets

- **Target**: < 15MB
- **Maximum acceptable**: < 20MB
- **Stretch goal**: < 10MB

**Size breakdown**:
- Go runtime: ~2MB
- Cobra: ~1MB
- SQLite: ~3MB
- Templates: ~500KB
- Rex code: ~3MB

## MVP Feature Matrix

| Feature | Phase | Priority | Complexity | Status |
|---------|-------|----------|------------|--------|
| Database schema | 1 | Critical | Medium | 🔴 Not started |
| YAML parser | 1 | Critical | Low | 🔴 Not started |
| Go embed templates | 1 | Critical | Low | 🔴 Not started |
| Auto-rebuild | 1 | Critical | Medium | 🔴 Not started |
| rex adr create | 1 | Critical | Medium | 🔴 Not started |
| rex adr list | 1 | Critical | Low | 🔴 Not started |
| rex rfc create | 2 | Critical | Medium | 🔴 Not started |
| rex rfc list | 2 | Critical | Low | 🔴 Not started |
| rex task create | 3 | Critical | High | 🔴 Not started |
| rex task list | 3 | Critical | Medium | 🔴 Not started |
| rex task stats | 3 | Critical | Medium | 🔴 Not started |
| rex task complete | 3 | Critical | Medium | 🔴 Not started |
| rex plan create | 4 | High | Low | 🔴 Not started |
| Binary releases | 4 | High | Low | 🔴 Not started |
| Migration guide | 5 | High | Low | 🔴 Not started |

Legend:
- 🔴 Not started
- 🟡 In progress
- 🟢 Complete

## Success Metrics

### MVP Launch Criteria

All must be ✅ to launch v1.0:

- [ ] All critical features implemented
- [ ] All tests passing (unit + integration)
- [ ] Performance targets met
- [ ] Binaries built for Linux, macOS, Windows
- [ ] Documentation complete
- [ ] Migration guide tested
- [ ] Zero known critical bugs
- [ ] At least one team successfully migrated from bash scripts

### Post-Launch Metrics

Track these after v1.0 release:

- **Adoption**: Number of repos using rex
- **Feedback**: GitHub issues and feature requests
- **Performance**: Real-world usage metrics
- **Stability**: Bug reports and crash rates

## MVP Delivery Schedule

**Optimistic** (all goes well): 7 weeks

**Realistic** (normal issues): 8.5 weeks

**Pessimistic** (major blockers): 10 weeks

### Week-by-Week Timeline

| Week | Phase | Deliverables |
|------|-------|-------------|
| 1 | 1 | Database + Connection |
| 2 | 1 | Parser + Templates |
| 3 | 1 | ADR Commands + Auto-rebuild |
| 4 | 2 | RFC Support |
| 5 | 3 | Task CRUD |
| 6 | 3 | Task Complete + Stats |
| 7 | 4 | Plan Support + Polish |
| 8 | 5 | Documentation |
| 9 | - | Buffer / Testing |

**Launch**: End of Week 8 (or Week 9 with buffer)

## Post-MVP Roadmap

After v1.0 launch, consider these features:

### v1.1 (Quick Wins)
- Shell completion (bash, zsh, fish)
- Config file validation
- Improved error messages
- Performance optimizations

### v1.2 (Enhanced Features)
- Full-text search (SQLite FTS5)
- Task dependency visualization
- Export to PDF/HTML

### v2.0 (Major Features)
- Template system (go-getter)
- Repo initialization
- GitHub Actions templates
- Docker templates

## Getting Started

Ready to build the MVP? Follow these steps:

1. **Read the RFC and ADRs**
   - RFC 0001: Overall architecture
   - ADR 0001-0005: Key decisions

2. **Review Implementation Guide**
   - Detailed technical instructions
   - Code examples
   - Testing strategies

3. **Set up development environment**
   ```bash
   go version  # Ensure Go 1.21+
   git clone github.com/donaldgifford/rex
   cd rex
   go mod download
   ```

4. **Start with Phase 1, Week 1**
   - Database schema
   - Connection management
   - Schema initialization

5. **Follow the critical path**
   - Don't skip ahead
   - Test thoroughly
   - Keep documentation updated

## Questions?

Refer to:
- **IMPLEMENTATION_GUIDE.md**: Detailed technical guide
- **RFC 0001**: Overall vision and architecture
- **ADRs**: Specific architectural decisions
- **Task Documentation**: docs/tasks/README.md

## References

- RFC 0001: Rex Documentation Management System
- All ADRs: docs/adr/
- Implementation Guide: docs/IMPLEMENTATION_GUIDE.md
- Current bash scripts: tools/ (for reference)
