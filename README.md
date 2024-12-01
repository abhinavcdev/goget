# goget

A modern replacement for wget with support for concurrency, resumable downloads, JSON output, and pluggable storage backends like S3.

## Features

- **Concurrent Downloads**: Download files using multiple connections for faster speeds
- **Resumable Downloads**: Continue interrupted downloads from where they left off
- **Progress Visualization**: Real-time progress bars for each download chunk
- **Pluggable Storage Backends**: Save downloads to local files or Amazon S3
- **JSON Output**: Machine-readable output for integration with other tools

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/goget.git

# Navigate to the project directory
cd goget

# Build the binary
go build -o goget

# Optionally, install it to your GOPATH
go install
```

## Usage

### Basic Download

```bash
# Download a file to the current directory with default settings
goget https://example.com/largefile.zip

# Specify an output filename
goget -o myfile.zip https://example.com/largefile.zip
```

### Advanced Options

```bash
# Set the number of concurrent connections
goget -c 16 https://example.com/largefile.zip

# Set the number of retries for failed chunks
goget -r 10 https://example.com/largefile.zip

# Enable JSON output for machine readability
goget --json https://example.com/largefile.zip
```

### Storage Backends

```bash
# Download to local file (default)
goget --output-type file https://example.com/largefile.zip

# Download to Amazon S3 (requires AWS credentials in environment)
goget --output-type s3 --s3-bucket mybucket --s3-region us-west-2 https://example.com/largefile.zip

# Download directly from S3 (requires AWS credentials in environment)
goget s3://mybucket/path/to/file.zip
goget --s3-region us-west-2 s3://mybucket/path/to/file.zip
```

## Full Command Reference

```
Usage:
  goget [URL] [flags]

Flags:
  -c, --concurrency int       Number of concurrent download connections (default 8)
  -h, --help                  Help for goget
      --json                  Enable machine-readable JSON output
  -o, --output string         Output file name (defaults to filename from URL)
      --output-type string    Output type: file, s3, or stdout (default "file")
  -r, --retries int           Number of retries on failed chunk download (default 5)
      --s3-bucket string      S3 bucket name (required when output-type is s3)
      --s3-region string      S3 region (used when output-type is s3) (default "us-east-1")
```

## Benchmark Results

Benchmark tests comparing goget against wget and curl for downloading a ~58MB file:

### HTTP Download Performance

| Tool | Configuration | Time (seconds) |
|------|--------------|----------------|
| goget | 1 connection | 8.67 |
| goget | 4 connections | 7.81 |
| goget | 8 connections | 4.05 |
| goget | 16 connections | 4.04 |
| wget | default | 7.59 |
| curl | default | N/A (failed) |

### S3 Download Performance

| Tool | Configuration | Time (seconds) |
|------|--------------|----------------|
| goget | 1 connection | 3.59 |
| goget | 4 connections | 3.74 |
| goget | 8 connections | 3.35 |
| goget | 16 connections | 3.55 |
| wget | with presigned URL | 22.80 |

### Key Findings

- For HTTP downloads, goget with 8+ connections is nearly twice as fast as wget
- For S3 downloads, goget is 6-7x faster than wget using presigned URLs
- Increasing concurrency beyond 8 connections doesn't provide significant benefits

## Development

### Project Structure

```
goget/
├── cmd/
│   └── root.go          # Cobra CLI command definitions
├── internal/
│   ├── downloader/
│   │   └── downloader.go  # Core download logic, concurrency, retries
│   ├── output/
│   │   ├── writer.go      # The pluggable writer interface
│   │   ├── file.go        # Implementation for local file output
│   │   └── s3.go          # Implementation for S3 output
│   └── ui/
│       └── progress.go    # Progress bar management
├── main.go                # Main application entry point
├── go.mod
└── go.sum
```

### Requirements

- Go 1.16 or higher
- Dependencies:
  - github.com/spf13/cobra
  - github.com/vbauerster/mpb/v8

## License

MIT
# Initial project setup on December 1, 2024
