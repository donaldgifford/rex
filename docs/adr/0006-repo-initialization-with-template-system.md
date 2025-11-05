# 0006. Repo Initialization with Template System

Date: 2025-11-05

## Status

Proposed

## Context

Rex currently focuses on documentation management (ADRs, RFCs, Tasks, Plans). However, setting up a new repository involves much more than just documentation structure:

- **CI/CD**: GitHub Actions workflows, GitLab CI, CircleCI configs
- **Build Tools**: Makefiles, Taskfiles, build scripts
- **Containers**: Dockerfiles, docker-compose.yml, .dockerignore
- **Linting**: .golangci.yml, .eslintrc, .prettierrc, markdownlint configs
- **Git Config**: .gitignore, .gitattributes, .editorconfig
- **Development**: mise.toml, .nvmrc, .tool-versions
- **Package Management**: go.mod, package.json, Cargo.toml
- **Testing**: Test configs, coverage settings
- **Security**: Dependabot, CodeQL, security scanning configs

Current pain points:
- Developers copy-paste config files from other repos
- Config files can be outdated or inconsistent
- No standard "blessed" setup for team projects
- Manual setup takes 30+ minutes for a new repo
- Different projects use different tooling (inconsistent)

Users need a way to:
1. Initialize standard repo tooling alongside documentation
2. Choose from community templates or company-specific templates
3. Keep templates updated without manual copying
4. Customize templates for specific project needs

## Decision

We will extend Rex to support comprehensive repository initialization with a dual-template system:

### 1. Embedded Templates (Built-in Defaults)

Use Go embed for standard, opinionated templates bundled with rex binary:

```
templates/
├── docs/           # Documentation (already decided)
│   ├── adr/
│   ├── rfc/
│   ├── tasks/
│   └── plans/
├── github/         # GitHub-specific
│   ├── workflows/
│   │   ├── ci.yml
│   │   ├── release.yml
│   │   └── dependabot.yml
│   └── CODEOWNERS
├── configs/        # Common config files
│   ├── .golangci.yml
│   ├── .markdownlint.yaml
│   ├── .editorconfig
│   ├── .gitignore
│   └── .prettierrc.yaml
├── docker/         # Container templates
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── .dockerignore
└── makefiles/      # Build automation
    ├── Makefile
    └── Taskfile.yml
```

### 2. Remote Templates (via go-getter)

Use HashiCorp's go-getter library to fetch templates from remote sources:

```bash
# Fetch from Git repository
rex init --template=github.com/company/repo-templates

# Fetch from specific branch/tag
rex init --template=github.com/company/repo-templates?ref=v1.2.0

# Fetch subdirectory
rex init --template=github.com/company/repo-templates//golang

# Multiple sources (S3, HTTP, local)
rex init --template=s3::https://s3.amazonaws.com/bucket/templates
rex init --template=https://example.com/templates.tar.gz
rex init --template=file:///local/path/templates
```

### 3. Template Selection System

Allow users to choose which templates to initialize:

```bash
# Interactive mode - prompt for each category
rex init

# Minimal - docs only (current behavior)
rex init --minimal

# Full - all embedded templates
rex init --full

# Selective - choose specific categories
rex init --docs --github --docker

# Remote template
rex init --template=github.com/myorg/templates

# Combine embedded + remote (remote overrides)
rex init --github --template=github.com/myorg/gh-actions
```

### 4. Template Overrides

Support `.rex.yaml` config for default template sources:

```yaml
# .rex.yaml
templates:
  source: github.com/myorg/templates
  categories:
    github: github.com/myorg/gh-actions
    docker: github.com/myorg/docker-templates
  embedded:
    - docs
    - configs
```

### 5. Template Variables

Support variable substitution in templates:

```yaml
# Template file: github/workflows/ci.yml
name: CI

on:
  push:
    branches: [{{ .MainBranch }}]

jobs:
  test:
    runs-on: {{ .RunnerOS }}
    steps:
      - uses: actions/checkout@v4
      - name: Run tests
        run: {{ .TestCommand }}
```

```bash
# Provide variables
rex init --var MainBranch=main --var RunnerOS=ubuntu-latest --var TestCommand="go test ./..."
```

### 6. Implementation Phases

**Not in Phase 1-5** (Core documentation features first)

**Phase 6: Template System Foundation** (2-3 weeks)
- Integrate go-getter library
- Implement template fetching and caching
- Add template variable substitution
- Create category selection system
- Implement embedded template categories

**Phase 7: Standard Template Library** (2 weeks)
- Create opinionated templates for:
  - GitHub Actions (CI, release, dependabot)
  - Common configs (.gitignore, .editorconfig, linters)
  - Docker templates
  - Makefile templates
- Document template structure
- Add template validation

**Phase 8: Remote Template Support** (1 week)
- Full go-getter integration
- S3, HTTP, Git source support
- Template versioning and caching
- `.rex.yaml` template configuration

## Consequences

### Positive

- **Comprehensive Setup**: One command to initialize entire repo structure
- **Consistency**: Standard templates across team/organization
- **Flexibility**: Support both built-in and custom templates
- **Updatable**: Remote templates can be updated without rex binary updates
- **Extensible**: Organizations can host their own template repos
- **DRY**: No more copying config files between repos
- **Best Practices**: Embed opinionated, well-tested configs
- **Time Savings**: 30+ minutes of setup → 2 minutes with `rex init --full`
- **Onboarding**: New developers get correct setup immediately

### Negative

