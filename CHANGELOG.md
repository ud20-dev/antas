# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org) and the format is based on [Keep a Changelog](https://keepachangelog.com).

---

## [1.0.0] - 2026-07-10

### Added

- CLI entrypoint with `-f / --format` flag supporting `human` and `json` output modes
- PDF-to-PNG rendering via go-pdfium (WebAssembly and native CGO build variants)
- Output path convention: `{tmp}/{sha256_of_file}/{timestamp}-{pid}-{seq}/page_{n}.png` grouping renders by file content and guaranteeing uniqueness across concurrent calls
- Exit code contract: `0` success, `1` runtime error, `2` bad CLI usage — JSON format always produces valid JSON on stdout for exit codes 0 and 1
- `-v / --version` flag reporting build version, OS/arch, build time, and commit hash
- Docker image for the native CGO build (`antas-natif`) via multi-stage Dockerfile
