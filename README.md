# Elves

Elves is a command-line tool that generates a project directory structure from a JSON file. This is particularly useful for quickly scaffolding new projects based on a predefined template. The JSON file format is compatible with the output of the `tree -J` command.

## Installation

You can install Elves using one of the following methods:

### Go

```shell
go install github.com/nukoneko-tarou/elves@latest
```

### Homebrew

```shell
brew tap nukoneko-tarou/cli-tool
brew install nukoneko-tarou/cli-tool/elves
```

## Usage

To generate a directory structure, navigate to your desired parent directory and run the `create` command, specifying the path to your JSON file.

```shell
cd <target_directory>
elves create <path_to_json_file>
```

### Example

Using the `sample.json` file included in this repository (based on the [golang-standards/project-layout](https://github.com/golang-standards/project-layout)), you can create the following structure:

```shell
elves create ../sample.json
.
├── api
├── assets
├── build
│   ├── ci
│   └── package
├── cmd
│   └── _your_app_
├── configs
├── deployments
├── docs
├── examples
├── githooks
├── init
├── internal
│   ├── app
│   │   └── _your_app_
│   └── pkg
│       └── _your_private_lib_
├── pkg
│   └── _your_public_lib_
├── scripts
├── test
├── third_party
├── tools
├── vendor
├── web
│   ├── app
│   ├── static
│   └── template
└── website

30 directories, 0 files
```

## Options

The `create` command supports the following options:

### `--sub, -s`

Creates a new directory with the specified name in the current location and generates the project structure inside it.

**Example:**

```shell
elves create ./sample.json --sub new-project
```

### `--permission, -p`

Sets the permissions for the generated directories. The default is `0755`.

**Example:**

```shell
elves create ./sample.json --permission 0777
```

### `--gitkeep, -g`

Creates a `.gitkeep` file in each generated directory. This is useful for ensuring that empty directories are tracked by Git.

**Example:**

```shell
elves create ./sample.json --gitkeep
```