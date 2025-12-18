# AutoEnumeration - Quick Start Guide

## Installation

### 1. Clone and Install
```bash
./scripts/install.sh
```

This will:
- Install all Go-based security tools
- Build the AutoEnumeration binary
- Install it to `/usr/local/bin/autoenum`
- Create default configuration at `~/.autoenum/config.yaml`

### 2. Verify Installation
```bash
autoenum --help
autoenum list
```

## Basic Usage

### Run Your First Scan
```bash
# Quick scan (fast subdomain and port discovery)
autoenum scan -t example.com -m quick

# Standard scan (recommended)
autoenum scan -t example.com -m standard

# Deep scan (comprehensive with vulnerability detection)
autoenum scan -t example.com -m deep
```

### View Results
Results are saved in the `./output` directory:
- `scan_results_*.json` - Structured JSON data
- `scan_report_*.html` - Interactive HTML report (open in browser)

```bash
# Open HTML report
open output/scan_report_*.html
```

## Common Commands

### List Available Tools
```bash
autoenum list                    # All tools
autoenum list -c subdomain       # By category
```

### Install/Update Tools
```bash
autoenum install --all           # Install all tools
autoenum install --tools subfinder,httpx,nuclei
./scripts/update.sh              # Update all tools
```

### Custom Scans
```bash
# Custom output directory
autoenum scan -t example.com -o /path/to/output

# More concurrent workers
autoenum scan -t example.com --threads 20

# Skip specific tools
autoenum scan -t example.com --skip nuclei,nikto

# Run only specific tools
autoenum scan -t example.com --only subfinder,httpx,naabu

# Verbose output
autoenum scan -t example.com -v

# Custom timeout
autoenum scan -t example.com --timeout 7200
```

## Project Structure

```
AutoEnumerations/
├── cmd/                    # CLI commands
├── internal/               # Core implementation
│   ├── config/            # Configuration
│   ├── engine/            # Scan engine
│   ├── tools/             # Tool integrations
│   ├── reporter/          # Report generation
│   ├── logger/            # Logging
│   └── installer/         # Tool installation
├── scripts/               # Installation scripts
├── config.yaml            # Default configuration
├── README.md              # Full documentation
├── EXAMPLES.md            # Usage examples
└── Makefile               # Build automation
```

## Build from Source

```bash
# Using Makefile
make build          # Build binary
make install        # Build and install
make clean          # Clean artifacts
make test           # Run tests

# Manual build
go build -o autoenum main.go
sudo mv autoenum /usr/local/bin/
```

## Configuration

Edit `~/.autoenum/config.yaml` or create a custom config:

```yaml
output_dir: ./output
verbose: false
timeout: 3600

concurrency:
  max_workers: 10
  queue_size: 100

modes:
  custom:
    description: "My custom mode"
    tools:
      - subfinder
      - httpx
      - nuclei
```

Use custom config:
```bash
autoenum scan -t example.com --config /path/to/config.yaml
```

## Integrated Tools (30+)

### Subdomain Discovery
subfinder, assetfinder, amass, findomain, shuffledns, puredns, crobat

### Port Scanning
naabu, nmap, rustscan, masscan

### HTTP Probing
httpx, httprobe, feroxbuster, ffuf, gobuster, dirsearch

### DNS Enumeration
dnsx, dnsrecon, massdns, zdns

### Web Crawling
katana, gospider, hakrawler, waybackurls, gau, cariddi

### Vulnerability Scanning
nuclei, nikto

### Technology Detection
webanalyze, whatweb, wappalyzer

### Screenshots
gowitness, aquatone, eyewitness

### SSL/TLS
tlsx, sslyze, testssl.sh

### Cloud
cloudlist, s3scanner, cloud-enum

## Tips

1. **Start with quick mode** to map the attack surface
2. **Use standard mode** for regular reconnaissance
3. **Use deep mode** for comprehensive audits
4. **Check HTML reports** for visual analysis
5. **Customize scan modes** in config.yaml for your workflow

## Troubleshooting

### Dependencies Not Found
```bash
go mod download
go mod tidy
```

### Tools Not Installed
```bash
autoenum install --all
```

### Permission Issues
```bash
sudo chmod +x scripts/install.sh
sudo chmod +x scripts/update.sh
```

### Build Errors
```bash
make clean
make build
```

## Next Steps

1. Read `README.md` for full documentation
2. Check `EXAMPLES.md` for advanced workflows
3. Customize `config.yaml` for your needs
4. Run `autoenum list` to see all available tools

## Support

- GitHub Issues: Report bugs or request features
- Documentation: See README.md and EXAMPLES.md
- Configuration: Check config.yaml for all options

---

**Happy Hacking! 🎯**
