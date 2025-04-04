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

```bash
# Clone the repository
git clone https://github.com/yourusername/dcm.git
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

```bash
dcm --path /path/to/projects list
```

### Start Projects

Start a specific project:

```bash
# Using project name flag
dcm --path /path/to/projects start --project myproject

# Or using positional argument
dcm --path /path/to/projects start myproject
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

#### List Managed Projects

```bash
dcm list-managed
```

#### Use Managed Projects

Once projects are managed, you can start, stop, and check their status without specifying a path:

```bash
# Start a managed project
dcm start my-alias

# Stop a managed project
dcm stop my-alias

# Check status of a managed project
dcm status my-alias
```

#### Remove a Managed Project

```bash
dcm remove-managed my-alias
```

## Examples

```bash
# Find all Docker Compose projects in your development directory
dcm --path ~/dev list

# Start a specific project
dcm --path ~/dev start myapi

# Stop all projects
dcm --path ~/dev stop --all

# Check the status of all projects
dcm --path ~/dev status --all

# Add a project to managed projects
dcm --path ~/dev add-managed myapi --alias api

# Start a managed project (no path needed)
dcm start api

# List all managed projects
dcm list-managed
```

## Project Structure

The project follows a clean architecture with the following structure:

```
dcm/
├── cmd/
│   └── dcm/              # Application entry point
├── internal/
│   ├── cmd/              # Command implementations
│   ├── manager/          # Business logic
│   └── model/            # Data structures
└── pkg/
    └── formatter/        # Output formatting
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.