# Contributing to GoAstra

Thank you for your interest in contributing to GoAstra! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/goastra.git
   cd goastra
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/channdev/goastra.git
   ```

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+ (optional, for database features)

### Building the CLI

```bash
cd cli
go build -o goastra.exe ./goastra
```

### Running Tests

```bash
# Backend tests
cd app
go test ./...

# Frontend tests
cd web
npm test
```

## Making Changes

1. Create a new branch for your feature or fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit them:
   ```bash
   git add .
   git commit -m "Add your descriptive commit message"
   ```

3. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

4. Open a Pull Request on GitHub

## Code Style

### Go Code
- Follow standard Go conventions
- Run `go fmt` before committing
- Run `go vet` to check for issues

### TypeScript/Angular Code
- Follow Angular style guide
- Use TypeScript strict mode
- Run `npm run lint` before committing

## Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Write clear commit messages
- Update documentation if needed
- Add tests for new features
- Ensure all tests pass before submitting

## Reporting Issues

When reporting issues, please include:

- GoAstra version (`goastra version`)
- Go version (`go version`)
- Node.js version (`node --version`)
- Operating system
- Steps to reproduce the issue
- Expected vs actual behavior

## Feature Requests

Feature requests are welcome! Please open an issue with:

- Clear description of the feature
- Use case and benefits
- Any implementation ideas (optional)

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn and grow

## License

By contributing to GoAstra, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue or reach out to [@channdev](https://github.com/channdev) on GitHub.

Thank you for contributing!
