# kube-disk-stats

A CLI tool for querying Kubernetes node and pod disk usage statistics. kube-disk-stats connects to your Kubernetes cluster and retrieves storage information from the `/api/v1/nodes/{node}/proxy/stats/summary` endpoint.

## Features

- Query storage usage for all nodes or a specific node
- Display pod-level storage breakdown
- Container-level storage details (rootfs + logs)
- Color-coded output for high-usage nodes (green < 60%, yellow 60-79%, red >= 80%)
- Sort by top N consumers
- JSON output support
- Watch mode for continuous monitoring
- Multiple Kubernetes context support

## Installation

### From Binary

Download the latest release from the [GitHub Releases](https://github.com/aldi-f/kube-disk-stats/releases) page.

### From Source

```bash
go install github.com/aldi-f/kube-disk-stats@latest
```

### From Homebrew (Linux/Mac)

```bash
brew tap aldi-f/tap
brew install kube-disk-stats
```

## Usage

### Basic Usage

Display storage usage for all nodes:

```bash
kube-disk-stats
```

Query a specific node:

```bash
kube-disk-stats --node ip-10-244-12-44
```

Use a different Kubernetes context:

```bash
kube-disk-stats --context my-cluster-context
```

### Subcommands

Display pod storage usage:

```bash
kube-disk-stats pods
```

Display node storage usage:

```bash
kube-disk-stats nodes
```

Display container storage usage:

```bash
kube-disk-stats containers
```

### Filtering and Sorting

Show top 10 pods by storage usage:

```bash
kube-disk-stats pods --top 10
```

Show top 5 nodes by percentage usage:

```bash
kube-disk-stats nodes --top 5
```

### Output Formats

JSON output:

```bash
kube-disk-stats --output json
```

### Watch Mode

Watch mode continuously refreshes the display:

```bash
kube-disk-stats --watch --interval 10s
```

## CLI Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--context` | `-c` | Kubernetes context to use | current context |
| `--node` | `-n` | Query specific node | all nodes |
| `--output` | `-o` | Output format (table/json) | table |
| `--top` | `-t` | Show top N results | all |
| `--watch` | `-w` | Watch mode | false |
| `--interval` | `-i` | Refresh interval for watch mode | 5s |

## Examples

View top 10 pods with color output:

```bash
kube-disk-stats pods --top 10
```

Watch all nodes with 30-second refresh:

```bash
kube-disk-stats nodes --watch --interval 30s
```

Query specific context and export JSON:

```bash
kube-disk-stats --context prod --output json > stats.json
```

## Development

### Build

```bash
make build
```

### Build for all platforms

```bash
make build-all
```

### Test

```bash
make test
```

### Lint

```bash
make lint
```

### Create Release

```bash
make release
```

## GitHub Actions

This project includes GitHub Actions workflows for:

- **Build**: Runs tests, lints code, and builds binaries for multiple platforms on push/PR
- **Release**: Creates GitHub releases and builds release artifacts on tags

### Creating a Release

Push a git tag to trigger the release workflow:

```bash
git tag v1.0.0
git push origin v1.0.0
```

Or use the GitHub Actions workflow_dispatch to create a release manually.

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
