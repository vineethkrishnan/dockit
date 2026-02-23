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
```text
--- CONTAINERS ---
ID/NAME              STATE           SCORE           SIZE       REASON
reverent_raman      exited          SAFE            400 B      Container is stopped and old
angry_murdock       running         PROTECTED       12 MB      Container is currently active

--- IMAGES ---
REPO:TAG             DANGLING   SCORE           SIZE       REASON
postgres:15          false      PROTECTED       379 MB     Image is currently backing a container
<none>               true       SAFE            1.2 GB     Image is dangling and old
```

### Disk Summary

Get a quick high-level overview.

```bash
dockit summary
```
```text
TYPE            TOTAL           ACTIVE          SIZE            RECLAIMABLE
Images          2               1               1.6 GB          1.2 GB
Containers      2               1               12 MB           400 B
Local Volumes   0               0               0 B             0 B
Build Cache     0               0               0 B             0 B
--------------------------------------------------------------------------------
Total Space:                                    1.6 GB
Reclaimable Space:                              1.2 GB
```

### Safe Cleanup

Reclaim space safely. By default, `dockit clean` performs a dry-run. 

```bash
dockit clean
```
```text
--- DRY RUN: Cleanup Plan ---
dockit identified 2 resources that are eligible for deletion.
Only SAFE and REVIEW items are included. PROTECTED items are ignored.
Total space that would be reclaimed: 1.2 GB

TYPE            ID/NAME                                  SCORE      SIZE
Container       reverent_raman                           SAFE       400 B
Image           <none>                                   SAFE       1.2 GB

To apply these deletions, run: dockit clean --apply
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
```text
--- CONTAINER LOG SIZES (Total: 4.2 GB) ---
CONTAINER            SIZE            WARNINGS
angry_murdock        4.2 GB          🚨 EXCESSIVE - Consider adding 'log-opt max-size=10m'
reverent_raman       24 B            
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to set up your development environment, run tests, and submit Pull Requests. Remember to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
