# kube-disk-stats - Project Summary

## Project Overview

kube-disk-stats is a CLI tool for querying Kubernetes node and pod disk usage statistics. It connects to your Kubernetes cluster and retrieves storage information from the `/api/v1/nodes/{node}/proxy/stats/summary` endpoint.

## Project Structure

```
cli/
├── .github/
│   └── workflows/
│       ├── build.yml      # CI/CD: Build, test, lint on push/PR
│       └── release.yml   # Release workflow for tags
├── cmd/
│   └── root.go           # Cobra CLI commands and flags
├── internal/
│   ├── analyzer/
│   │   └── calculator.go # Storage calculations, age formatting, pod grouping
│   ├── display/
│   │   ├── color.go      # Color thresholds (green/yellow/red)
│   │   ├── json.go       # JSON output formatting
│   │   └── table.go      # Table output formatting
│   ├── k8s/
│   │   ├── client.go     # Kubernetes client with context support
│   │   └── stats.go      # Stats summary API fetching
│   └── models/
│       └── types.go      # Data structures
├── pkg/
│   └── sort/
│       └── sorter.go     # Top N sorting for pods/nodes
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
├── Makefile              # Build targets
├── .golangci.yml        # Linter configuration
├── .gitignore            # Git ignore patterns
├── LICENSE               # MIT License
├── README.md             # User documentation
├── CONTRIBUTING.md       # Contribution guidelines
└── homebrew-formula.rb   # Homebrew formula template
```

## Features Implemented

### Core Functionality
- ✅ Query storage usage for all nodes or specific node
- ✅ Display pod-level storage breakdown
- ✅ Container-level storage details (rootfs + logs)
- ✅ Color-coded output (green < 60%, yellow 60-79%, red >= 80%)
- ✅ Sort by top N consumers
- ✅ JSON output support
- ✅ Watch mode for continuous monitoring
- ✅ Multiple Kubernetes context support

### CLI Commands
- `kube-disk-stats` - Show all stats (nodes, pods)
- `kube-disk-stats nodes` - Show node storage usage only
- `kube-disk-stats pods` - Show pod storage usage only
- `kube-disk-stats containers` - Show container storage usage only
- `kube-disk-stats version` - Print version information

### CLI Flags
- `--context, -c` - Kubernetes context to use
- `--node, -n` - Query specific node
- `--output, -o` - Output format (table/json)
- `--top, -t` - Show top N results
- `--watch, -w` - Watch mode
- `--interval, -i` - Refresh interval for watch mode

## Dependencies

```go
github.com/fatih/color v1.17.0       // Terminal colors
github.com/spf13/cobra v1.8.1        // CLI framework
github.com/olekukonko/tablewriter v0.0.5 // Table formatting
k8s.io/client-go v0.30.3             // Kubernetes client
k8s.io/api v0.30.3                   // Kubernetes API types
k8s.io/apimachinery v0.30.3           // Kubernetes utilities
```

## GitHub Actions Workflows

### Build Workflow (`.github/workflows/build.yml`)
- Triggers on push/PR to main or develop branches
- Runs on Ubuntu, macOS, and Windows
- Matrix builds for amd64 and arm64 (except Windows arm64)
- Steps:
  1. Checkout code
  2. Set up Go 1.23
  3. Cache Go modules
  4. Run tests with race detector and coverage
  5. Upload coverage to Codecov
  6. Build binaries
  7. Upload artifacts
  8. Run golangci-lint

### Release Workflow (`.github/workflows/release.yml`)
- Triggers on tag push (v*.*.*) or workflow_dispatch
- Permissions: contents: write
- Steps:
  1. Create tag (if using workflow_dispatch)
  2. Create GitHub Release with auto-generated notes
  3. Build release binaries for all platforms
  4. Upload release assets (tar.gz and zip)

## Building and Releasing

### Using Make

```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make test         # Run tests
make lint         # Run linters
make install      # Install to GOPATH/bin
make release      # Create release packages
make clean        # Remove build artifacts
```

### Manual Build

```bash
# Build for current platform
go build -o dist/kube-disk-stats .

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o dist/kube-disk-stats-linux-amd64 .
```

### Creating a Release

```bash
# Tag and push
git tag v1.0.0
git push origin v1.0.0

# Or use GitHub Actions workflow_dispatch manually
```

The release workflow will:
1. Create GitHub release with changelog
2. Build binaries for all platforms
3. Upload release artifacts

## Installation Methods

### From Binary
Download from GitHub Releases

### From Source
```bash
go install github.com/aldi-f/kube-disk-stats@latest
```

### Using Homebrew
```bash
brew tap aldi-f/tap
brew install kube-disk-stats
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

## Future Enhancements

Potential improvements:
1. Add support for filtering by namespace
2. Add support for filtering by pod label
3. Add support for exporting to CSV
4. Add support for persistent storage (PVC) statistics
5. Add support for customizing color thresholds
6. Add support for aggregating statistics across multiple clusters
7. Add support for historical data tracking
8. Add support for alerts/notifications

## Notes

- The tool uses the Kubernetes `/api/v1/nodes/{node}/proxy/stats/summary` endpoint
- Requires appropriate Kubernetes RBAC permissions to read node stats
- Works both inside Kubernetes clusters (in-cluster config) and locally (kubeconfig)
- Follows Kubernetes API conventions and best practices
- Uses semantic versioning for releases