- **Scope Creep**: Rex becomes more than documentation tool
- **Maintenance**: Must maintain template library
- **Complexity**: More templates = more code paths to test
- **Breaking Changes**: Template updates could break existing workflows
- **Network Dependency**: Remote templates require internet (cached after first fetch)
- **Go-getter Dependency**: Adds external dependency (~500KB compiled)
- **Template Quality**: Must ensure embedded templates are high-quality
- **Documentation Burden**: Must document template structure and variables

### Neutral

- **Opinionated vs. Flexible**: Balance between "blessed" templates and customization
- **Template Versioning**: Remote templates need versioning strategy
- **Cache Management**: Go-getter caches templates - need cleanup command
- **Override Behavior**: How should embedded + remote templates merge?
- **Template Discovery**: How do users find available templates?

## Alternatives Considered

### Alternative 1: Only Embedded Templates

**Pros**:
- Simpler implementation
- No network dependency
- Always available offline
- No go-getter dependency

**Cons**:
- Cannot use custom/company templates
- Templates locked to rex version
- No extensibility for organizations
- All templates bloat binary

**Decision**: Rejected - organizations need custom templates

### Alternative 2: Only Remote Templates

**Pros**:
- No template maintenance in rex
- Organizations have full control
- Smaller binary size

**Cons**:
- Requires network access always
- No opinionated defaults
- Poor onboarding (must know template URL)
- Inconsistent experience

**Decision**: Rejected - need good defaults for new users

### Alternative 3: Git Submodules for Templates

**Pros**:
- Standard Git mechanism
- Version controlled
- Well understood

**Cons**:
- Requires git repo to be initialized
- Complex submodule management
- Not beginner-friendly
- Doesn't support S3/HTTP sources

**Decision**: Rejected - too complex, limited sources

### Alternative 4: Package Manager Approach (like npm init)

**Pros**:
- Familiar to developers
- Rich template ecosystem possible

**Cons**:
- Requires package registry infrastructure
- More complex than needed
- Vendor lock-in to registry

**Decision**: Rejected - go-getter is simpler and more flexible

### Alternative 5: Template as Code (Programmatic)

**Example**: Write templates in Go code, not files

**Pros**:
- Type-safe
- Compile-time validation
- No parsing needed

**Cons**:
- Not human-friendly
- Cannot be externally hosted
- Hard to customize
- Requires Go knowledge

**Decision**: Rejected - templates should be accessible files

### Alternative 6: Template Registry Service

**Pros**:
- Centralized discovery
- Template ratings/reviews
- Version management

**Cons**:
- Requires infrastructure
- Network dependency
- Complexity overhead
- Not needed for MVP

**Decision**: Rejected - can add later if needed

## Implementation Details

### Go-getter Integration

```go
import "github.com/hashicorp/go-getter"

func FetchTemplate(source string, dest string) error {
    client := &getter.Client{
        Src:  source,
        Dst:  dest,
        Mode: getter.ClientModeAny,
    }
    return client.Get()
}
```

### Template Structure

```
template-repo/
├── template.yaml        # Metadata
├── README.md           # Documentation
├── github/
│   └── workflows/
│       └── ci.yml
├── docker/
│   └── Dockerfile
└── configs/
    └── .gitignore
```

**template.yaml**:
```yaml
name: "Go Microservice Template"
version: "1.0.0"
description: "Standard template for Go microservices"
variables:
  - name: MainBranch
    description: "Main branch name"
    default: "main"
  - name: GoVersion
    description: "Go version"
    default: "1.21"
categories:
  - github
  - docker
  - configs
```

### Template Cache

```
~/.rex/cache/templates/
└── github.com/
    └── company/
        └── templates/
            └── 1.0.0/
                ├── template.yaml
                └── ...
```

### CLI Commands

```bash
# Initialize with defaults
rex init

# Full initialization
rex init --full

# Selective categories
rex init --github --docker

# Remote template
rex init --template=github.com/org/templates

# With variables
rex init --var GoVersion=1.21

# List available embedded templates
rex template list

# Show template details
rex template show github.com/org/templates

# Clear template cache
rex template cache clear

# Update cached templates
rex template cache update
```

## Future Enhancements

1. **Template Inheritance**: Allow templates to extend/override other templates
2. **Interactive Wizard**: TUI for template selection and variable input
3. **Template Validation**: Validate templates before applying
4. **Dry Run**: Preview what would be created (`rex init --dry-run`)
5. **Template Marketplace**: Community registry of templates
6. **AI-Generated Templates**: LLM to generate custom templates based on description
7. **Template Testing**: Framework for testing template outputs
8. **Hooks**: Run scripts after template initialization

## References

- go-getter: https://github.com/hashicorp/go-getter
- GitHub Actions Templates: https://github.com/actions/starter-workflows
- Cookiecutter (Python): https://github.com/cookiecutter/cookiecutter
- Yeoman (Node): https://yeoman.io/
- cargo-generate (Rust): https://github.com/cargo-generate/cargo-generate
- RFC 0001: Rex Documentation Management System
- ADR 0002: Use Go Embed for Templates

## Notes

- **Phase Timing**: This is a future enhancement (Phase 6+), not part of core MVP (Phases 1-5)
- **Go-getter Support**: Supports Git, Mercurial, HTTP, S3, GCS, and local files
- **Security**: Remote templates should be verified (checksums, signatures) in future
- **Conflicts**: Need strategy for merging templates with existing files
- **Rollback**: Should support undoing template application if something goes wrong
