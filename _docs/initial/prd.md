# 🐳 dockit

## Product Requirements Document (PRD)

**Version:** 1.0\
**Language:** Go\
**Type:** CLI Tool

------------------------------------------------------------------------

# 1. Product Vision

`dockit` is a safe, intelligent, audit-first Docker disk analysis and
cleanup CLI built in Go.

It provides:

-   Full Docker disk visibility
-   Risk-aware cleanup recommendations
-   Log size detection
-   Safe interactive cleanup
-   JSON output for automation

It replaces blind usage of:

    docker system prune -a

with a transparent and controlled workflow.

------------------------------------------------------------------------

# 2. Problem Statement

Docker environments accumulate disk usage from:

-   Unused images
-   Stopped containers
-   Orphan volumes
-   Build cache
-   Large container logs

Common issues:

-   VPS disk full
-   CI runners failing
-   Servers crashing
-   No clear visibility into space usage
-   Over-aggressive pruning

There is no intelligent CLI that:

-   Audits first
-   Scores deletion risk
-   Analyzes log growth
-   Provides structured output

------------------------------------------------------------------------

# 3. Target Users

## Primary

-   DevOps engineers
-   Backend developers
-   Startup teams
-   VPS administrators

## Secondary

-   CI/CD maintainers
-   Homelab users
-   Docker-heavy SaaS operators

------------------------------------------------------------------------

# 4. Goals

## 4.1 Functional Goals

-   Provide disk usage breakdown
-   Detect reclaimable space
-   Classify safe-to-delete resources
-   Detect large Docker logs
-   Provide interactive cleanup
-   Support machine-readable output

## 4.2 Non-Goals (v1)

-   Kubernetes support
-   Multi-node cluster management
-   Docker Swarm orchestration
-   Cloud billing analysis

------------------------------------------------------------------------

# 5. Core Features (MVP)

## 5.1 Disk Summary

Command:

    dockit summary

Displays:

-   Images size
-   Containers size
-   Volumes size
-   Build cache size
-   Container logs size
-   Total usage
-   Reclaimable estimate

------------------------------------------------------------------------

## 5.2 Deep Analysis

Command:

    dockit analyze

Per-resource breakdown:

### Images

-   ID
-   Repo:Tag
-   Size
-   Created date
-   Referenced by running container
-   Safe score

### Containers

-   Name
-   Running / Stopped
-   Writable layer size
-   Log file size
-   Uptime
-   Restart policy

### Volumes

-   Name
-   Size
-   Attached containers
-   Last modified

------------------------------------------------------------------------

## 5.3 Safe-To-Delete Classification

Each resource classified as:

  Status      Definition
  ----------- ----------------------
  SAFE        No active references
  REVIEW      Possibly used
  PROTECTED   Currently in use

------------------------------------------------------------------------

## 5.4 Log Analysis

Command:

    dockit logs

Must:

-   Identify top largest logs
-   Detect containers writing excessive logs
-   Suggest log rotation

------------------------------------------------------------------------

## 5.5 Cleanup Mode

Command:

    dockit clean

Default: Dry-run only.

Apply mode:

    dockit clean --apply

Safety requirements:

-   Cannot delete PROTECTED resources
-   Double confirmation required
-   Show reclaimed estimate

------------------------------------------------------------------------

## 5.6 JSON Output

    dockit summary --json
    dockit analyze --json

For CI automation and monitoring integration.

------------------------------------------------------------------------

# 6. CLI Design

## Commands

    dockit summary
    dockit analyze
    dockit clean
    dockit logs
    dockit report --json
    dockit version

## Global Flags

    --json
    --verbose
    --dry-run
    --apply
    --threshold-days
    --log-threshold-mb

------------------------------------------------------------------------

# 7. Technical Architecture (Go)

## 7.1 Dependencies

-   Docker Go SDK
-   Cobra (CLI framework)
-   Optional: Viper (config)
-   Humanize package

## 7.2 Project Structure

    dockit/
     ├── cmd/
     ├── internal/
     │    ├── dockerclient/
     │    ├── analyzer/
     │    ├── scorer/
     │    ├── cleaner/
     │    ├── logger/
     │    └── formatter/
     ├── pkg/models/
     └── main.go

------------------------------------------------------------------------

# 8. Performance Requirements

-   Handle 200+ containers
-   Handle 1000+ images
-   Memory usage \< 150MB
-   Summary execution \< 5 seconds

------------------------------------------------------------------------

# 9. Security & Safety

-   Never delete running container resources
-   Never delete attached volumes
-   Default mode is non-destructive
-   Graceful permission handling

------------------------------------------------------------------------

# 10. Distribution Strategy

Initial:

-   GitHub Releases (static binary)
-   Linux (amd64, arm64)
-   macOS
-   Windows

Phase 2:

-   Homebrew formula
-   apt repository
-   Docker image distribution

------------------------------------------------------------------------

# 11. Version Roadmap

-   v0.1.0 → Summary + Analyze
-   v0.2.0 → Log analysis
-   v0.3.0 → Cleanup mode
-   v1.0.0 → Stable scoring engine + JSON output

------------------------------------------------------------------------

# 12. Long-Term Vision

dockit can evolve into:

-   Docker observability dashboard
-   Log growth monitoring agent
-   SaaS Docker health platform
-   Small-team DevOps toolkit

------------------------------------------------------------------------

**End of PRD**
