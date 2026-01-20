# tfpl - Terraform Polylith CLI

`tfpl` (Terraform Polylith) is a CLI tool for working with polylith-style Terraform repositories. It simplifies running Terraform/OpenTofu commands on components, bases, and projects organized in a polylith structure.

## Features

- Run `init`, `fmt`, and `validate` commands on components, bases, or projects
- Recursive search for nested modules (e.g., `iac/components/azurerm/storage-account`)
- Name clash detection with helpful error messages
- Support for both Terraform and OpenTofu
- Configuration file support (`.tfpl.yml`)
- Clean, intuitive CLI interface

## Installation

### Using go install

```bash
go install github.com/TechnicallyJoe/sturdy-parakeet/tools/tfpl@latest
```

### Build from source

```bash
git clone https://github.com/TechnicallyJoe/sturdy-parakeet.git
cd sturdy-parakeet/tools/tfpl
go build -o tfpl
```

Then move the `tfpl` binary to a directory in your PATH:

```bash
# Linux/macOS
sudo mv tfpl /usr/local/bin/

# Windows
# Move tfpl.exe to a directory in your PATH
```

## Requirements

- Go 1.25 or later (for building)
- Terraform CLI or OpenTofu CLI installed and in PATH

## Usage

### Commands

| Command | Aliases | Description |
|---------|---------|-------------|
| `tfpl init` | | Run terraform/tofu init on a module |
| `tfpl fmt` | | Run terraform/tofu fmt on a module |
| `tfpl val` | `validate` | Run terraform/tofu validate on a module |
| `tfpl config` | | Show current configuration |

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--component` | `-c` | Component name to operate on |
| `--base` | `-b` | Base name to operate on |
| `--project` | `-p` | Project name to operate on |
| `--path` | | Explicit path (mutually exclusive with -c, -b, -p) |
| `--init` | `-i` | Run init before the command (for `fmt` and `val`) |
| `--version` | `-v` | Show version |
| `--help` | `-h` | Show help |

### Examples

```bash
# Run fmt on a component
tfpl fmt -c storage-account

# Run validate on a base
tfpl val -b k8s-argocd

# Run init then validate on a base
tfpl val -i -b k8s-argocd

# Run init on a base
tfpl init -b k8s-argocd

# Run fmt on explicit path
tfpl fmt --path iac/components/azurerm/storage-account

# Show current configuration
tfpl config
```

## Configuration File

Create a `.tfpl.yml` file in your repository root to configure tfpl:

```yaml
# The root directory containing the polylith structure (components, bases, projects)
# Default: "" (repository root)
root: iac

# The Terraform binary to use: "terraform" or "tofu"
# Default: "terraform"
binary: terraform
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `root` | string | `""` | Directory containing components/bases/projects (relative to repo root) |
| `binary` | string | `"terraform"` | Binary to use: `"terraform"` or `"tofu"` |

tfpl will automatically search for `.tfpl.yml` starting from the current directory and walking up the directory tree until it finds one or reaches the filesystem root.

## Expected Directory Structure

tfpl expects your Terraform repository to follow a polylith-style structure:

```
repository-root/
├── .tfpl.yml              # Configuration file (optional)
└── iac/                   # Configured as "root: iac" (optional)
    ├── components/
    │   ├── aws/
    │   │   ├── s3-bucket/
    │   │   │   ├── main.tf
    │   │   │   └── variables.tf
    │   │   └── lambda/
    │   └── azurerm/
    │       ├── storage-account/
    │       │   ├── main.tf
    │       │   └── variables.tf
    │       └── naming/
    ├── bases/
    │   ├── azsloth/
    │   │   └── main.tf
    │   └── azsloth-docker-translator/
    └── projects/
        └── spacelift-modules/
            └── main.tf
```

### Module Discovery

- **Components** are searched in `<root>/components/`
- **Bases** are searched in `<root>/bases/`
- **Projects** are searched in `<root>/projects/`

tfpl recursively searches subdirectories to find modules with matching names. A directory is considered a valid module if it contains `.tf` or `.tf.json` files.

## Handling Edge Cases

### Nested Subfolders

Components, bases, or projects can be nested in subfolders:

```bash
# These both work
tfpl fmt -c storage-account  # Finds iac/components/azurerm/storage-account
tfpl fmt -c s3-bucket        # Finds iac/components/aws/s3-bucket
```

### Name Clashes

If multiple modules share the same name in different locations, tfpl will detect the clash and provide a helpful error:

```
Error: multiple components named 'naming' found - name clash detected:
  1. /path/to/repo/iac/components/azurerm/naming
  2. /path/to/repo/iac/components/aws/naming

Please use --path to specify the exact path
```

To resolve, use the `--path` flag:

```bash
tfpl fmt --path iac/components/azurerm/naming
```

### Flag Mutual Exclusivity

- The `--path` flag is mutually exclusive with `-c`, `-b`, and `-p`
- Only one of `-c`, `-b`, or `-p` can be specified at a time
- If you violate these rules, tfpl will show a clear error message

## Development

### Project Structure

```
tools/tfpl/
├── cmd/
│   └── root.go           # Main command definitions using cobra
├── internal/
│   ├── config/
│   │   └── config.go     # YAML configuration loading
│   ├── finder/
│   │   └── finder.go     # Module discovery with recursive search
│   └── terraform/
│       └── terraform.go  # Terraform/tofu command execution
├── main.go               # Entry point
├── go.mod                # Go modules
└── README.md             # This file
```

### Running Tests

```bash
cd tools/tfpl
go test ./...
```

### Building

```bash
cd tools/tfpl
go build
```

## License

See repository license.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
