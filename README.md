# 🐳 dockit

[![CI](https://github.com/vineethkrishnan/dockit/actions/workflows/ci.yml/badge.svg)](https://github.com/vineethkrishnan/dockit/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

`dockit` is a safe, intelligent, audit-first Docker disk analysis and cleanup CLI built in Go.

It replaces the blind usage of `docker system prune -a` with a transparent, controlled, and risk-aware workflow. Know exactly what is taking up space, and safely reclaim it.

## 🌟 Features

*   **Full Visibility**: Get a complete breakdown of Docker disk usage (images, containers, volumes, build cache).
*   **Intelligent Scoring**: Resources are classified into `SAFE`, `REVIEW`, and `PROTECTED` to prevent accidental deletion of critical infrastructure.
*   **Log Detection**: Detect runaway container logs that are silently filling up your disk.
*   **Safe Cleanup**: Interactive prompt-based cleanup. Never delete running containers or attached volumes.
*   **Automation Ready**: Full JSON output support (`--json`) for CI pipelines and monitoring tools.

## 🚀 Installation

*Pre-built binaries will be available in the Releases page soon.*

To build from source:

```bash
git clone https://github.com/vineethkrishnan/dockit.git
cd dockit
make build
./dockit version
```

## 📖 Usage

### Deep Analysis

Audit your Docker environment to see where space is being used and what is safe to delete.

```bash
dockit analyze
```

### Disk Summary

Get a quick high-level overview.

```bash
dockit summary
```

### Safe Cleanup

Reclaim space safely. By default, `dockit clean` performs a dry-run. 

```bash
dockit clean
```

To actually apply the deletions:
```bash
dockit clean --apply
```

### Find Large Logs

Identify containers generating massive logs.

```bash
dockit logs
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to set up your development environment, run tests, and submit Pull Requests. Remember to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
