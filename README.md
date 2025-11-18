# Scaffold - Go Project Generator

A command-line tool for scaffolding new Go projects following the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) structure.

## Features

- Creates well-structured Go projects with standard directory layout
- Support for different project types (CLI, Library, Service)
- Interactive mode with beautiful prompts
- Dry-run mode to preview project structure
- Minimal mode for simpler projects
- Automatic generation of:
  - `go.mod` with correct module path
  - `README.md` with project structure documentation
  - `Makefile` with common development tasks
  - `.gitignore` with sensible defaults
  - Sample `main.go` in `/cmd` directory

## Installation

### From Source

```bash
git clone https://github.com/daedalus/scaffold.git
cd scaffold
make install
```

### Using Go Install

```bash
go install github.com/daedalus/scaffold/cmd/scaffold@latest
```

## Usage

### Interactive Mode

The easiest way to create a new project is using interactive mode. Simply run:

```bash
scaffold new
```

This automatically enters interactive mode and will prompt you for:
- Project name
- Module path
- Author name
- License
- Project type (CLI, Library, or Service)

You can also use the `-i` flag with a project name:

```bash
scaffold new myproject -i
```

### Non-Interactive Mode

You can also provide all options via flags:

```bash
scaffold new myproject \
  --module-path github.com/username/myproject \
  --author "Your Name" \
  --license MIT \
  --type cli
```

### Project Types

**CLI Application** (`--type cli`)
- Full directory structure including `/cmd`, `/internal`, `/pkg`, `/api`, `/configs`, `/scripts`, `/build`, `/test`, `/docs`, `/examples`
- Sample main.go with CLI boilerplate
- Ideal for command-line tools

**Library** (`--type library`)
- Focused structure for reusable code: `/internal`, `/pkg`, `/cmd`, `/scripts`, `/test`, `/docs`, `/examples`
- Emphasizes public (`/pkg`) and private (`/internal`) library code
- Perfect for Go packages and modules

**Service/API** (`--type service`)
- Full structure with `/web` and `/api` directories plus all infrastructure directories
- Suitable for web services, APIs, and long-running services

### Dry Run

Preview what will be created without making any changes:

```bash
scaffold new myproject --dry-run
```

## Examples

### Create a project interactively

```bash
scaffold new
# Follow the prompts to configure your project
```

### Create a CLI application

```bash
scaffold new mycli \
  --type cli \
  --module-path github.com/user/mycli \
  --author "John Doe"
```

### Create a library

```bash
scaffold new mylib \
  --type library \
  --module-path github.com/user/mylib
```

### Interactive mode with a specific project name

```bash
scaffold new myservice -i
# Follow the prompts to configure your project
```

## Generated Project Structure

### CLI Application Structure

```
myproject/
├── cmd/
│   └── myproject/
│       └── main.go
├── internal/
├── pkg/
├── api/
├── configs/
├── scripts/
├── build/
├── test/
├── docs/
├── examples/
├── README.md
├── Makefile
├── .gitignore
└── go.mod
```

### Library Structure

```
mylib/
├── internal/
├── pkg/
├── cmd/
├── scripts/
├── test/
├── docs/
├── examples/
├── README.md
├── Makefile
├── .gitignore
└── go.mod
```

### Service/API Structure

```
myservice/
├── cmd/
│   └── myservice/
│       └── main.go
├── internal/
├── pkg/
├── api/
├── web/
├── configs/
├── scripts/
├── build/
├── test/
├── docs/
├── examples/
├── README.md
├── Makefile
├── .gitignore
└── go.mod
```

## Directory Purposes

| Directory | Purpose |
|-----------|---------|
| `/cmd` | Main applications for the project |
| `/internal` | Private application and library code |
| `/pkg` | Library code that's ok to use by external applications |
| `/api` | API definitions (OpenAPI/Swagger specs, protocol definitions) |
| `/web` | Web application specific components (services only) |
| `/configs` | Configuration file templates |
| `/scripts` | Scripts for build, install, analysis, etc. |
| `/build` | Packaging and continuous integration |
| `/test` | Additional external test apps and test data |
| `/docs` | Design and user documents |
| `/examples` | Examples for your applications and/or libraries |

## Command Reference

### Global Flags

- `--help, -h` - Show help information
- `--version` - Show version information

### `new` Command

Create a new Go project.

**Usage:** `scaffold new [project-name] [flags]`

**Flags:**
- `-m, --module-path string` - Go module path (e.g., github.com/user/project)
- `-a, --author string` - Author name
- `-l, --license string` - License type (default "MIT")
- `-t, --type string` - Project type: cli, library, or service (default "cli")
- `--dry-run` - Preview what would be created without creating anything
- `-i, --interactive` - Use interactive mode to configure the project

## Development

### Prerequisites

- Go 1.21 or higher

### Building

```bash
make build
```

### Running Tests

```bash
make test
```

### Installing Locally

```bash
make install
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Related Projects

- [golang-standards/project-layout](https://github.com/golang-standards/project-layout) - Standard Go project layout reference

## Author

Your Name

> [!NOTE]
> This project was created with the assistance of Claude Code, an AI-powered development assistant by Anthropic.
