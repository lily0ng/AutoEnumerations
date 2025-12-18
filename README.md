# AutoEnumeration

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**A comprehensive reconnaissance framework integrating 30+ security tools for automated enumeration and vulnerability discovery**

</div>

## 🎯 Features

- **🔧 Modular Architecture**: Plugin-based design for easy tool integration
- **⚡ Parallel Execution**: Concurrent tool execution with intelligent rate limiting
- **🎨 Multiple Scan Modes**: Quick, Standard, and Deep scanning profiles
- **📊 Rich Reporting**: JSON and HTML report generation
- **🔄 Result Aggregation**: Automatic deduplication and correlation of findings
- **⏸️ Resume Capability**: Continue interrupted scans
- **🎛️ Flexible Configuration**: YAML-based configuration with per-tool settings
- **📈 Progress Tracking**: Real-time scan progress and status updates

## 📦 Tool Categories

### Subdomain Discovery
- subfinder, assetfinder, amass, findomain, shuffledns, puredns, crobat

### Port Scanning
- naabu, nmap, rustscan, masscan

### HTTP Probing & Web Discovery
- httpx, httprobe, feroxbuster, ffuf, gobuster, dirsearch

### DNS Enumeration
- dnsx, dnsrecon, massdns, zdns

### Web Crawling
- katana, gospider, hakrawler, waybackurls, gau, cariddi

### Vulnerability Scanning
- nuclei, nikto

### Technology Detection
- webanalyze, whatweb, wappalyzer

### Screenshot Capture
- gowitness, aquatone, eyewitness

### SSL/TLS Analysis
- tlsx, sslyze, testssl.sh

### Cloud Enumeration
- cloudlist, s3scanner, cloud-enum

## 🚀 Installation

### Prerequisites

- Go 1.21 or later
- Git
- Bash (for installation script)

### Quick Install

```bash
git clone https://github.com/yourusername/AutoEnumerations.git
cd AutoEnumerations
chmod +x scripts/install.sh
./scripts/install.sh
```

### Manual Build

```bash
go mod download
go build -o autoenum main.go
sudo mv autoenum /usr/local/bin/
```

## 📖 Usage

### Basic Scan

```bash
autoenum scan -t example.com
```

### Scan Modes

#### Quick Scan (Fast subdomain and port discovery)
```bash
autoenum scan -t example.com -m quick
```

#### Standard Scan (Full enumeration with web discovery)
```bash
autoenum scan -t example.com -m standard
```

#### Deep Scan (Comprehensive with vulnerability detection)
```bash
autoenum scan -t example.com -m deep
```

### Advanced Options

```bash
# Custom output directory
autoenum scan -t example.com -o /path/to/output

# Specify number of concurrent workers
autoenum scan -t example.com --threads 20

# Skip specific tools
autoenum scan -t example.com --skip nuclei,nikto

# Run only specific tools
autoenum scan -t example.com --only subfinder,httpx,naabu

# Verbose output
autoenum scan -t example.com -v

# Custom timeout (seconds)
autoenum scan -t example.com --timeout 7200
```

### Tool Management

#### List Available Tools
```bash
autoenum list
```

#### List Tools by Category
```bash
autoenum list -c subdomain
```

#### Install Tools
```bash
# Install all tools
autoenum install --all

# Install specific tools
autoenum install --tools subfinder,httpx,nuclei

# Update existing tools
autoenum install --all --update
```

### Configuration

Create or edit `~/.autoenum/config.yaml` or use a custom config:

```bash
autoenum scan -t example.com --config /path/to/config.yaml
```

Example configuration:

```yaml
output_dir: ./output
verbose: false
timeout: 3600

concurrency:
  max_workers: 10
  queue_size: 100

rate_limit:
  requests_per_second: 100
  burst: 200

modes:
  custom:
    description: "My custom scan mode"
    tools:
      - subfinder
      - httpx
      - nuclei

tools:
  nuclei:
    enabled: true
    priority: 6
    timeout: 900
    args:
      severity: "critical,high"
```

## 📊 Output

AutoEnumeration generates two types of reports:

### JSON Report
- Structured data for programmatic processing
- Located at: `output/scan_results_YYYYMMDD_HHMMSS.json`

### HTML Report
- Beautiful, interactive web report
- Located at: `output/scan_report_YYYYMMDD_HHMMSS.html`

## 🔧 Architecture

```
AutoEnumeration/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── scan.go            # Scan command
│   ├── install.go         # Install command
│   └── list.go            # List command
├── internal/
│   ├── config/            # Configuration management
│   ├── engine/            # Core scan engine
│   │   ├── engine.go      # Main engine
│   │   └── aggregator.go  # Result aggregation
│   ├── tools/             # Tool integrations
│   │   ├── registry.go    # Tool registry
│   │   ├── subdomain.go   # Subdomain tools
│   │   ├── portscanning.go
│   │   ├── http.go
│   │   └── ...
│   ├── reporter/          # Report generation
│   ├── logger/            # Logging utilities
│   └── installer/         # Tool installation
├── scripts/               # Bash scripts
│   ├── install.sh         # Installation script
│   └── update.sh          # Update script
├── config.yaml            # Default configuration
└── main.go                # Entry point
```

## 🎯 Workflow Example

```bash
# 1. Install AutoEnumeration and tools
./scripts/install.sh

# 2. Run a quick reconnaissance
autoenum scan -t target.com -m quick -o ./target_recon

# 3. Review findings and run deep scan
autoenum scan -t target.com -m deep -o ./target_deep

# 4. View HTML report
open ./target_deep/scan_report_*.html

# 5. Update tools periodically
./scripts/update.sh
```

## 🔄 Integration with Other Tools

AutoEnumeration outputs can be easily piped to other tools:

```bash
# Extract subdomains for further processing
cat output/subfinder.txt | httprobe | tee live_hosts.txt

# Feed results to other scanners
cat output/httpx.txt | nuclei -t cves/
```

## 🛠️ Development

### Adding New Tools

1. Create tool definition in appropriate category file
2. Implement execute function
3. Register tool in registry
4. Update configuration defaults

Example:

```go
r.Register(&Tool{
    Name:        "mytool",
    Category:    "subdomain",
    Description: "My custom tool",
    InstallCmd:  "go install github.com/user/mytool@latest",
    ExecuteFunc: executeMyTool,
})
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o autoenum main.go
```

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## ⚠️ Disclaimer

This tool is for authorized security testing only. Always obtain proper authorization before scanning any targets.

## 📧 Contact

- GitHub Issues: [Report bugs or request features](https://github.com/yourusername/AutoEnumerations/issues)

## 🙏 Acknowledgments

- ProjectDiscovery team for their amazing tools
- All the open-source security tool developers
- The bug bounty and security research community

---

<div align="center">
Made with ❤️ for the security community
</div>
