# 0003. Use Cobra for CLI Framework

Date: 2025-11-05

## Status

Accepted

## Context

Rex is being refactored from bash scripts to a Go CLI tool. We need a CLI framework that provides:
- Command structure with subcommands (`rex adr create`, `rex task list`)
- Flag parsing with type safety (`--type core --status in_progress`)
- Automatic help generation (`rex help`, `rex adr --help`)
- Input validation and error messages
- POSIX-compliant flag handling
- Completion support (bash, zsh, fish)
- Good developer experience and community support

The CLI should feel professional and follow conventions that developers expect from modern CLI tools like `git`, `docker`, and `kubectl`.

## Decision

We will use Cobra (https://github.com/spf13/cobra) as the CLI framework for rex:

1. **Command Structure**: Organize commands using Cobra's command tree structure
2. **Flag Handling**: Use Cobra/pflag for robust flag parsing
3. **Help Generation**: Let Cobra auto-generate help text from command definitions
4. **Viper Integration**: Use Viper (Cobra companion) for configuration file support
5. **Completion**: Enable shell completion for all major shells
6. **Conventions**: Follow Cobra best practices for command naming and structure

**Command Structure Example**:
```
rex
├── init
├── rebuild
├── adr
│   ├── create
│   ├── list
│   └── update
├── rfc
│   ├── create
│   ├── list
│   └── update
├── task
│   ├── create
│   ├── complete
│   ├── list
│   ├── stats
│   └── update
└── plan
    ├── create
    └── update
```

## Consequences

### Positive

- **Proven**: Used by kubectl, GitHub CLI, Hugo, and many production tools
- **Rich Features**: Commands, subcommands, flags, persistent flags, aliases
- **Auto Help**: Generates help text automatically from command definitions
- **Completion**: Built-in shell completion for bash, zsh, fish, PowerShell
- **Viper Integration**: Seamless config file support via Viper
- **Documentation**: Excellent documentation and examples
- **Community**: Large community, active maintenance, well-tested
- **Flag Parsing**: Robust pflag library for POSIX and GNU-style flags
- **Conventions**: Encourages CLI best practices and consistency
- **Testing**: Good support for testing CLI commands

### Negative

- **Dependency**: Adds external dependency (~100KB compiled)
- **Opinionated**: Must follow Cobra patterns and conventions
- **Learning Curve**: Developers need to learn Cobra API (moderate)
- **Over-Engineering**: May be more than needed for simple CLI (rex is not simple)
- **Flag Complexity**: Persistent flags vs. local flags can be confusing initially

### Neutral

- **Code Structure**: Encourages one file per command (more files but clearer)
- **Init Code**: Requires some boilerplate for command registration
- **Versioning**: Must manage Cobra dependency version (standard Go practice)

## Alternatives Considered

### Alternative 1: Standard Library flag Package

**Pros**:
- No dependencies
- Simple and minimal
- Part of standard library

**Cons**:
- No subcommands support
- Basic flag parsing only
- No help generation
- No completion support
- Manual implementation of everything

**Decision**: Rejected - too basic for rex's needs

### Alternative 2: urfave/cli

**Pros**:
- Popular alternative to Cobra
- Simpler API than Cobra
- Good documentation

**Cons**:
- Less feature-rich than Cobra
- Smaller community
- Weaker completion support
- Less integration with config libraries

**Decision**: Rejected - Cobra is more mature and feature-complete

### Alternative 3: spf13/pflag with Custom Framework

**Pros**:
- Use pflag for parsing
- Custom command structure
- Full control

**Cons**:
- Must implement subcommand routing
- Must implement help generation
- Must implement completion
- Reinventing the wheel
- More maintenance burden

**Decision**: Rejected - not worth the effort to build custom framework

### Alternative 4: alecthomas/kong

**Pros**:
- Tag-based struct definition
- Less boilerplate
- Modern approach

**Cons**:
- Smaller community
- Less mature than Cobra
- Different paradigm (learning curve)
- Limited completion support

**Decision**: Rejected - prefer Cobra's maturity and adoption

### Alternative 5: go-arg

**Pros**:
- Minimal boilerplate
- Struct-tag based
- Simple API

**Cons**:
- No subcommands support
- No completion
- Too simple for rex

**Decision**: Rejected - insufficient features

## References

- Cobra GitHub: https://github.com/spf13/cobra
- Cobra Documentation: https://cobra.dev/
- Viper (Config): https://github.com/spf13/viper
- pflag: https://github.com/spf13/pflag
- RFC 0001: Rex Documentation Management System
- Example Cobra Apps: kubectl, GitHub CLI (gh), Hugo
