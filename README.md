# Elves

[![Go Version](https://img.shields.io/github/go-mod/go-version/nukoneko-tarou/elves)](https://go.dev/)
[![CI](https://github.com/nukoneko-tarou/elves/actions/workflows/go-test.yml/badge.svg)](https://github.com/nukoneko-tarou/elves/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nukoneko-tarou/elves)](https://goreportcard.com/report/github.com/nukoneko-tarou/elves)
[![License](https://img.shields.io/github/license/nukoneko-tarou/elves)](LICENSE)

A fast, lightweight CLI tool that scaffolds project directory structures from JSON files compatible with UNIX `tree -J` output.

Ideal for scaffolding boilerplate projects, replicating directory templates, or sharing project structures across teams.

---

## Features

- ⚡ **High Performance & Concurrent**: Tuned concurrent workers for fast directory tree creation with minimal OS lock contention.
- 🌲 **`tree -J` Compatible**: Works seamlessly with standard `tree -J` JSON output.
- 🛡️ **Secure by Default**: Built-in boundary validation prevents path traversal attacks (`..` or separator injections).
- 🔍 **Dry-Run Preview**: Preview the directory tree in terminal ASCII art before touching the filesystem.
- 📁 **Flexible Configuration**: Supports custom subdirectories, permissions, `.gitkeep` auto-creation, and concurrency tuning.

---

## Installation

### Homebrew (macOS / Linux)

```shell
brew tap nukoneko-tarou/cli-tool
brew install nukoneko-tarou/cli-tool/elves
```

### Go Install

```shell
go install github.com/nukoneko-tarou/elves@latest
```

### Pre-built Binaries

Download pre-compiled binaries from [GitHub Releases](https://github.com/nukoneko-tarou/elves/releases).

---

## Quick Start

### 1. Generate JSON from an existing project

Export your project directory structure using UNIX `tree`:

```shell
tree -J -d > structure.json
```

### 2. Scaffold directories using Elves

Navigate to the target destination and run `create`:

```shell
elves create structure.json
```

---

## Usage & Commands

```shell
elves [command] [flags]
```

### Available Commands

| Command | Description |
| :--- | :--- |
| `create <file>` | Creates the directory structure from the specified JSON file |
| `version` | Displays the current version of Elves |
| `help` | Shows help message for commands |

---

### `create` Options

```shell
elves create <path_to_json_file> [flags]
```

| Flag | Shorthand | Default | Description |
| :--- | :---: | :---: | :--- |
| `--dry-run` | `-d` | `false` | Previews the directory tree in ASCII format without making changes |
| `--sub` | `-s` | `""` | Creates the directory structure inside a specified subdirectory |
| `--permission` | `-p` | `"755"` | Sets the file mode permission for generated directories (e.g. `777`) |
| `--gitkeep` | `-g` | `false` | Automatically creates `.gitkeep` files in every generated directory |
| `--concurrency` | `-c` | `2` | Number of concurrent workers for directory creation (`1` for synchronous) |

---

## Examples

### 1. Previewing with `--dry-run`

Inspect the generated tree structure before creating directories on disk:

```shell
$ elves create ./sample.json --dry-run
.
├── api
├── assets
├── build
│   ├── ci
│   └── package
├── cmd
│   └── _your_app_
...
└── website
```

### 2. Scaffolding in a Subdirectory

Create a new project folder `my-new-app` and build the hierarchy inside it:

```shell
elves create ./sample.json --sub my-new-app
```

### 3. Preserving Empty Directories in Git

Generate `.gitkeep` files inside each directory so Git tracks the structure:

```shell
elves create ./sample.json --gitkeep
```

### 4. Custom Permissions

Set custom directory permissions (e.g. `0777`):

```shell
elves create ./sample.json --permission 777
```

### 5. High-Throughput Scaffolding

Adjust worker concurrency for large trees or network filesystems:

```shell
elves create ./sample.json --concurrency 4
```

---

## JSON Format Specification

Elves parses JSON compatible with the output of `tree -J`:

```json
[
  {
    "type": "directory",
    "name": ".",
    "contents": [
      {
        "type": "directory",
        "name": "cmd",
        "contents": [
          { "type": "directory", "name": "app" }
        ]
      },
      { "type": "file", "name": "README.md" }
    ]
  },
  { "type": "report", "directories": 2, "files": 1 }
]
```

- Top-level array containing the root directory node.
- Entries with `"type": "directory"` and `"name": "<dir_name>"` are created.
- Child entries are nested within `"contents": [...]`.

---

## License

This project is licensed under the [MIT License](LICENSE).