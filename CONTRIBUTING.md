# Contributing to LiteKPI

Thank you for your interest in contributing to LiteKPI! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/devbydaniel/litekpi/issues)
2. If not, create a new issue using the bug report template
3. Include as much detail as possible: steps to reproduce, expected behavior, actual behavior, and your environment

### Suggesting Features

1. Check if the feature has already been requested in [Issues](https://github.com/devbydaniel/litekpi/issues)
2. If not, create a new issue using the feature request template
3. Describe the problem you're trying to solve and your proposed solution

### Pull Requests

1. Fork the repository
2. Create a new branch from `main` (`git checkout -b feature/your-feature`)
3. Make your changes
4. Run tests and linting (`make test && make lint`)
5. Commit your changes with a clear message
6. Push to your fork and submit a pull request

## Development Setup

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker](https://www.docker.com/) & Docker Compose
- [Air](https://github.com/cosmtrek/air) (for Go hot-reload)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/litekpi.git
cd litekpi

# Set up environment
cp .env.example .env

# Install dependencies
make install

# Start dev services (PostgreSQL + Mailcatcher)
make dev-services

# Run database migrations
make migrate

# Start development servers
make dev
```

### Available Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start all services with hot-reload |
| `make dev-backend` | Start backend only |
| `make dev-frontend` | Start frontend only |
| `make dev-services` | Start PostgreSQL + Mailcatcher |
| `make dev-stop` | Stop dev services |
| `make test` | Run all tests |
| `make lint` | Lint code |
| `make fmt` | Format code |
| `make migrate` | Run database migrations |
| `make api-gen` | Generate TypeScript API client |

### Project Structure

See [AGENTS.md](AGENTS.md) for detailed architecture documentation.

## Style Guidelines

### Go (Backend)

- Follow standard Go conventions
- Run `make fmt` before committing
- Run `make lint` to check for issues

### TypeScript (Frontend)

- Follow the existing code style
- Use TypeScript strict mode
- Prefer functional components with hooks

### Commits

- Use clear, descriptive commit messages
- Start with a verb in present tense (e.g., "Add feature", "Fix bug", "Update docs")
- Reference issues when applicable (e.g., "Fix login redirect (#123)")

## Testing

- Write tests for new features
- Ensure existing tests pass before submitting a PR
- Run `make test` to execute the test suite

## Questions?

If you have questions, feel free to open an issue with the "question" label.
