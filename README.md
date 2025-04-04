# Docker Compose Manager (DCM)

DCM is a command-line utility for managing multiple Docker Compose projects across different directories.

## Features

- Find all Docker Compose projects in a specified directory
- List, start, stop, and check status of projects
- Save and manage favorite Docker Compose projects with aliases
- Use managed projects without specifying paths
- Colorful output with emojis for better visualization
- Concurrent operations for faster management of multiple projects

## Installation

### Using Go

```bash
# Install directly using Go
go install github.com/mitas/dcm/cmd/dcm@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/mitas/dcm.git
cd dcm

# Build using Makefile
make build

# Optional: Install to your $GOPATH/bin
make install
```

### Manual Build

```bash
# Clone the repository
git clone https://github.com/mitas/dcm.git
cd dcm

# Build the application
go build -o dcm ./cmd/dcm

# Optional: Move the binary to your PATH
sudo mv dcm /usr/local/bin/
```

## Usage

### Global Flags

```
-p, --path string   Root path to search for Docker Compose projects
-c, --config string Path to config file (default is ~/.config/dcm/config.yaml)
```

Note: When using managed projects, the `--path` flag is not required.

### List Docker Compose Projects

List all Docker Compose projects in a directory:

```bash
dcm --path /path/to/projects list
```

Example output:
```
📋 Found 17 Docker Compose projects:
📁 1. project-a (/path/to/projects/project-a/docker-compose.yml)
📁 2. project-b (/path/to/projects/project-b/docker-compose.yml)
📁 3. project-c (/path/to/projects/project-c/docker-compose.yml)
📁 4. project-d (/path/to/projects/sub/project-d/docker-compose.yml)
📁 5. project-e (/path/to/projects/sub/project-e/docker-compose.yml)
...
```

### Start Projects

Start a specific project:

```bash
# Using project name flag
dcm --path /path/to/projects start --project myproject

# Or using positional argument
dcm --path /path/to/projects start myproject
```

Example output:
```
🔄 Starting Docker Compose project: myproject
✅ Successfully started myproject
```

Start all projects:

```bash
dcm --path /path/to/projects start --all
```

### Stop Projects

Stop a specific project:

```bash
# Using project name flag
dcm --path /path/to/projects stop --project myproject

# Or using positional argument
dcm --path /path/to/projects stop myproject
```

Example output:
```
🔄 Stopping Docker Compose project: myproject
✅ Successfully stopped myproject
```

Stop all projects:

```bash
dcm --path /path/to/projects stop --all
```

### Check Status

Check status of a specific project:

```bash
# Using project name flag
dcm --path /path/to/projects status --project myproject

# Or using positional argument
dcm --path /path/to/projects status myproject
```

Example output:
```
🔄 Checking status of Docker Compose project: myproject

=== Status of myproject (/path/to/projects/myproject) ===
🟢 myproject_redis: running (Up 7 seconds)
🟢 myproject_db: running (Up 7 seconds)
🟢 myproject_api: running (Up 7 seconds)
🟢 myproject_web: running (Up 7 seconds)
```

Check status of all projects:

```bash
dcm --path /path/to/projects status --all
```

### Managed Projects

#### Add a Project to Managed Projects

```bash
# Add a project with default alias (project name)
dcm --path /path/to/projects add-managed myproject

# Add a project with custom alias
dcm --path /path/to/projects add-managed myproject --alias prod-api
```

Example output:
```
✅ Project 'myproject' added to managed projects with alias 'prod-api'
```

#### List Managed Projects

```bash
# Using the full command
dcm list-managed

# Using aliases
dcm lsm
dcm lm
```

Example output:
```
📋 Managed Projects:
📌 1. project-a (alias) -> project-a (/path/to/projects/project-a)
📌 2. api (alias) -> myproject (/path/to/projects/myproject)
📌 3. prod-api (alias) -> myproject (/path/to/projects/myproject)
```

> Note: Each Docker Compose file path can only be added once to managed projects.

#### Use Managed Projects

Once projects are managed, you can start, stop, and check their status without specifying a path:

```bash
# Start a managed project
dcm start prod-api
```

Example output:
```
🔄 Starting managed Docker Compose project: prod-api
✅ Successfully started myproject
```

```bash
# Check status of a managed project
dcm status prod-api
```

Example output:
```
🔄 Checking status of managed Docker Compose project: prod-api

=== Status of myproject (/path/to/projects/myproject) ===
🟢 myproject_redis: running (Up 5 seconds)
🟢 myproject_db: running (Up 5 seconds)
🟢 myproject_api: running (Up 5 seconds)
🟢 myproject_web: running (Up 5 seconds)
```

```bash
# Stop a managed project
dcm stop prod-api
```

Example output:
```
🔄 Stopping managed Docker Compose project: prod-api
✅ Successfully stopped myproject
```

#### Remove a Managed Project

```bash
dcm remove-managed prod-api
```

Example output:
```
✅ Project with alias 'prod-api' removed from managed projects
```

## Complete Example Workflow

First, list all projects in your development directory:
```bash
dcm --path ~/dev list
```

Start a specific project:
```bash
dcm --path ~/dev start myproject
```

Check the status:
```bash
dcm --path ~/dev status myproject
```

Add the project to managed projects with an alias:
```bash
dcm --path ~/dev add-managed myproject --alias api
```

List managed projects:
```bash
dcm list-managed
```

Use the managed project (no path needed):
```bash
dcm start api
dcm status api
dcm stop api
```

Remove from managed projects when no longer needed:
```bash
dcm remove-managed api
```

## Build and Development

### Makefile Commands

The project includes a Makefile with the following commands:

```bash
# Build the binary
make build

# Format, lint, and build the code
make all

# Run tests
make test

# Format the code
make fmt

# Run linter
make lint

# Build for multiple platforms (Linux, macOS, Windows)
make release

# Show all available commands
make help
```

## Project Structure

The project follows a clean architecture with the following structure:

```
dcm/
├── cmd/
│   └── dcm/              # Application entry point
├── internal/
│   ├── cli/              # CLI interface
│   ├── cmd/              # Command implementations
│   ├── config/           # Configuration handling
│   ├── manager/          # Business logic
│   └── model/            # Data structures
└── pkg/
    └── formatter/        # Output formatting
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.