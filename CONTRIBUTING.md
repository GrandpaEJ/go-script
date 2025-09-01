# Contributing to Go-Script

Thank you for your interest in contributing to Go-Script! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/go-script.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit and push your changes
7. Create a pull request

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/GrandpaEJ/go-script.git
cd go-script

# Build the project
go build -o gos ./cmd/gos

# Run tests
go test ./tests/...

# Test with examples
./gos run examples/hello.gos
```

## Project Structure

```
go-script/
├── cmd/gos/           # CLI application
├── pkg/
│   ├── ast/           # Abstract Syntax Tree definitions
│   ├── lexer/         # Lexical analysis
│   ├── parser/        # Syntax analysis
│   ├── codegen/       # Go code generation
│   ├── runtime/       # Runtime support
│   └── stdlib/        # Standard library
├── examples/          # Example programs
├── tests/            # Test suite
└── docs/             # Documentation
```

## Contributing Guidelines

### Code Style

- Follow Go conventions and use `gofmt`
- Write clear, descriptive variable and function names
- Add comments for public APIs and complex logic
- Keep functions focused and reasonably sized

### Testing

- Write tests for new features
- Ensure all existing tests pass
- Add integration tests for end-to-end functionality
- Test edge cases and error conditions

### Documentation

- Update README.md for new features
- Add examples for new language constructs
- Update LANGUAGE_SPEC.md for syntax changes
- Include inline code documentation

## Types of Contributions

### Bug Reports

When reporting bugs, please include:

- Go-Script version
- Operating system and version
- Minimal code example that reproduces the issue
- Expected vs actual behavior
- Error messages (if any)

### Feature Requests

For new features, please:

- Check if the feature already exists or is planned
- Describe the use case and motivation
- Provide examples of how it would be used
- Consider backward compatibility

### Code Contributions

#### Language Features

- New syntax constructs
- Built-in functions
- Type system improvements
- Error handling enhancements

#### Tooling

- CLI improvements
- Build system enhancements
- IDE integration
- Debugging tools

#### Performance

- Lexer/parser optimizations
- Code generation improvements
- Runtime performance
- Memory usage optimization

## Development Workflow

### Adding a New Language Feature

1. **Design**: Update `LANGUAGE_SPEC.md` with the new syntax
2. **Lexer**: Add new tokens if needed in `pkg/lexer/token.go`
3. **Parser**: Implement parsing logic in `pkg/parser/parser.go`
4. **AST**: Add new AST nodes in `pkg/ast/ast.go`
5. **Codegen**: Implement Go code generation in `pkg/codegen/generator.go`
6. **Tests**: Add comprehensive tests
7. **Examples**: Create example programs
8. **Documentation**: Update relevant documentation

### Testing Your Changes

```bash
# Run unit tests
go test ./pkg/lexer/
go test ./pkg/parser/
go test ./pkg/codegen/

# Run integration tests
go test ./tests/

# Test with examples
./gos run examples/hello.gos
./gos run examples/calculator.gos
./gos run examples/structs.gos

# Test CLI commands
./gos version
./gos help
```

### Submitting Changes

1. Ensure all tests pass
2. Update documentation
3. Add examples if applicable
4. Write clear commit messages
5. Create a pull request with:
   - Clear description of changes
   - Motivation for the changes
   - Testing performed
   - Any breaking changes

## Code Review Process

- All changes require review before merging
- Reviewers will check for:
  - Code quality and style
  - Test coverage
  - Documentation updates
  - Backward compatibility
  - Performance implications

## Release Process

1. Version bump in relevant files
2. Update CHANGELOG.md
3. Tag the release
4. Build and test release binaries
5. Update documentation

## Community

- Be respectful and inclusive
- Help others learn and contribute
- Share knowledge and best practices
- Provide constructive feedback

## Getting Help

- Check existing issues and documentation
- Ask questions in GitHub discussions
- Join community channels (if available)
- Reach out to maintainers

## License

By contributing to Go-Script, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to Go-Script!
