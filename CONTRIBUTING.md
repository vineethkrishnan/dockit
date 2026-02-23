# Contributing to dockit

First off, thank you for considering contributing to `dockit`! It's people like you that make `dockit` such a great tool.

## Getting Started

### Prerequisites

*   Go 1.23+ installed.
*   Docker installed and running locally.
*   `make` utility.
*   [golangci-lint](https://golangci-lint.run/) (optional, but our CI uses it).

### Local Development Setup

1.  **Fork** the repository on GitHub.
2.  **Clone** your fork locally:
    ```bash
    git clone https://github.com/your-username/dockit.git
    cd dockit
    ```
3.  **Install dependencies**:
    ```bash
    go mod download
    ```
4.  **Run tests** to ensure everything is working:
    ```bash
    make test
    ```
5.  **Build** the project:
    ```bash
    make build
    ```

## Development Workflow

1.  Create a descriptive branch for your feature or bugfix:
    ```bash
    git checkout -b feat/add-log-rotation-suggestion
    ```
2.  Make your changes.
3.  Format and lint your code:
    ```bash
    make lint
    # or just run go fmt ./...
    ```
4.  Run tests:
    ```bash
    make test
    ```
5.  Commit your changes using conventional commits (e.g., `feat: added log rotation`, `fix: handled nil pointer in analyzer`).
6.  Push to your fork and submit a Pull Request.

## Pull Request Guidelines

*   Ensure all tests pass. If you're adding a new feature, add a test for it.
*   Update documentation if you are changing CLI flags or behaviors.
*   Your PR should ideally address a single concern.
*   The project uses `release-please` for automated changelog generation and semantic versioning, which is why following [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) is highly recommended.

Once your PR is submitted, CI will run checks. Note that we require approvals and passing CI before merging.

Thank you for contributing!
