# Gloc - GitHub Repository Code Analyzer

A simple command-line tool that clones a GitHub repository and analyzes its lines of code using `cloc`. It provides a clean, tabular view of the codebase composition.

## Prerequisites

- Go 1.21 or later
- Git
- cloc (install with `sudo apt-get install cloc` on Ubuntu/Debian)

## Installation

```bash
go install bybobby.dev/gloc@latest
```

## Usage

Basic usage:

```bash
gloc https://github.com/username/repository
```

Show raw output:

```bash
gloc --raw https://github.com/username/repository
```

## Output Format

The tool displays a clean, tabular view of the codebase:

```
Language    Files    Code
JavaScript    12    1234
TypeScript     8     789
CSS            5     234
---            ---    ---
Total         25    2257
```

## Features

- Clones the specified GitHub repository to a temporary directory
- Analyzes the code using `cloc`
- Shows a clean, tabular view of files and code by language
- Option to show raw command output with `--raw` flag
- Automatically cleans up temporary files after analysis
- Validates GitHub repository URL
- Checks for required dependencies (git and cloc)
