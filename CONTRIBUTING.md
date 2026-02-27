# Contributing to kube-disk-stats

Thank you for your interest in contributing to kube-disk-stats! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:
- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive criticism

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- A clear and descriptive title
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, Kubernetes version)
- Any relevant logs or error messages

### Suggesting Enhancements

Enhancement suggestions are welcome! Please:

- Use a clear and descriptive title
- Provide a detailed description of the enhancement
- Explain why this enhancement would be useful
- Provide examples if applicable

### Pull Requests

#### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linter (`make lint`)
6. Commit your changes (`git commit -m 'Add some amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

#### Pull Request Guidelines

- Keep PRs focused and limited in scope
- Write clear commit messages
- Update documentation if needed
- Add tests for new features
- Ensure all tests pass
- Follow the existing code style

## Development Setup

### Prerequisites

- Go 1.23 or higher
- Docker (optional, for containerized development)
- Make (optional, for building)

### Setting Up

```bash
# Clone your fork
git clone https://github.com/kube-disk-stats.git
cd kube-disk-stats

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

## Code Style

- Follow Go conventions and best practices
- Use `gofmt` for formatting
- Write meaningful commit messages
- Add comments for complex logic
- Keep functions focused and small

## Testing

- Write tests for new features
- Ensure all tests pass before submitting a PR
- Test on multiple platforms if possible
- Test against different Kubernetes versions if applicable

## Documentation

- Update README.md for user-facing changes
- Add inline comments for complex code
- Update examples in documentation

## Release Process

Releases are managed via git tags and GitHub Actions. To create a release:

```bash
# Tag the release
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions workflow will:
1. Build binaries for all platforms
2. Create a GitHub release
3. Build and push Docker images

## Questions?

Feel free to open an issue or start a discussion if you have questions about contributing!

Thank you for contributing to kube-disk-stats! 🎉
